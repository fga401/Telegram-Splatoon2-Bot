package stage

import "github.com/spf13/viper"

type DumperConfig struct {
	StageFile string
}

func ReadConfig(viper viper.Viper) Config {
	return Config{
		dumper: DumperConfig{
			StageFile: viper.GetString("service.stage.stageFileName"),
		},
	}
}

type Config struct {
	dumper DumperConfig
}
