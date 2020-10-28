package botutil

import (
	botapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap"
	"strings"
	log "telegram-splatoon2-bot/logger"
)

type Handler func(message *botapi.Update) error

type DescribedHandler struct {
	handler Handler
	des     string
}

type UpdateRouter struct {
	commandHandlers       map[string]*DescribedHandler
	caseSensitiveCommand  map[string]struct{}
	callbackQueryHandlers map[string]*DescribedHandler
	defaultCommandHandler *DescribedHandler
	textHandler           *DescribedHandler
}

func NewUpdateRouter() *UpdateRouter {
	return &UpdateRouter{
		commandHandlers:       make(map[string]*DescribedHandler),
		caseSensitiveCommand:  make(map[string]struct{}),
		callbackQueryHandlers: make(map[string]*DescribedHandler),
	}
}

func (ur *UpdateRouter) AddCommandHandler(cmd string, handler Handler, des string) {
	lowerCmd := strings.ToLower(cmd)
	ur.commandHandlers[lowerCmd] = &DescribedHandler{handler, des}
}

func (ur *UpdateRouter) SetDefaultCommandHandler(handler Handler, des string) {
	ur.defaultCommandHandler = &DescribedHandler{handler, des}
}

func (ur *UpdateRouter) AddCommandHandlerWithCaseSensitivity(cmd string, handler Handler, caseSensitive bool, des string) {
	if caseSensitive {
		ur.caseSensitiveCommand[cmd] = struct{}{}
	} else {
		cmd = strings.ToLower(cmd)
	}
	ur.commandHandlers[cmd] = &DescribedHandler{handler, des}
}

func (ur *UpdateRouter) SetTextHandler(handler Handler, des string) {
	ur.textHandler = &DescribedHandler{handler, des}
}

func (ur *UpdateRouter) AddCallbackQueryHandler(prefix string, handler Handler, des string) {
	ur.callbackQueryHandlers[prefix] = &DescribedHandler{handler, des}
}

func (ur *UpdateRouter) Route(update *botapi.Update) *DescribedHandler {
	if update.Message != nil {
		msg := update.Message
		if !msg.IsCommand() {
			return ur.textHandler
		}

		cmd := msg.Command()
		if _, in := ur.caseSensitiveCommand[cmd]; !in {
			cmd = strings.ToLower(cmd)
		}
		handler, in := ur.commandHandlers[cmd]
		if !in {
			log.Info("command not existed", zap.String("Command", cmd))
			if ur.defaultCommandHandler == nil {
				return nil
			}
			handler = ur.defaultCommandHandler
		}
		return handler
	}
	if update.CallbackQuery != nil {
		callbackQuery := update.CallbackQuery
		prefix := CallbackHelper.GetPrefix(callbackQuery.Data)
		handler, in := ur.callbackQueryHandlers[prefix]
		if !in {
			log.Warn("prefix not existed", zap.String("prefix", prefix))
			return nil
		}
		return handler
	}
	log.Debug("unsupported update", zap.Object("update", log.WrapUpdate(update)))
	return nil
}

type callbackHelper struct {}
var CallbackHelper callbackHelper
func (callbackHelper)SetPrefix(prefix, text string) string {
	return prefix + ":" + text
}

func (callbackHelper)GetPrefix(data string) string {
	index := strings.Index(data, ":")
	if index == -1 {
		return ""
	}
	return data[:index]
}

func (callbackHelper)GetText(data string) string {
	index := strings.Index(data, ":")
	if index == -1 {
		return data
	}
	return data[index+1:]
}
