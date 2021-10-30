package nintendo

import (
	"net/http"

	proxyClient "telegram-splatoon2-bot/common/proxyclient"
)

type impl struct {
	client     *http.Client
	retryTimes int
	appVersion string
}

// New returns a new Service object.
func New(config Config) Service {
	client := proxyClient.New(proxyClient.Config{
		EnableHTTP2: false,
		Timeout:     config.Timeout,
	})
	return &impl{
		client:     client,
		retryTimes: config.RetryTimes,
		appVersion: config.AppVersion,
	}
}
