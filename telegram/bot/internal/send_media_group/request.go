package send_media_group

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"strconv"

	botApi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mailru/easyjson"
	"github.com/pkg/errors"
)

type RequestBuilder struct {
	bot         *botApi.BotAPI
	config      Config
	buf         *bytes.Buffer
	writer      *multipart.Writer
	mediaConfig []easyjson.RawMessage
}

	func NewRequestBuilder(bot *botApi.BotAPI, config Config) *RequestBuilder {
	buf := new(bytes.Buffer)
	return &RequestBuilder{
		bot:         bot,
		config:      config,
		buf:         buf,
		writer:      multipart.NewWriter(buf),
		mediaConfig: make([]easyjson.RawMessage, 0),
	}
}

func (r *RequestBuilder) Build() (*http.Request, error) {
	// add chat_id
	err := r.writer.WriteField("chat_id", r.config.ChatID)
	if err != nil {
		return nil, errors.Wrap(err, "can't add chat_id")
	}
	// add disable_notification
	if r.config.DisableNotification != nil {
		err := r.writer.WriteField("disable_notification", strconv.FormatBool(*r.config.DisableNotification))
		if err != nil {
			return nil, errors.Wrap(err, "can't add disable_notification")
		}
	}
	// add reply_to_message_id
	if r.config.ReplyToMessageID != nil {
		err := r.writer.WriteField("reply_to_message_id", strconv.FormatInt(*r.config.ReplyToMessageID, 10))
		if err != nil {
			return nil, errors.Wrap(err, "can't add reply_to_message_id")
		}
	}
	// add allow_sending_without_reply
	if r.config.AllowSendingWithoutReply != nil {
		err := r.writer.WriteField("allow_sending_without_reply", strconv.FormatBool(*r.config.AllowSendingWithoutReply))
		if err != nil {
			return nil, errors.Wrap(err, "can't add allow_sending_without_reply")
		}
	}
	// add media
	mediaJson, err := json.Marshal(r.mediaConfig)
	if err != nil {
		return nil, errors.Wrap(err, "can't marshal media")
	}
	err = r.writer.WriteField("media", string(mediaJson))
	if err != nil {
		return nil, errors.Wrap(err, "can't add media")
	}
	err = r.writer.Close()
	if err != nil {
		return nil, errors.Wrap(err, "can't close request writer")
	}

	request, err := http.NewRequest("POST", sendMediaGroupUrl(r.bot), r.buf)
	if err != nil {
		return nil, errors.Wrap(err, "can't new a request")
	}

	request.Header.Set("content-Type", r.writer.FormDataContentType())
	return request, nil
}

func (r *RequestBuilder) AddFile(file FileConfig) error {
	basicFile := file.File()
	fileWriter, err := r.writer.CreateFormFile(basicFile.Name, basicFile.Name)
	if err != nil {
		return errors.Wrap(err, "can't add "+basicFile.Name)
	}
	_, err = fileWriter.Write(basicFile.Data)
	if err != nil {
		return errors.Wrap(err, "can't add data of "+basicFile.Name)
	}
	raw, err := file.InputMediaConfig().MarshalJSON()
	if err != nil {
		return errors.Wrap(err, "can't add data of "+basicFile.Name)
	}
	r.mediaConfig = append(r.mediaConfig, raw)
	return nil
}
