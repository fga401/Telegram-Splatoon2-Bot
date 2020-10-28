package service

import (
	botapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"strconv"
	"telegram-splatoon2-bot/botutil"
	log "telegram-splatoon2-bot/logger"
	"telegram-splatoon2-bot/nintendo"
	"telegram-splatoon2-bot/service/cache"
	"telegram-splatoon2-bot/service/db"
	"time"
)

var (
	bot          *botapi.BotAPI
	Cache        *cache.CacheImpl
	AccountTable *db.AccountTableImpl
	UserTable    *db.UserTableImpl
	RuntimeTable *db.RuntimeTableImpl
	Transactions *db.TransactionImpl
)

var (
	userMaxAccount   int
	userAllowPolling bool
	defaultAdmin     int64
	storeChannelID   int64

	callbackQueryCachedSecond  int
	retryTimes                 int
	updateFailureRetryInterval time.Duration
	updateDelayInSecond        int64

	proposedStageNumber int
)

func InitService(b *botapi.BotAPI) {
	InitImageClient()
	db.InitDatabaseInstance()
	cache.InitCache()
	Cache = cache.Cache
	AccountTable = db.AccountTable
	UserTable = db.UserTable
	RuntimeTable = db.RuntimeTable
	Transactions = db.Transactions

	var err error
	// default value
	bot = b
	userMaxAccount = viper.GetInt("account.maxAccount")
	userAllowPolling = viper.GetBool("account.allowPolling")
	callbackQueryCachedSecond = viper.GetInt("service.callbackQueryCachedSecond")
	retryTimes = viper.GetInt("service.retryTimes")
	defaultAdmin, err = strconv.ParseInt(viper.GetString("admin"), 10, 64)
	if err != nil {
		panic(errors.Wrap(err, "viper get admin failed"))
	}
	storeChannelID, err = strconv.ParseInt(viper.GetString("store_channel"), 10, 64)
	if err != nil {
		panic(errors.Wrap(err, "viper get store_channel failed"))
	}
	updateFailureRetryInterval = viper.GetDuration("service.updateFailureRetryInterval")
	updateDelayInSecond = viper.GetInt64("service.updateDelayInSecond")
	proposedStageNumber = viper.GetInt("service.stage.proposedNumber")

	// markup
	initMarkup()

	//service
	loadUsers()

	// salmon
	salmonScheduleRepo, err = NewSalmonScheduleRepo(admins)
	if err != nil {
		panic(errors.Wrap(err, "can't init NewSalmonScheduleRepo"))
	}

	// stage
	stageScheduleRepo, err = NewStageScheduleRepo(admins)
	if err != nil {
		panic(errors.Wrap(err, "can't init NewStageScheduleRepo"))
	}

	Scheduler.tryStart()
}

func UpdateCookies(runtime *db.Runtime) (string, error) {
	// fetch
	var iksm string
	err := retry(func() error {
		var err error
		iksm, _, _, err = nintendo.GetCookiesAndNames(runtime.SessionToken, runtime.Language)
		return err
	}, retryTimes)
	if err != nil {
		return "", errors.Wrap(err, "can't get cookie")
	}
	// update
	if iksm == runtime.IKSM {
		return iksm, nil
	}
	err = RuntimeTable.UpdateRuntimeIKSM(runtime.Uid, iksm)
	if err != nil {
		return "", errors.Wrap(err, "can't update iksm to db")
	}
	log.Info("cookie updated", zap.Int64("user", runtime.Uid), zap.String("cookie", iksm))
	Cache.DeleteRuntime(runtime.Uid)
	return iksm, nil
}

func FetchRuntime(uid int64) (*db.Runtime, error) {
	runtime, err := Cache.GetRuntime(uid)
	// found in cache
	if err == nil && runtime != nil {
		return runtime, nil
	}
	if err != nil {
		log.Warn("can't fetch Runtime from cache", zap.Int64("uid", uid), zap.Error(err))
	}
	// try to fetch from db
	// todo: add metrics
	log.Info("runtime cache missed", zap.Int64("uid", uid))
	runtime, err = RuntimeTable.GetRuntime(uid)
	if err != nil {
		return nil, err
	}
	// set cache
	err = Cache.SetRuntime(runtime)
	if err != nil {
		log.Warn("can't set Runtime to cache", zap.Object("runtime", runtime), zap.Error(err))
	}
	return runtime, nil
}

func retry(handler func() error, times int) error {
	var err error
	for i := 0; i < times; i++ {
		err = handler()
		if err == nil {
			return nil
		}
	}
	return err
}

func sendWithRetry(bot *botapi.BotAPI, msg botapi.Chattable) error {
	err := retry(func() error {
		_, err := bot.Send(msg)
		if is, sec := botutil.IsTooManyRequestError(err); is {
			// todo: more info?
			log.Warn("send message blocked by telegram request limits", zap.Int("after", sec))
			time.Sleep(time.Duration(sec) * time.Second)
		}
		return err
	}, retryTimes)
	if err != nil {
		err = errors.Wrap(err, "can't send message")
	}
	return err
}

func sendWithRetryAndResponse(bot *botapi.BotAPI, msg botapi.Chattable) (*botapi.Message, error) {
	var respMsg botapi.Message
	err := retry(func() error {
		var err error
		respMsg, err = bot.Send(msg)
		if is, sec := botutil.IsTooManyRequestError(err); is {
			// todo: more info?
			log.Warn("send message blocked by telegram request limits", zap.Int("after", sec))
			time.Sleep(time.Duration(sec) * time.Second)
		}
		return err
	}, retryTimes)
	if err != nil {
		err = errors.Wrap(err, "can't send message")
	}
	return &respMsg, err
}

type timeHelper struct {updateInterval int64}
var TimeHelper = timeHelper{updateInterval: int64(2 * time.Hour.Seconds())}

func (helper timeHelper)getSplatoonNextUpdateTime(t time.Time) time.Time {
	nowTimestamp := t.Unix()
	nextTimestamp := (nowTimestamp/helper.updateInterval + 1) * helper.updateInterval
	return time.Unix(nextTimestamp, 0)
}

func (timeHelper)getLocalTime(timestamp int64, offsetInMinute int) time.Time {
	return time.Unix(timestamp, 0).In(time.FixedZone("", offsetInMinute * 60))
}

// func(iksm string, timezone int, acceptLang string, args ...interface{}) (result interface{}, error)
type Retriever func(string, int, string, ...interface{}) (interface{}, error)

func FetchResourceWithUpdate(uid int64, retriever Retriever, args ...interface{}) (interface{}, error) {
	runtime, err := FetchRuntime(uid)
	if err != nil {
		return nil, errors.Wrap(err, "can't fetch runtime")
	}

	var result interface{}
	err = retry(func() error {
		result, err = retriever(runtime.IKSM, runtime.Timezone, runtime.Language, args...)
		return err
	}, retryTimes)

	if errors.Is(err, &nintendo.ExpirationError{}) {
		// todo: add metric
		var iksm string
		iksm, err = UpdateCookies(runtime)
		if err != nil {
			return nil, errors.Wrap(err, "cookie expired and can't update it")
		}
		runtime.IKSM = iksm
		err = retry(func() error {
			result, err = retriever(runtime.IKSM, runtime.Timezone, runtime.Language, args...)
			return err
		}, retryTimes)
	}

	if err != nil {
		return nil, errors.Wrap(err, "can't get resources from nintendo")
	}

	return result, nil
}