package cache

import (
	"github.com/dgraph-io/ristretto"
	gocache "github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var (
	GoCache   *gocache.Cache
	Ristretto *ristretto.Cache
)

func InitCache() {
	var err error
	numCounters := viper.GetInt64("cache.ristretto.numCounters")
	maxCost := viper.GetInt64("cache.ristretto.maxCost")
	metrics := viper.GetBool("cache.ristretto.metrics")
	Ristretto, err = ristretto.NewCache(&ristretto.Config{
		NumCounters: numCounters,
		MaxCost:     maxCost,
		BufferItems: 64,
		Metrics:     metrics,
		OnEvict:     nil,
		KeyToHash:   nil,
		Cost:        nil,
	})
	if err != nil {
		panic(errors.Wrap(err, "can't init ristretto"))
	}

	expireTime := viper.GetDuration("cache.goCache.expire")
	purgeTime := viper.GetDuration("cache.goCache.purge")
	GoCache = gocache.New(expireTime, purgeTime)
}
