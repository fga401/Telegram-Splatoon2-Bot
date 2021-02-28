package test

import (
	"strconv"

	"go.uber.org/zap"
	"telegram-splatoon2-bot/common/log"
)

// IntToKey converts an int to a cache key.
func IntToKey(i int) []byte {
	return []byte(strconv.Itoa(i))
}

// IntToValue converts an int to a cache value.
func IntToValue(i int) []byte {
	return []byte(strconv.Itoa(i * 2))
}

// KeyToInt converts a cache key to a int.
func KeyToInt(key []byte) int {
	i, err := strconv.Atoi(string(key))
	if err != nil {
		log.Panic("can't convert key to int", zap.Error(err))
	}
	return i
}

// ValueToInt converts a cache value to a int.
func ValueToInt(key []byte) int {
	i, err := strconv.Atoi(string(key))
	if err != nil {
		log.Panic("can't convert value to int", zap.Error(err))
	}
	return i / 2
}
