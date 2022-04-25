package sendmediagroup

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	botApi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	log "telegram-splatoon2-bot/common/log"
	"telegram-splatoon2-bot/common/util"
	botUtils "telegram-splatoon2-bot/telegram/bot/internal/limit"
)

// Do builds a request from config and sends it.
// If too files sent and blocked by telegram, it will wait a moment and retry
func Do(bot *botApi.BotAPI, config Config, retryTimes int) ([]*botApi.Message, error) {
	builder := NewRequestBuilder(bot, config)
	for _, file := range config.File {
		err := builder.AddFile(file)
		if err != nil {
			return nil, errors.Wrap(err, "can't build request")
		}
	}
	req, err := builder.Build()
	if err != nil {
		return nil, errors.Wrap(err, "can't build request")
	}

	var apiResp botApi.APIResponse
	err = util.Retry(func() error {
		var resp *http.Response
		var err error
		resp, err = bot.Client.Do(req)
		if err != nil {
			return errors.Wrap(err, "can't get response")
		}
		defer resp.Body.Close()
		// parse APIResponse
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, "can't read response")
		}
		err = json.Unmarshal(data, &apiResp)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("can't unmarshal response, resp=%s", data))
		}
		if !apiResp.Ok {
			if is, sec := botUtils.IsTooManyRequestString(apiResp.Description); is {
				log.Warn("upload media blocked by telegram request limits", zap.Int("after", sec))
				time.Sleep(time.Duration(sec) * time.Second)
			}
			return errors.Errorf(apiResp.Description)
		}
		return err
	}, retryTimes)
	if err != nil {
		return nil, err
	}
	// parse Message
	var messages []*botApi.Message
	err = json.Unmarshal(apiResp.Result, &messages)
	if err != nil {
		return nil, errors.Wrap(err, "can't unmarshal message")
	}
	return messages, nil
}
