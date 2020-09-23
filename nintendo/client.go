package nintendo

import (
	"crypto/tls"
	"github.com/spf13/viper"
	"net/http"
	"telegram-splatoon2-bot/common"
)

var client *http.Client

func init() {
	// disable http 2
	useProxy := viper.GetBool("nintendo.useProxy")
	proxy := common.GetProxy()
	if viper.InConfig("nintendo.proxyUrl"){
		proxy = common.GetProxyWithUrl(viper.GetString("nintendo.proxyUrl"))
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
