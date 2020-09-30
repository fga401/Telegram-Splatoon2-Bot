package bot

import (
	botapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap"
	"strings"
	log "telegram-splatoon2-bot/logger"
)

type Handler func(message *botapi.Update, bot *botapi.BotAPI) error

type DescribedHandler struct {
	handler Handler
	des string
}

type UpdateRouter struct {
	commandHandlers       map[string]*DescribedHandler
	caseSensitiveCommand  map[string]struct{}
	defaultCommandHandler *DescribedHandler
	textHandler           *DescribedHandler
}

func NewUpdateRouter() *UpdateRouter {
	return &UpdateRouter{
		commandHandlers:      make(map[string]*DescribedHandler),
		caseSensitiveCommand: make(map[string]struct{}),
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

func (ur *UpdateRouter) Route(update *botapi.Update) *DescribedHandler {
	if update.Message == nil {
		return nil
	}

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
