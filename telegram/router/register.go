package router

import (
	"regexp"
	"strings"

	botApi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"telegram-splatoon2-bot/common/enum"
	log "telegram-splatoon2-bot/common/log"
	callbackQueryUtil "telegram-splatoon2-bot/telegram/callbackquery"
)

type updateType enum.Enum
type updateTypeEnumStruct struct {
	Text            updateType
	Command         updateType
	CallbackQuery   updateType
	UnsupportedType updateType
}

var updateTypeEnum = enum.Assign(&updateTypeEnumStruct{}).(*updateTypeEnumStruct)

func (r *impl) RegisterCommand(command string, handler Handler, options ...Option) {
	if hasOption(options, OptionEnum.AsDefault) {
		if err := r.registerDefaultCommand(handler); err != nil {
			log.Panic(`can't register "`+command+`" command handler`, zap.Error(err))
		}
	}
	if hasOption(options, OptionEnum.Regexp) {
		if err := r.registerRegexpCommand(command, handler); err != nil {
			log.Panic(`can't register "`+command+`" regular expression command handler`, zap.Error(err))
		}
	}
	if !hasOption(options, OptionEnum.CaseSensitive) {
		command = strings.ToLower(command)
	}
	if err := r.registerCommand(command, handler); err != nil {
		log.Panic(`can't register "`+command+`" command handler`, zap.Error(err))
	}
	return
}

func (r *impl) RegisterCallbackQuery(prefix string, handler Handler) {
	if err := r.registerCallbackQuery(prefix, handler); err != nil {
		log.Panic(`can't register "`+prefix+`" callback query handler`, zap.Error(err))
	}
	return
}

func (r *impl) RegisterText(handler Handler) {
	if r.textHandler != nil {
		log.Panic("text handler has already been registered")
	}
	r.textHandler = handler
}

func (r *impl) registerRegexpCommand(command string, handler Handler) error {
	if !strings.HasPrefix(command, "^") {
		command = "^" + command
	}
	if !strings.HasSuffix(command, "$") {
		command = command + "$"
	}
	re, err := regexp.Compile(command)
	if err != nil {
		return errors.Wrap(err, "can't compile regular expression")
	}
	r.regexpCommandHandlers = append(r.regexpCommandHandlers, regexpHandler{
		re:      re,
		handler: handler,
	})
	return nil
}

func (r *impl) registerCommand(command string, handler Handler) error {
	if _, found := r.commandHandlers[command]; found {
		return errors.New("command has already been registered")
	}
	r.commandHandlers[command] = handler
	return nil
}

func (r *impl) registerDefaultCommand(handler Handler) error {
	if r.defaultCommandHandler != nil {
		return errors.New("default command handler has already been registered")
	}
	r.defaultCommandHandler = handler
	return nil
}

func (r *impl) registerCallbackQuery(prefix string, handler Handler) error {
	if _, found := r.callbackQueryHandlers[prefix]; found {
		return errors.New("prefix has already been registered")
	}
	r.callbackQueryHandlers[prefix] = handler
	return nil
}

func hasOption(options []Option, target Option) bool {
	for _, option := range options {
		if option == target {
			return true
		}
	}
	return false
}

func getUpdateType(update botApi.Update) updateType {
	if update.Message != nil {
		if update.Message.IsCommand() {
			return updateTypeEnum.Command
		}
		return updateTypeEnum.Text
	}
	if update.CallbackQuery != nil {
		return updateTypeEnum.CallbackQuery
	}
	return updateTypeEnum.UnsupportedType
}

func (r *impl) route(update botApi.Update) Handler {
	switch getUpdateType(update) {
	case updateTypeEnum.Text:
		if r.textHandler == nil {
			log.Warn("no text handler",
				zap.Object("update", log.UpdateLogger(update)),
			)
		}
		return r.textHandler
	case updateTypeEnum.Command:
		cmd := update.Message.Command()
		if _, in := r.commandHandlers[cmd]; !in {
			cmd = strings.ToLower(cmd)
		}
		handler := r.commandHandlers[cmd]
		if handler == nil {
			for _, reHandler := range r.regexpCommandHandlers {
				if reHandler.re.Match([]byte(cmd)) {
					handler = reHandler.handler
					break
				}
			}
		}
		if handler == nil {
			handler = r.defaultCommandHandler
		}
		if handler == nil {
			log.Warn("command not existed",
				zap.String("command", cmd),
				zap.Object("update", log.UpdateLogger(update)),
			)
		}
		return handler
	case updateTypeEnum.CallbackQuery:
		callbackQuery := update.CallbackQuery
		prefix := callbackQueryUtil.GetPrefix(callbackQuery.Data)
		handler := r.callbackQueryHandlers[prefix]
		if handler == nil {
			log.Warn("callback query prefix not existed",
				zap.String("prefix", prefix),
				zap.Object("update", log.UpdateLogger(update)),
			)
		}
		return handler
	default:
		return nil
	}
}
