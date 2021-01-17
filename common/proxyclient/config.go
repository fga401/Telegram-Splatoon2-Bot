package proxyclient

import "time"

type Config struct {
	// The Client will use proxy if EnableProxy is true
	EnableProxy bool
	// Proxy url. If it's an empty string, it will be set by environment variables.
	ProxyUrl string

	EnableHttp2 bool
	Timeout     time.Duration
}
