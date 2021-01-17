package fastcache

import (
	"github.com/VictoriaMetrics/fastcache"
	"telegram-splatoon2-bot/driver/cache"
	"telegram-splatoon2-bot/driver/cache/internal/cleaner"

	"time"
)

type impl struct {
	cache   *fastcache.Cache
	cleaner *cleaner.Cleaner
}

// New returns a Cache implemented by fastcache given capacity.
func New(config Config) cache.Cache {
	ret := &impl{
		cache: fastcache.New(config.MaxBytes),
	}
	ret.cleaner = cleaner.New(ret)
	return ret
}

// Set adds a key-value pair to cache.
func (impl *impl) Set(key, value []byte) {
	impl.cache.Set(key, value)
}

// SetExpiration adds a key-value pair to cache and deletes it after duration.
// Pair will be cleaned in time basing on a min heap and a resettable timer.
func (impl *impl) SetExpiration(key, value []byte, duration time.Duration) {
	if duration <= 0 {
		return
	}
	impl.cache.Set(key, value)
	impl.cleaner.Set(key, time.Now().Add(duration))
}

// Get returns the value against key.
func (impl *impl) Get(key []byte) []byte {
	ret := impl.cache.Get(nil, key)
	if len(ret) == 0 {
		return nil
	}
	return ret
}

// HasGet will return the value and true if key is in Cache. Otherwise it'll return nil and false.
func (impl *impl) HasGet(key []byte) ([]byte, bool) {
	return impl.cache.HasGet(nil, key)
}

// Has will return true if key is in Cache.
func (impl *impl) Has(key []byte) bool {
	return impl.cache.Has(key)
}

// Del deletes the key-value pair in Cache.
func (impl *impl) Del(key []byte) {
	impl.cache.Del(key)
}
