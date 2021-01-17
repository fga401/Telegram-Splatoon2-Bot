package fastcache

import "github.com/spf13/viper"

type Config struct {
	// MaxBytes is the capacity of cache in bytes.
	// MaxBytes must be smaller than the available RAM size for the app, since the cache holds data in memory.
	// If MaxBytes is less than 32MB, then the minimum cache capacity is 32MB.
	MaxBytes int
}

func LoadFastCacheConfig(viper viper.Viper) Config {
	return Config{
		MaxBytes: viper.GetInt("cache.fastcache.maxBytes"),
	}
}

