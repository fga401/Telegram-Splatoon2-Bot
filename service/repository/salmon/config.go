package salmon

// DumperConfig sets up a salmon Dumper.
type DumperConfig struct {
	// WeaponFile is the path of the weapon dumping file.
	WeaponFile string
	// StageFile is the path of the stage dumping file.
	StageFile string
}

// Config sets up a salmon Repository.
type Config struct {
	// Dumper sets up a salmon Dumper.
	Dumper DumperConfig
	// RandomWeaponPath is the path of the random weapon image.
	RandomWeaponPath string
	// GrizzcoWeaponPath is the path of the grizzco weapon image.
	GrizzcoWeaponPath string
}
