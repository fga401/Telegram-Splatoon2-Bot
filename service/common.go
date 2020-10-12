package service

import (
	botapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"strconv"
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
	userMaxAccount            int
	userAllowPolling          bool
	callbackQueryCachedSecond int
	retryTimes                int
	defaultAdmin              int64
	storeChannelID            int64
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
	// markup
	initMarkup()

	//service
	loadUsers()
	tryStartJobScheduler()
}

func tryStartJobScheduler() {
	if salmonSchedules == nil {
		log.Info("start salmon job scheduler")
		startSalmonJobScheduler()
	}
}

func updateCookies(runtime *db.Runtime) (string, error) {
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

func fetchRuntime(uid int64) (*db.Runtime, error) {
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

type I18nKeys struct {
	Key  string
	Args []interface{}
}

func NewI18nKey(key string, args ...interface{}) I18nKeys {
	return I18nKeys{
		Key:  key,
		Args: args,
	}
}

func getI18nText(lang string, user *botapi.User, keys ...I18nKeys) []string {
	tag, err := language.Parse(lang)
	zapFields := make([]zap.Field, 0, 3)
	if user != nil {
		zapFields = append(zapFields, zap.Object("user", log.WrapUser(user)))
	}
	if err != nil {
		zapFields = append(zapFields, zap.String("language", lang), zap.Error(err))
		log.Warn("parse language failed", zapFields...)
		tag = language.English
	}
	printer := message.NewPrinter(tag)
	ret := make([]string, 0, len(keys))
	for _, key := range keys {
		ret = append(ret, printer.Sprintf(key.Key, key.Args...))
	}
	return ret
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
		return err
	}, retryTimes)
	if err != nil {
		err = errors.Wrap(err, "can't send message")
	}
	return &respMsg, err
}

var updateInterval = int64(2 * time.Hour.Seconds())

func getSplatoonNextUpdateTime(t time.Time) time.Time {
	nowTimestamp := t.Unix()
	nextTimestamp := (nowTimestamp/updateInterval + 1) * updateInterval
	// nextTimestamp += 5 // 5s delay
	return time.Unix(nextTimestamp, 0)
}

func getLocalTime(timestamp int64, offsetInMinute int) time.Time {
	return time.Unix(timestamp, 0).UTC().Add(time.Duration(offsetInMinute) * time.Minute)
}