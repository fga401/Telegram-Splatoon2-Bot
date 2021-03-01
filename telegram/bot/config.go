package bot

type CallbackQueryConfig struct {
	Text      string
	ShowAlert bool
	CacheTime int
}

type Config struct {
	RetryTimes                 int
	DefaultCallbackQueryConfig CallbackQueryConfig
}
