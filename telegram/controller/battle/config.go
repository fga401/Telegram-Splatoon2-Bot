package battle

// Config sets up a Battle.
type Config struct {
	// MaxResultsPerMessage sets the max number of results presented in one telegram message.
	MaxResultsPerMessage int
	// MinLastResults sets the min number of results shown by /battle_last.
	MinLastResults int
	// PollingMaxWorker sets the max number of goroutine to send polled battles .
	// If PollingMaxWorker == 0, there is no limitation.
	PollingMaxWorker int32
}
