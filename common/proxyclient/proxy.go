package proxyclient

import (
	"net/http"
	"net/url"
	"os"
	"strings"

	"go.uber.org/zap"
	log "telegram-splatoon2-bot/common/log"
)

func defaultProxyUrl() *url.URL {
	proxies := []string{"socks5_proxy", "https_proxy", "http_proxy"}
	var proxyUrl *url.URL
	for _, proxy := range proxies {
		proxyEnv := os.Getenv(strings.ToLower(proxy))
		if proxyEnv == "" {
			proxyEnv = os.Getenv(strings.ToUpper(proxy))
		}
		if proxyEnv != "" {
			innerProxyUrl, err := url.Parse(proxyEnv)
			if err != nil {
				continue
			}
			proxyUrl = innerProxyUrl
			break
		}
	}
	return proxyUrl
}

func proxyUrl(u string) func(*http.Request) (*url.URL, error) {
	proxyUrl := defaultProxyUrl()
	if u != "" {
		var err error
		proxyUrl, err = url.Parse(u)
		if err != nil{
			log.Warn("can't parse proxy url",zap.String("url", u), zap.Error(err))
		}
	}
	return http.ProxyURL(proxyUrl)
}

