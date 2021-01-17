package test

import (
	"strconv"

	"go.uber.org/zap"
	"telegram-splatoon2-bot/common/log"
)

func IntToKey(i int) []byte {
	return []byte(strconv.Itoa(i))
}

func IntToValue(i int) []byte {
	return []byte(strconv.Itoa(i * 2))
}

func KeyToInt(key []byte) int {
	i, err := strconv.Atoi(string(key))
	if err != nil {
		log.Panic("can't convert key to int", zap.Error(err))
	}
	return i
}

func ValueToInt(key []byte) int {
	i, err := strconv.Atoi(string(key))
	if err != nil {
		log.Panic("can't convert value to int", zap.Error(err))
	}
	return i / 2
}
