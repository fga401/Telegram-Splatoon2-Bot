package database

// Config sets up a Database.
type Config struct {
	// URL of connection.
	URL          string
	// Driver of Database.
	Driver       string
	// MaxIdleConns, the maximum number of connections in the idle connection pool.
	MaxIdleConns int
	// MaxOpenConns, the maximum number of open connections to the database.
	MaxOpenConns int
}
