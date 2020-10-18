package botutil

import (
	"bytes"
	"fmt"
	botapi "github.com/go-telegram-bot-api/telegram-bot-api"
	json "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"io/ioutil"
	"mime/multipart"
	"regexp"
	"strconv"
	log "telegram-splatoon2-bot/logger"
	"time"
)

func getUrlToSendMediaGroup(bot *botapi.BotAPI) string {
	return fmt.Sprintf("https://api.telegram.org/bot%s/sendMediaGroup", bot.Token)
}

type InputMediaPhoto struct {
	Name    string
	Caption string
	Data    []byte
}

// todo: other parameters
type InputMediaVideo struct {
	Name    string
	Caption string
	Data    []byte
}

func addPhoto(bodyWriter *multipart.Writer, media *[]map[string]interface{}, img *InputMediaPhoto) error {
	fileWriter, err := bodyWriter.CreateFormFile(img.Name, img.Name)
	if err != nil {
		return err
	}
	_, err = fileWriter.Write(img.Data)
	if err != nil {
		return err
	}
	*media = append(*media, map[string]interface{}{
		"type":       "photo",
		"media":      "attach://" + img.Name,
		"caption":    img.Caption,
		"parse_mode": "MarkdownV2",
	})
	return nil
}

func addVideo(bodyWriter *multipart.Writer, media *[]map[string]interface{}, video *InputMediaVideo) error {
	fileWriter, err := bodyWriter.CreateFormFile(video.Name, video.Name)
	if err != nil {
		return err
	}
	_, err = fileWriter.Write(video.Data)
	if err != nil {
		return err
	}
	*media = append(*media, map[string]interface{}{
		"type":       "photo",
		"media":      "attach://" + video.Name,
		"caption":    video.Caption,
		"parse_mode": "MarkdownV2",
	})
	return nil
}

type UploadMediaRequest struct {
	buffer []byte
	writer *multipart.Writer
}

func NewUploadMediaRequest(chatID string, inputMedia ...interface{}) (*UploadMediaRequest, error) {
	if len(inputMedia) < 2 || len(inputMedia) > 10 {
		return nil, errors.Errorf("number of image is illegal")
	}
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	fileWriter, err := bodyWriter.CreateFormField("chat_id")
	if err != nil {
		return nil, errors.Wrap(err, "can't add chat_id")
	}
	_, err = fileWriter.Write([]byte(chatID))
	if err != nil {
		return nil, errors.Wrap(err, "can't add chat_id")
	}

	mediaJsonObject := make([]map[string]interface{}, 0)
	for _, media := range inputMedia {
		switch m := media.(type) {
		case InputMediaPhoto:
			err = addPhoto(bodyWriter, &mediaJsonObject, &m)
		case *InputMediaPhoto:
			err = addPhoto(bodyWriter, &mediaJsonObject, m)
		case InputMediaVideo:
			err = addVideo(bodyWriter, &mediaJsonObject, &m)
		case *InputMediaVideo:
			err = addVideo(bodyWriter, &mediaJsonObject, m)
		default:
			continue
		}
		if err != nil {
			return nil, err
		}
	}

	fileWriter, err = bodyWriter.CreateFormField("media")
	if err != nil {
		return nil, errors.Wrap(err, "can't add chat_id")
	}
	mediaJson, err := json.Marshal(mediaJsonObject)
	if err != nil {
		return nil, errors.Wrap(err, "can't marshal media description")
	}
	_, err = fileWriter.Write(mediaJson)
	if err != nil {
		return nil, errors.Wrap(err, "can't add chat_id")
	}
	_ = bodyWriter.Close()
	return &UploadMediaRequest{bodyBuf.Bytes(), bodyWriter}, nil
}

func DoUploadMedia(bot *botapi.BotAPI, req *UploadMediaRequest) ([]*botapi.Message, error) {
	url := getUrlToSendMediaGroup(bot)
	contentType := req.writer.FormDataContentType()

	var apiResp botapi.APIResponse
	for {
		resp, err := bot.Client.Post(url, contentType, bytes.NewBuffer(req.buffer))
		if err != nil {
			return nil, errors.Wrap(err, "can't get response")
		}
		defer resp.Body.Close()
		// parse APIResponse
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.Wrap(err, "can't read response")
		}
		err = json.Unmarshal(data, &apiResp)
		if err != nil {
			return nil, errors.Wrap(err, "can't unmarshal response")
		}
		if !apiResp.Ok {
			if is, sec := IsTooManyRequestString(apiResp.Description); is {
				log.Warn("upload media blocked by telegram request limits", zap.Int("after", sec))
				time.Sleep(time.Duration(sec) * time.Second)
				continue
			} else {
				return nil, errors.Errorf(apiResp.Description)
			}
		}
		break
	}

	// parse Message
	var messages []*botapi.Message
	err := json.Unmarshal(apiResp.Result, &messages)
	if err != nil {
		return nil, errors.Wrap(err, "can't unmarshal message")
	}
	return messages, nil
}

var tooManyRequestRegExp = regexp.MustCompile(`Too Many Requests: retry after (?P<sec>\d+)`)

func IsTooManyRequestError(e error) (bool, int) {
	if e == nil {
		return false, 0
	}
	results := tooManyRequestRegExp.FindStringSubmatch(e.Error())
	if len(results) > 0 {
		sec, _ := strconv.Atoi(results[1])
		return true, sec
	}
	return false, 0
}

func IsTooManyRequestString(s string) (bool, int) {
	results := tooManyRequestRegExp.FindStringSubmatch(s)
	if len(results) > 0 {
		sec, _ := strconv.Atoi(results[1])
		return true, sec
	}
	return false, 0
}
