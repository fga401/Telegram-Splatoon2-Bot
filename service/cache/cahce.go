package cache

import (
	"github.com/VictoriaMetrics/fastcache"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	gocache "github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"telegram-splatoon2-bot/service/db"
)

var (
	Cache *CacheImpl
)

//type CacheInterface interface {
//	SetProofKey(user *tgbotapi.User, proofKey []byte) error
//	GetProofKey(user *tgbotapi.User) ([]byte, error)
//	SetRuntime(user *tgbotapi.User, runtime *db.Runtime) error
//	GetRuntime(user *tgbotapi.User) (*db.Runtime, error)
//}

type CacheImpl struct {
	fastCache *fastcache.Cache
	gocache   *gocache.Cache
}

func InitCache() {
	Cache = &CacheImpl{
		fastCache: fastcache.New(
			viper.GetInt("cache.fastcache.maxBytes")),
		gocache: gocache.New(
			viper.GetDuration("cache.gocache.expiration"),
			viper.GetDuration("cache.gocache.cleanUp")),
	}
}

func (c *CacheImpl) SetProofKey(user *tgbotapi.User, proofKey []byte) error {
	if user == nil {
		return errors.Errorf("user is nil")
	}
	key := userToStringKey(user)
	c.gocache.SetDefault(key, proofKey)
	return nil
}

// GetProofKey returns proof key if existed, or nil if not found
func (c *CacheImpl) GetProofKey(user *tgbotapi.User) ([]byte, error) {
	if user == nil {
		return nil, errors.Errorf("user is nil")
	}
	key := userToStringKey(user)
	proofKeyInterface, found := c.gocache.Get(key)
	if !found {
		return nil, nil
	}
	proofKey := proofKeyInterface.([]byte)
	return proofKey, nil
}

func (c *CacheImpl) SetRuntime(user *tgbotapi.User, runtime *db.Runtime) error {
	if user == nil {
		return errors.Errorf("user is nil")
	}
	key, err := userToBytesKey(user)
	if err != nil {
		return errors.Wrap(err, "can't convert user to key")
	}
	value, err := serializeRuntime(runtime)
	if err != nil {
		return errors.Wrap(err, "can't convert runtime  to value")
	}
	c.fastCache.Set(key, value)
	return nil
}

// GetRuntime returns proof key if existed, or nil if not found
func (c *CacheImpl) GetRuntime(user *tgbotapi.User) (*db.Runtime, error) {
	if user == nil {
		return nil, errors.Errorf("user is nil")
	}
	key, err := userToBytesKey(user)
	if err != nil {
		return nil, errors.Wrap(err, "can't convert user to key")
	}
	value, found := c.fastCache.HasGet(nil, key)
	if !found {
		return nil, nil
	}
	runtime, err := deserializeRuntime(value)
	if err != nil {
		return nil, errors.Wrap(err, "can't convert value to runtime ")
	}
	return runtime, nil
}

func (c *CacheImpl) DeleteRuntime(user *tgbotapi.User) {
	if user != nil {
		key, err := userToBytesKey(user)
		if err != nil {
			c.fastCache.Del(key)
		}
	}
}
