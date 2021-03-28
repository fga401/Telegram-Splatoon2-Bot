package router

import (
	"net/http"
	"regexp"
	"strings"

	botApi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	log "telegram-splatoon2-bot/common/log"
)

type regexpHandler struct {
	re      *regexp.Regexp
	handler Handler
}

type impl struct {
	commandHandlers       map[string]Handler
	defaultCommandHandler Handler
	callbackQueryHandlers map[string]Handler
	regexpCommandHandlers []regexpHandler
	textHandler           Handler
	config                Config
	bot                   *botApi.BotAPI
}

// New returns a Router object.
func New(bot *botApi.BotAPI, config Config) Router {
	return &impl{
		commandHandlers:       make(map[string]Handler),
		callbackQueryHandlers: make(map[string]Handler),
		config:                config,
		bot:                   bot,
	}
}

func (r *impl) Run() {
	updatesChan, err := r.getUpdateChan()
	if err != nil {
		log.Panic("can't run router", zap.Error(err))
	}
	r.dispatch(updatesChan)
}

func (r *impl) getUpdateChan() (botApi.UpdatesChannel, error) {
	switch r.config.Mode {
	case ModeEnum.Polling:
		updatesChan, err := r.bot.GetUpdatesChan(r.config.Polling)
		if err != nil {
			return nil, errors.Wrap(err, "can't get bot update channel by polling")
		}
		return updatesChan, nil
	case ModeEnum.WebHook:
		if !strings.HasSuffix(r.config.WebHook.URL, r.bot.Token) {
			r.config.WebHook.URL += r.bot.Token
		}
		var config botApi.WebhookConfig
		if r.config.WebHook.Cert == "" {
			config = botApi.NewWebhook(r.config.WebHook.URL)
			go func() {
				err := http.ListenAndServe("0.0.0.0:"+r.config.WebHook.Port, nil)
				if err != nil {
					log.Panic("can't start server", zap.Error(err))
				}
			}()
		} else {
			config = botApi.NewWebhookWithCert(r.config.WebHook.URL, r.config.WebHook.Cert)
			go func() {
				err := http.ListenAndServeTLS("0.0.0.0:"+r.config.WebHook.Port, r.config.WebHook.Cert, r.config.WebHook.Key, nil)
				if err != nil {
					log.Panic("can't start server", zap.Error(err))
				}
			}()
		}
		_, err := r.bot.SetWebhook(config)
		if err != nil {
			return nil, errors.Wrap(err, "can't set webhook")
		}
		updatesChan := r.bot.ListenForWebhook("/" + r.bot.Token)
		return updatesChan, nil
	default:
		return nil, errors.New("unknown mode")
	}
}

func (r *impl) dispatch(updatesChan botApi.UpdatesChannel) {
	if r.config.MaxWorker <= 0 {
		for update := range updatesChan {
			go r.routine(update)
		}
	} else {
		worker := make(chan botApi.Update, r.config.MaxWorker)
		for i := int32(0); i < r.config.MaxWorker; i++ {
			go func() {
				for update := range worker {
					r.routine(update)
				}
			}()
		}
		for update := range updatesChan {
			worker <- update
		}
	}
}

func (r *impl) routine(update botApi.Update) {
	handler := r.route(update)
	if handler != nil {
		err := handler(update)
		if err != nil {
			log.Error("can't handle update",
				zap.Object("update", log.UpdateLogger(update)),
				zap.Error(err),
			)
		} else {
			log.Info("handle update successfully",
				zap.Object("update", log.UpdateLogger(update)),
			)
		}
	}
}
