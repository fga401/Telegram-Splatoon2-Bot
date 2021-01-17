package gocache

import (
	"time"

	gocache "github.com/patrickmn/go-cache"
	"telegram-splatoon2-bot/driver/cache"
)

type impl struct {
	*gocache.Cache
}

// New returns a Cache implemented by go-cache given default expiration duration and cleanup interval.
func New(config Config) cache.Cache {
	return &impl{
		Cache: gocache.New(config.Expiration, config.CleanUp),
	}
}

// Set adds a key-value pair to cache.
func (impl *impl) Set(key, value []byte) {
	impl.Cache.Set(string(key), value, -1)
}

// SetExpiration adds a key-value pair to cache and deletes it after duration.
// Clean up will be executed in a periodical task, so it may not be cleaned in time
func (impl *impl) SetExpiration(key, value []byte, duration time.Duration) {
	if duration <= 0 {
		return
	}
	impl.Cache.Set(string(key), value, duration)
}

// Get returns the value against key.
func (impl *impl) Get(key []byte) []byte {
	ret, _ := impl.Cache.Get(string(key))
	return ret.([]byte)
}

// HasGet will return the value and true if key is in Cache. Otherwise it'll return nil and false.
func (impl *impl) HasGet(key []byte) ([]byte, bool) {
	ret, found := impl.Cache.Get(string(key))
	return ret.([]byte), found
}

// Has will return true if key is in Cache.
func (impl *impl) Has(key []byte) bool {
	_, found := impl.Cache.Get(string(key))
	return found
}

// Del deletes the key-value pair in Cache.
func (impl *impl) Del(key []byte) {
	impl.Cache.Delete(string(key))
}
