package language

// Config sets up a language Service.
type Config struct {
	// SupportedLanguages are the IETF of all allowed language.
	SupportedLanguages []string
	// LocalePath is the path of translation JSON files.
	LocalePath string
}
