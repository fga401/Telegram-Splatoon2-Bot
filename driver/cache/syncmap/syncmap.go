package syncmap

import (
	"sync"
	"time"

	"telegram-splatoon2-bot/driver/cache"
	"telegram-splatoon2-bot/driver/cache/internal/cleaner"
)

type impl struct {
	syncMap sync.Map
	cleaner *cleaner.Cleaner
}

func (impl *impl) Set(key, value []byte) {
	impl.syncMap.Store(string(key), value)
}

func (impl *impl) SetExpiration(key, value []byte, duration time.Duration) {
	impl.syncMap.Store(string(key), value)
	impl.cleaner.Set(key, time.Now().Add(duration))
}

func (impl *impl) Get(key []byte) []byte {
	value, _ := impl.syncMap.Load(string(key))
	return value.([]byte)
}

func (impl *impl) HasGet(key []byte) ([]byte, bool) {
	value, ok := impl.syncMap.Load(string(key))
	return value.([]byte), ok
}

func (impl *impl) Has(key []byte) bool {
	_, ok := impl.syncMap.Load(string(key))
	return ok
}

func (impl *impl) Del(key []byte) {
	impl.syncMap.Delete(string(key))
}

func (impl *impl) Range(f func([]byte, []byte) bool) {
	impl.syncMap.Range(func(k, v interface{}) bool {
		return f([]byte(k.(string)), v.([]byte))
	})
}

// New returns a new IterableCache implemented by sync.map.
func New() cache.IterableCache {
	ret := &impl{}
	ret.cleaner = cleaner.New(ret)
	return ret
}
