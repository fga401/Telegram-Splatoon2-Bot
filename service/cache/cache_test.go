package cache

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	log "telegram-splatoon2-bot/logger"
	"telegram-splatoon2-bot/service/db"
	"testing"
)

func prepareTest() {
	viper.SetConfigName("dev")
	viper.SetConfigType("json")
	viper.AddConfigPath("../config/")
	viper.AddConfigPath("./config/")
	viper.ReadInConfig()
	log.InitLogger()
	InitCache()
}

func TestSetAndGetValue(t *testing.T) {
	prepareTest()
	expected := &db.Runtime{
		Uid:          123456,
		SessionToken: "654321",
		IKSM:         "123456789",
		Language:     "987654321",
		Timezone:     0,
	}
	user := &tgbotapi.User{
		ID:           147258369,
		FirstName:    "",
		LastName:     "",
		UserName:     "",
		LanguageCode: "",
		IsBot:        false,
	}
	err := Cache.SetRuntime(user, expected)
	assert.Nil(t, err)
	actual, err := Cache.GetRuntime(user)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}
