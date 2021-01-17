package log

import (
	botApi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap/zapcore"
)

type UpdateLogger botApi.Update
type MessageLogger botApi.Message
type CallbackQueryLogger botApi.CallbackQuery
type UserLogger botApi.User
type ChatLogger botApi.Chat

func (l UpdateLogger) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddInt("id", l.UpdateID)
	if l.Message != nil {
		_ = encoder.AddObject("message", (*MessageLogger)(l.Message))
	}
	if l.CallbackQuery != nil {
		_ = encoder.AddObject("callback_query", (*CallbackQueryLogger)(l.CallbackQuery))
	}
	return nil
}

func (l *MessageLogger) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddInt("id", l.MessageID)
	if l.From != nil {
		_ = encoder.AddObject("from", UserLogger(*l.From))
	}
	encoder.AddInt("date", l.Date)
	if l.From != nil {
		_ = encoder.AddObject("chat", (*ChatLogger)(l.Chat))
	}
	encoder.AddString("text", l.Text)
	return nil
}

func (l UserLogger) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddInt("id", l.ID)
	encoder.AddString("username", l.UserName)
	encoder.AddString("language_code", l.LanguageCode)
	return nil
}

func (l ChatLogger) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddInt64("id", l.ID)
	encoder.AddString("type", l.Type)
	encoder.AddString("title", l.Title)
	return nil
}

func (l CallbackQueryLogger) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddString("id", l.ID)
	if l.From != nil {
		_ = encoder.AddObject("from", UserLogger(*l.From))
	}
	if l.Message != nil {
		_ = encoder.AddObject("message", (*MessageLogger)(l.Message))
	}
	if l.InlineMessageID != "" {
		encoder.AddString("inline_message_id", l.InlineMessageID)
	}
	encoder.AddString("chat_instance", l.ChatInstance)
	encoder.AddString("data", l.Data)
	return nil
}

func UserPtrLogger(ptr *botApi.User) UserLogger {
	return UserLogger(*ptr)
}

