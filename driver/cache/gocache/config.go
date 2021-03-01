package gocache

import (
	"time"
)

// Config sets up a Cache implemented by gocache.
type Config struct {
	// Expiration is default expiration time.
	// If the expiration duration is less than one,
	// the items in the cache never expire by default.
	Expiration time.Duration
	// CleanUp is the interval between two automatic cleanup
	// If the cleanup interval is less than one, expired items are not deleted from the cache.
	CleanUp time.Duration
}
