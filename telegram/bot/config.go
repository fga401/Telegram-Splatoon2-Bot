package bot

// CallbackQueryConfig sets up an AnswerCallbackQuery request.
// More info: https://core.telegram.org/bots/api#answercallbackquery
type CallbackQueryConfig struct {
	Text      string
	ShowAlert bool
	CacheTime int
}

// Config sets up a Bot
type Config struct {
	// RetryTimes identifies the times of retry after failure.
	RetryTimes int
	// DefaultCallbackQueryConfig is the default config using for all AnswerCallbackQuery requests.
	DefaultCallbackQueryConfig CallbackQueryConfig
}
