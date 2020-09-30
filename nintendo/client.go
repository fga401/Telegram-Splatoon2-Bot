package nintendo

import (
	"crypto/tls"
	"github.com/spf13/viper"
	"net/http"
	proxy2 "telegram-splatoon2-bot/common/proxy"
)

var client *http.Client

func InitClient() {
	// disable http 2
	useProxy := viper.GetBool("nintendo.useProxy")
	proxy := proxy2.GetProxy()
	if viper.InConfig("nintendo.proxyUrl"){
		proxy = proxy2.GetProxyWithUrl(viper.GetString("nintendo.proxyUrl"))
	}
	if !useProxy {
		proxy = nil
	}
	client = &http.Client{
		Transport: &http.Transport{
			TLSNextProto: make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
			Proxy: proxy,
		},
	}
}
