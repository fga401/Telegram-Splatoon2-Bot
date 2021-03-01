package language

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type Language string

func (l Language) IETF() string {
	return string(l)
}

func (l Language) Tag() language.Tag {
	return language.MustParse(string(l))
}

type Service interface {
	Supported() []Language
	All() []Language
	Printer(language Language) *message.Printer
}

type impl struct {
	printers   map[Language]*message.Printer
	supported  []Language
	localePath string
}

func NewService(config Config) Service {
	printers := make(map[Language]*message.Printer)
	supported := make([]Language, 0)
	ret := &impl{
		localePath: config.LocalePath,
	}
	supportedSet := make(map[string]struct{})
	for _, tag := range config.SupportedLanguages {
		supportedSet[tag] = struct{}{}
	}
	for _, l := range ret.All() {
		if _, ok := supportedSet[l.IETF()]; ok {
			supported = append(supported, l)
		}
	}
	ret.supported = supported
	cat := ret.loadCatalog()
	opt := message.Catalog(cat)
	for _, l := range ret.All() {
		if _, ok := supportedSet[l.IETF()]; ok {
			printers[l] = message.NewPrinter(l.Tag(), opt)
		}
	}
	ret.printers = printers
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
