package repository

import "time"

// Repository fetches data from external service periodically and maintains the status.
type Repository interface {
	// Name is the unique ID of Repository.
	Name() string
	// Update fetches data from external service.
	Update() error
	// NextUpdateTime returns the next update time.
	NextUpdateTime() time.Time
}
