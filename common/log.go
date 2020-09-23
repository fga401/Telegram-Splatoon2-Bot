package common

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"strings"
	"time"
)

var Logger *zap.Logger

func init() {
	// todo
	level := strings.ToLower(viper.GetString("log.lever"))
	cfg := zap.NewProductionConfig()
	switch level {
	case "debug":
		cfg.Level.SetLevel(zap.DebugLevel)
	case "info":
		cfg.Level.SetLevel(zap.InfoLevel)
	}
	logger ,err := cfg.Build()
	if err != nil {
		log.Fatal(errors.Wrap(err, "can't initialize zap logger"))
	}
	Logger = logger
}

type UpdateWrapper tgbotapi.Update
type MessageWrapper tgbotapi.Message

func (u *UpdateWrapper) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt("id", u.UpdateID)
	return nil
}

func (m *MessageWrapper) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt("id", m.MessageID)
	enc.AddString("text", m.Text)
	enc.AddString("fro,", m.From.UserName)
	enc.AddTime("date", time.Unix(int64(m.Date), 0))
	return nil
}

func WrapMessage(m *tgbotapi.Message) *MessageWrapper{
	return (*MessageWrapper)(m)
}
