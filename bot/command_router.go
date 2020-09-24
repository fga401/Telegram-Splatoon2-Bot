package bot

import (
	botapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.uber.org/zap"
	"strings"
	log "telegram-splatoon2-bot/logger"
)

type Handler func(message *botapi.Update, bot *botapi.BotAPI)

type CommandRouter struct {
	handlers             map[string]Handler
	caseSensitiveCommand map[string]struct{}
	defaultHandler       Handler
}

func NewCommandRouter() *CommandRouter {
	return &CommandRouter{
		handlers:             make(map[string]Handler),
		caseSensitiveCommand: make(map[string]struct{}),
	}
}

func (cr *CommandRouter) Add(cmd string, handler Handler) {
	lowerCmd := strings.ToLower(cmd)
	cr.handlers[lowerCmd] = handler
}

func (cr *CommandRouter) AddDefault(handler Handler) {
	cr.defaultHandler = handler
}

func (cr *CommandRouter) AddWithCaseSensitivity(cmd string, handler Handler, caseSensitive bool) {
	if caseSensitive {
		cr.caseSensitiveCommand[cmd] = struct{}{}
	} else {
		cmd = strings.ToLower(cmd)
	}
	cr.handlers[cmd] = handler
}

func (cr *CommandRouter) Run(update *botapi.Update, bot *botapi.BotAPI) bool {
	cmd := update.Message.Command()
	if _, in := cr.caseSensitiveCommand[cmd]; !in {
		cmd = strings.ToLower(cmd)
	}
	handler, in := cr.handlers[cmd]
	if !in {
		log.Info("command not existed", zap.String("Command", cmd))
		if cr.defaultHandler == nil {
			return false
		}
		handler = cr.defaultHandler
	}
	go handler(update, bot)
	return true
}
