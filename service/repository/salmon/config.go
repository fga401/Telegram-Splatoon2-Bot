package salmon

type DumperConfig struct {
	WeaponFile string
	StageFile  string
}

type Config struct {
	Dumper DumperConfig
	RandomWeaponPath string
	GrizzcoWeaponPath string
}
