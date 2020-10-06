package logger

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap/zapcore"
	"time"
)

type UpdateWrapper tgbotapi.Update
type MessageWrapper tgbotapi.Message
type UserWrapper tgbotapi.User
type CallbackQueryWrapper tgbotapi.CallbackQuery

func (u *UpdateWrapper) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt("id", u.UpdateID)
	if u.Message!= nil {
		_ = enc.AddObject("message",WrapMessage(u.Message))
	}
	if u.CallbackQuery!= nil {
		_ = enc.AddObject("callback_query",WrapCallbackQuery(u.CallbackQuery))
	}
	//enc.AddBool("message", u.Message != nil)
	//enc.AddBool("edited_message", u.EditedMessage != nil)
	//enc.AddBool("channel_post", u.ChannelPost != nil)
	//enc.AddBool("edited_channel_post", u.EditedChannelPost != nil)
	//enc.AddBool("inline_query", u.InlineQuery != nil)
	//enc.AddBool("chosen_inline_result", u.ChosenInlineResult != nil)
	//enc.AddBool("callback_query", u.CallbackQuery != nil)
	//enc.AddBool("shipping_query", u.ShippingQuery != nil)
	//enc.AddBool("pre_checkout_query", u.PreCheckoutQuery != nil)
	//enc.AddBool("poll", u.Poll != nil)
	//enc.AddBool("poll_answer", u.PollAnswer != nil)
	return nil
}
func WrapUpdate(u *tgbotapi.Update) *UpdateWrapper {
	return (*UpdateWrapper)(u)
}

func (m *MessageWrapper) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt("id", m.MessageID)
	enc.AddString("text", m.Text)
	enc.AddString("from,", m.From.UserName)
	enc.AddTime("date", time.Unix(int64(m.Date), 0))
	return nil
}
func WrapMessage(m *tgbotapi.Message) *MessageWrapper {
	return (*MessageWrapper)(m)
}

func (u *UserWrapper) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt("id", u.ID)
	enc.AddString("user_name", u.UserName)
	return nil
}
func WrapUser(u *tgbotapi.User) *UserWrapper {
	return (*UserWrapper)(u)
}

func (q *CallbackQueryWrapper) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("id", q.ID)
	enc.AddString("data", q.Data)
	return enc.AddObject("message",WrapMessage(q.Message))
}
func WrapCallbackQuery(u *tgbotapi.CallbackQuery) *CallbackQueryWrapper {
	return (*CallbackQueryWrapper)(u)
}