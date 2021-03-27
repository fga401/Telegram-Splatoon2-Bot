package battle

import "time"

// Config sets up a battle poller
type Config struct {
	// RefreshmentTime sets the interval between tow refreshment.
	RefreshmentTime time.Duration
	// MinBattleTime sets the min intervals of different modes.
	MinBattleTime MinBattleTime
	// MaxWorker sets the max number of goroutine to process request.
	// If MaxWorker <= 0, there is no limitation.
	MaxWorker int32
	// MaxWorker sets the max interval between tow battles.
	// If the time no new battles is longer than MaxIdleTime, the polling will be canceled.
	MaxIdleTime time.Duration
}

// MinBattleTime sets the min intervals of different modes.
type MinBattleTime struct {
	Zone      time.Duration
	Tower     time.Duration
	Clam      time.Duration
	Rainmaker time.Duration
	Waiting   time.Duration
}
