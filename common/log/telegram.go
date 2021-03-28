package log

import (
	botApi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap/zapcore"
)

// UpdateLogger wraps Update as zapcore.ObjectMarshaler.
type UpdateLogger botApi.Update

// MessageLogger wraps Message as zapcore.ObjectMarshaler.
type MessageLogger botApi.Message

// CallbackQueryLogger wraps CallbackQuery as zapcore.ObjectMarshaler.
type CallbackQueryLogger botApi.CallbackQuery

// UserLogger wraps User as zapcore.ObjectMarshaler.
type UserLogger botApi.User

// ChatLogger wraps Chat as zapcore.ObjectMarshaler.
type ChatLogger botApi.Chat

// MarshalLogObject encodes UpdateLogger for logging.
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

// MarshalLogObject encodes MessageLogger for logging.
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

// MarshalLogObject encodes UserLogger for logging.
func (l UserLogger) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddInt("id", l.ID)
	encoder.AddString("username", l.UserName)
	encoder.AddString("language_code", l.LanguageCode)
	return nil
}

// MarshalLogObject encodes ChatLogger for logging.
func (l ChatLogger) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddInt64("id", l.ID)
	encoder.AddString("type", l.Type)
	encoder.AddString("title", l.Title)
	return nil
}

// MarshalLogObject encodes CallbackQueryLogger for logging.
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

// UserPtrLogger convert a pointer of User to UserLogger.
func UserPtrLogger(ptr *botApi.User) UserLogger {
	return UserLogger(*ptr)
}
