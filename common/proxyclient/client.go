package proxyclient

import (
	"crypto/tls"
	"net/http"
)

func New(config Config) *http.Client {
	proxy := proxyUrl(config.ProxyUrl)
	if !config.EnableProxy {
		proxy = nil
	}
	var TLSNextProto map[string]func(authority string, c *tls.Conn) http.RoundTripper
	if !config.EnableHttp2 {
		TLSNextProto = make(map[string]func(authority string, c *tls.Conn) http.RoundTripper)
	}
	return &http.Client{
		Timeout: config.Timeout,
		Transport: &http.Transport{
			TLSNextProto: TLSNextProto,
			Proxy:        proxy,
		},
	}
}
