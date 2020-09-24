package common

import (
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"os"
	"strings"
	log "telegram-splatoon2-bot/logger"
)

func getProxyUrl() *url.URL {
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

func GetProxy() func(*http.Request) (*url.URL, error) {
	proxyUrl := getProxyUrl()
	if proxyUrl == nil{
		return nil
	}
	return http.ProxyURL(proxyUrl)
}

func GetProxyWithUrl(u string) func(*http.Request) (*url.URL, error) {
	proxyUrl, err := url.Parse(u)
	if err != nil{
		log.Warn("GetProxyWithUrl failed",zap.String("url", u), zap.Error(err))
	}
	return http.ProxyURL(proxyUrl)
}

