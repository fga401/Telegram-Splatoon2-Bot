package proxyclient

import (
	"net/http"
	"net/url"
	"os"
	"strings"

	"go.uber.org/zap"
	log "telegram-splatoon2-bot/common/log"
)

func defaultProxyURL() *url.URL {
	proxies := []string{"socks5_proxy", "https_proxy", "http_proxy"}
	var proxyURL *url.URL
	for _, proxy := range proxies {
		proxyEnv := os.Getenv(strings.ToLower(proxy))
		if proxyEnv == "" {
			proxyEnv = os.Getenv(strings.ToUpper(proxy))
		}
		if proxyEnv != "" {
			innerProxyURL, err := url.Parse(proxyEnv)
			if err != nil {
				continue
			}
			proxyURL = innerProxyURL
			break
		}
	}
	return proxyURL
}

func proxyURL(u string) func(*http.Request) (*url.URL, error) {
	proxyURL := defaultProxyURL()
	if u != "" {
		var err error
		proxyURL, err = url.Parse(u)
		if err != nil {
			log.Warn("can't parse proxy url", zap.String("url", u), zap.Error(err))
		}
	}
	return http.ProxyURL(proxyURL)
}
