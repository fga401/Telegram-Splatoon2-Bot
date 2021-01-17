package telegram

import (
	"strconv"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"telegram-splatoon2-bot/common/log"
	"telegram-splatoon2-bot/common/proxyclient"
)

type UploaderConfig struct {
	storeChannelID int64
	proxy          proxyclient.Config
}

func ReadUploaderConfig(viper viper.Viper) UploaderConfig {
	storeChannelID, err := strconv.ParseInt(viper.GetString("store_channel"), 10, 64)
	if err != nil {
		log.Panic("can't parse Image Config: StoreChannelID", zap.Error(err))
	}
	return UploaderConfig{
		storeChannelID: storeChannelID,
		proxy: proxyclient.Config{
			EnableProxy: viper.GetBool("nintendo.enable"), // todo(refactor)
			ProxyUrl:    viper.GetString("nintendo.proxyUrl"),
		},
	}
}

