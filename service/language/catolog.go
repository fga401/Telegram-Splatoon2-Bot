package language

import (
	"os"

	"github.com/mailru/easyjson"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/text/feature/plural"
	"golang.org/x/text/language"
	"golang.org/x/text/message/catalog"
	"telegram-splatoon2-bot/common/log"
)

func packagePath(base string, tag string) string {
	return base + "/" + tag + ".json"
}

func (svc *impl) loadPackage(lang Language) (langPackage, error) {
	path := packagePath(svc.localePath, lang.IETF())
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "can't open file: "+path)
	}
	ret := langPackage{}
	err = easyjson.UnmarshalFromReader(file, &ret)
	if err != nil {
		return nil, errors.Wrap(err, "can't unmarshal langPackage")
	}
	return ret, nil
}

func (svc *impl) loadCatalog() catalog.Catalog {
	fallbackTag := English.Tag()
	builder := catalog.NewBuilder(catalog.Fallback(fallbackTag))
	langs := make(map[string]langPackage)
	for _, lang := range svc.Supported() {
		pkg, err := svc.loadPackage(lang)
		if err != nil {
			log.Warn("can't load language package", zap.String("tag", lang.IETF()), zap.Error(err))
			continue
		}
		langs[lang.IETF()] = pkg
	}
	for t, msgs := range langs {
		tag := language.MustParse(t)
		for _, msg := range msgs {
			var varArgs []catalog.Message
			if len(msg.Vars) > 0 {
				for _, v := range msg.Vars {
					var pluralArgs []interface{}
					for _, c := range v.Plural.Cases {
						pluralArgs = append(pluralArgs, c.Cond, c.Text)
					}
					pluralMsg := plural.Selectf(v.Plural.Arg, v.Plural.Format, pluralArgs...)
					varMsg := catalog.Var(v.Key, pluralMsg)
					varArgs = append(varArgs, varMsg)
				}
			}
			varArgs = append(varArgs, catalog.String(msg.Text))
			err := builder.Set(tag, msg.Key, varArgs...)
			if err != nil {
				log.Panic("can't load language", zap.String("tag", t), zap.String("key", msg.Key), zap.Error(err))
			}
		}
	}
	return builder
}
