package nintendo

import "time"

// Config sets up a Service.
type Config struct {
	// Timeout of request.
	Timeout    time.Duration
	// RetryTimes after failure.
	RetryTimes int
}
