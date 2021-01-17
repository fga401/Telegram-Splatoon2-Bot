package database

type Config struct {
	URL          string
	Driver       string
	MaxIdleConns int
	MaxOpenConns int
}
