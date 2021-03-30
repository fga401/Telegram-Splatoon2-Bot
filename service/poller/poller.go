package poller

import "telegram-splatoon2-bot/service/user"

// Poller polls results from external service
type Poller interface {
	// Start adds a user to the poller. It would not be blocked.
	Start(id user.ID)
	// Start removes a user from the poller. It would not be blocked.
	Stop(id user.ID)
}
