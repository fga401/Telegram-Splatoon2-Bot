package nintendo

import (
	"net/http"

	proxyClient "telegram-splatoon2-bot/common/proxyclient"
)

type impl struct {
	client     *http.Client
	retryTimes int
}

func New(config Config) Service {
	client := proxyClient.New(proxyClient.Config{
		EnableHttp2: false,
		Timeout:     config.Timeout,
	})
	return &impl{
		client:     client,
		retryTimes: config.RetryTimes,
	}
}
