package proxyclient

import "time"

// Config sets up a proxy client.
type Config struct {
	// EnableProxy identifies whether The Client uses proxy.
	EnableProxy bool
	// ProxyURL If it's an empty string, it will be set by environment variables.
	ProxyURL string

	// EnableHTTP2 identifies whether The Client enables HTTP/2.
	EnableHTTP2 bool
	// Timeout of request.
	Timeout time.Duration
}
