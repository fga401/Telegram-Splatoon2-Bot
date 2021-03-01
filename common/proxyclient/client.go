package proxyclient

import (
	"crypto/tls"
	"net/http"
)

// New returns a new proxy http client.
func New(config Config) *http.Client {
	proxy := proxyURL(config.ProxyURL)
	if !config.EnableProxy {
		proxy = nil
	}
	var TLSNextProto map[string]func(authority string, c *tls.Conn) http.RoundTripper
	if !config.EnableHTTP2 {
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
