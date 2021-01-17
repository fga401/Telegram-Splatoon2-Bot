package language

import (
	"go.uber.org/zap"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"telegram-splatoon2-bot/common/log"
)

type Language string

func (l Language) IETF() string {
	return string(l)
}

type Service interface {
	Supported() []Language
	All() []Language
	Printer(language Language) *message.Printer
}

type impl struct {
	printers  map[Language]*message.Printer
	supported []Language
}

func NewService(config Config) Service {
	printers := make(map[Language]*message.Printer)
	supported := make([]Language, 0)
	ret := &impl{}
	for _, tag := range config.SupportedLanguages {
		for _, l := range ret.All() {
			if tag == l.IETF() {
				supported = append(supported, l)
				tag, err := language.Parse(tag)
				if err != nil {
					log.Panic("can't parse IETF language tag", zap.Error(err))
				}
				printers[l] = message.NewPrinter(tag)
			}
		}
	}
	ret.printers = printers
	ret.supported = supported
	return ret
}

func (svc *impl) Supported() []Language {
	return svc.supported
}

func (svc *impl) All() []Language {
	return all
}

func (svc *impl) Printer(language Language) *message.Printer {
	return svc.printers[language]
}
