package dump

// Config sets up a Dumper.
type Config struct {
	targets map[string]string
}

// AddTarget adds a pair of key and file to the Dumper.
func (c *Config) AddTarget(key string, fileName string) *Config {
	if c.targets == nil {
		c.targets = make(map[string]string)
	}
	c.targets[key] = fileName
	return c
}
