package sendmediagroup

import (
	"fmt"

	botApi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func sendMediaGroupURL(bot *botApi.BotAPI) string {
	return fmt.Sprintf("https://api.telegram.org/bot%s/sendMediaGroup", bot.Token)
}
