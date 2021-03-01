package stage

// DumperConfig sets up a stage Dumper.
type DumperConfig struct {
	// StageFile is the path of the stage dumping file.
	StageFile string
}

// Config sets up a stage Repository.
type Config struct {
	// Dumper sets up a stage Dumper.
	Dumper DumperConfig
}
