package dump

type Config struct {
	targets map[string]string
}

func (c *Config) AddTarget(key string, fileName string) *Config {
	if c.targets == nil {
		c.targets = make(map[string]string)
	}
	c.targets[key] = fileName
	return c
}