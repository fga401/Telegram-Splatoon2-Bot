package language

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// Language of User.
type Language string

// IETF returns IETF tag.
func (l Language) IETF() string {
	return string(l)
}

// Tag returns language.Tag.
func (l Language) Tag() language.Tag {
	return language.MustParse(string(l))
}

// Service manages language and translation.
type Service interface {
	// Supported returns all supported language
	Supported() []Language
	// All returns all defined language.
	All() []Language
	// Printer returns a message.Printer against the language.
	Printer(language Language) *message.Printer
}

type impl struct {
	printers   map[Language]*message.Printer
	supported  []Language
	localePath string
}

// NewService returns a new Service.
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
