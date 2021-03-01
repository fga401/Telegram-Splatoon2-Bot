package downloader

import (
	"telegram-splatoon2-bot/common/proxyclient"
)

// Config sets up a Downloader.
type Config struct {
	// Proxy config of http client.
	Proxy      proxyclient.Config
	// RetryTimes after failure.
	RetryTimes int
}
