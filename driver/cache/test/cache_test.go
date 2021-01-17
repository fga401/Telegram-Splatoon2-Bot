package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"telegram-splatoon2-bot/driver/cache"
	"telegram-splatoon2-bot/driver/cache/fastcache"
	"telegram-splatoon2-bot/driver/cache/gocache"
	"telegram-splatoon2-bot/driver/cache/syncmap"
)

func TestAllCache(t *testing.T) {
	testCache(t, syncmap.New())
	fmt.Println("SyncMap passed.")
	testCache(t, fastcache.New(fastcache.Config{
		MaxBytes: 1 << 20,
	}))
	fmt.Println("FastCache passed.")
	testCache(t, gocache.New(gocache.Config{
		Expiration: 0,
		CleanUp:    time.Second,
	}))
	fmt.Println("GoCache passed.")
}

func testCache(t *testing.T, cache cache.Cache) {
	const N = 10000
	for i := 0; i < N; i++ {
		cache.Set(IntToKey(i), IntToValue(i))
	}
	for i := 0; i < N; i += 2 {
		key := IntToKey(i)
		cache.Del(key)
		require.False(t, cache.Has(key), "Key should be deleted.")
	}
	for i := 1; i < N; i += 2 {
		key := IntToKey(i)
		expected := IntToValue(i)
		require.True(t, cache.Has(key), "Key should be found.")
		actual := cache.Get(key)
		require.Equal(t, expected, actual, "Value should be equal.")
	}
	const M = 30
	for i := M; i >= 0; i-- {
		cache.SetExpiration(IntToKey(i), IntToValue(i), time.Second*time.Duration(i))
	}
	<-time.After(time.Second * M / 2)
	for i := M; i >= 0; i-- {
		key := IntToKey(i)
		if i > M/2 {
			require.True(t, cache.Has(key), "Half of key should be existed.")
		} else {
			require.False(t, cache.Has(key), "Half of key should be expired.")
		}
	}
}
