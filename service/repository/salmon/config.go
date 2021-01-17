package salmon

import "github.com/spf13/viper"

type DumperConfig struct {
	WeaponFile string
	StageFile  string
}

func ReadConfig(viper viper.Viper) Config {
	return Config{
		dumper: DumperConfig{
			StageFile:  viper.GetString("service.salmon.stageFileName"),
			WeaponFile: viper.GetString("service.salmon.weaponFileName"),
		},
	}
}

type Config struct {
	dumper DumperConfig
}
