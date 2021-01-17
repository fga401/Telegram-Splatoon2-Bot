package cache

import (
	"time"
)

// Cache provides a in-memory cache for other service.
type Cache interface {
	// Set adds a key-value pair to cache.
	Set(key, value []byte)
	// SetExpiration adds a key-value pair to cache and deletes it after duration.
	// Clean up strategies of each Cache is different.
	// If duration <= 0, do nothing
	SetExpiration(key, value []byte, duration time.Duration)
	// Get returns the value against key.
	// If not found, return nil
	Get(key []byte) []byte
	// HasGet will return the value and true if key is in Cache. Otherwise it'll return nil and false.
	HasGet(key []byte) ([]byte, bool)
	// Has will return true if key is in Cache.
	Has(key []byte) bool
	// Del deletes the key-value pair in Cache.
	Del(key []byte)
}

// IterableCache provides an iterable in-memory cache for other service.
type IterableCache interface {
	Cache
	// Range calls f sequentially for each key and value present in the map.
	// If f returns false, range stops the iteration.
	Range(f func([]byte, []byte) bool)
}
