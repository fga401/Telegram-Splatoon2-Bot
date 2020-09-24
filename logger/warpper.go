package logger

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap/zapcore"
	"time"
)

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

func WrapMessage(m *tgbotapi.Message) *MessageWrapper {
	return (*MessageWrapper)(m)
}

