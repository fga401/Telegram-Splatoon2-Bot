package downloader

import (
	"telegram-splatoon2-bot/common/proxyclient"
)

type Config struct {
	Proxy proxyclient.Config
	RetryTimes int
}
