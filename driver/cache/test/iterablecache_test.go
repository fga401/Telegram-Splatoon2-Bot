package test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"telegram-splatoon2-bot/driver/cache"
	"telegram-splatoon2-bot/driver/cache/syncmap"
)

func TestAllIterableCache(t *testing.T) {
	testIterableCache(t, syncmap.New())
}

func testIterableCache(t *testing.T, cache cache.IterableCache) {
	const N = 10000
	for i := 0; i < N; i++ {
		cache.Set(IntToKey(i), IntToValue(i))
	}
	for i := 0; i < N; i += 2 {
		key := IntToKey(i)
		cache.Del(key)
		require.False(t, cache.Has(key), "Key should be deleted.")
	}
	keys := getKeys(cache)
	set := make(map[string]struct{})
	require.Equal(t, N/2, len(keys), "Half of keys should be retrieved.")
	for _, key := range keys{
		set[string(key)] = struct{}{}
		i := KeyToInt(key)
		expected := IntToValue(i)
		require.True(t, cache.Has(key), "Key should be found.")
		actual := cache.Get(key)
		require.Equal(t, expected, actual, "Value should be equal.")
	}
	require.Equal(t, len(keys), len(set), "All keys should be different.")
}

func getKeys(cache cache.IterableCache) [][]byte {
	ret := make([][]byte, 0)
	cache.Range(func(key []byte, value []byte) bool {
		ret = append(ret, key)
		return true
	})
	return ret
}
