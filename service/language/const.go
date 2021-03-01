package language

var all = []Language{
	English,
	Japanese,
	SimplifiedChinese,
	TraditionalChinese,
}

// all available timezones
const (
	English            = Language("en")
	Japanese           = Language("ja")
	SimplifiedChinese  = Language("zh-CN")
	TraditionalChinese = Language("zh-TW")
)

// ByIETF returns a language given the IETF
func ByIETF(ietf string) Language {
	switch ietf {
	case "en":
		return English
	case "ja":
		return Japanese
	case "zh-CN":
		return SimplifiedChinese
	case "zh-TW":
		return TraditionalChinese
	default:
		return English
	}
}
