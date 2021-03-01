package repository

import "time"

type Repository interface {
	Name() string
	Update() error
	NextUpdateTime() time.Time
}

