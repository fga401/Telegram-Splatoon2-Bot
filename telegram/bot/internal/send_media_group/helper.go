package send_media_group

import (
	"fmt"

	botApi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func sendMediaGroupUrl(bot *botApi.BotAPI) string {
	return fmt.Sprintf("https://api.telegram.org/bot%s/sendMediaGroup", bot.Token)
}