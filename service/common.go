package service

import (
	botapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	log "telegram-splatoon2-bot/logger"
	"telegram-splatoon2-bot/service/cache"
	"telegram-splatoon2-bot/service/db"
)

var (
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
)

func InitService() {
	db.InitDatabaseInstance()
	cache.InitCache()
	Cache = cache.Cache
	AccountTable = db.AccountTable
	UserTable = db.UserTable
	RuntimeTable = db.RuntimeTable
	Transactions = db.Transactions

	// default value
	userMaxAccount = viper.GetInt("account.maxAccount")
	userAllowPolling = viper.GetBool("account.allowPolling")
	callbackQueryCachedSecond = viper.GetInt("service.callbackQueryCachedSecond")
	retryTimes = viper.GetInt("service.retryTimes")

	// markup
	initMarkup()
}

func fetchRuntime(user *botapi.User) (*db.Runtime, error) {
	runtime, err := Cache.GetRuntime(user)
	// found in cache
	if err == nil && runtime != nil {
		return runtime, nil
	}
	if err != nil {
		log.Warn("can't fetch Runtime from cache", zap.Object("user", log.WrapUser(user)), zap.Error(err))
	}
	// try to fetch from db
	// todo: add metrics
	log.Info("runtime cache missed", zap.Object("user", log.WrapUser(user)))
	runtime, err = RuntimeTable.GetRuntime(int64(user.ID))
	if err != nil {
		return nil, err
	}
	// set cache
	err = Cache.SetRuntime(user, runtime)
	if err != nil {
		log.Warn("can't set Runtime to cache", zap.Object("runtime", runtime), zap.Error(err))
	}
	return runtime, nil
}

type I18nKeys struct {
	Key  string
	Args []interface{}
}

func getI18nText(lanague string, user *botapi.User, keys ...I18nKeys) []string {
	tag, err := language.Parse(lanague)
	zapFields := make([]zap.Field, 0, 3)
	if user != nil {
		zapFields = append(zapFields, zap.Object("user", log.WrapUser(user)))
	}
	if err != nil {
		zapFields = append(zapFields, zap.String("language", lanague), zap.Error(err))
		log.Warn("parse language failed",
			zap.Object("user", log.WrapUser(user)),
			zap.String("language", lanague),
			zap.Error(err))
		tag = language.English
	}
	printer := message.NewPrinter(tag)
	ret := make([]string, 0, len(keys))
	for _, key := range keys {
		ret = append(ret, printer.Sprintf(key.Key, key.Args...))
	}
	return ret
}

func Retry(handler func() error, times int) error {
	var err error
	for i := 0; i < times; i++ {
		err = handler()
		if err == nil {
			return nil
		}
	}
	return err
}

func SendWithRetry(bot *botapi.BotAPI, msg botapi.Chattable) error {
	err := Retry(func() error {
		_, err := bot.Send(msg)
		return err
	}, retryTimes)
	if err != nil {
		err = errors.Wrap(err, "can't send message")
	}
	return err
}

func SendWithRetryAndResponse(bot *botapi.BotAPI, msg botapi.Chattable) (*botapi.Message, error) {
	var respMsg botapi.Message
	err := Retry(func() error {
		var err error
		respMsg, err = bot.Send(msg)
		return err
	}, retryTimes)
	if err != nil {
		err = errors.Wrap(err, "can't send message")
	}
	return &respMsg, err
}
