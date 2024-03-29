package telegram

import (
	"bytes"
	"image"
	"image/png"
	"strconv"

	botApi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	log "telegram-splatoon2-bot/common/log"
	imageSvc "telegram-splatoon2-bot/service/image"
	"telegram-splatoon2-bot/telegram/bot"
)

type telegramUploader struct {
	storeChannelID int64

	bot bot.Bot
}

// NewUploader returns a new Uploader.
func NewUploader(bot bot.Bot, config Config) imageSvc.Uploader {
	return &telegramUploader{
		storeChannelID: config.StoreChannelID,
		bot:            bot,
	}
}

func (s *telegramUploader) Upload(img image.Image) (imageSvc.Identifier, error) {
	uuid4, err := uuid.NewRandom()
	if err != nil {
		return "", errors.Wrap(err, "can't generate uuid4")
	}
	name := uuid4.String()
	buf := bytes.NewBuffer(nil)
	err = png.Encode(buf, img)
	if err != nil {
		return "", errors.Wrap(err, "can't encode image")
	}
	msg := botApi.NewPhotoUpload(s.storeChannelID, botApi.FileBytes{Name: name + ".png", Bytes: buf.Bytes()})
	respMsg, err := s.bot.Send(msg)
	if err != nil {
		return "", errors.Wrap(err, "no response photo")
	}
	photo := *respMsg.Photo
	return imageSvc.Identifier(photo[len(photo)-1].FileID), nil
}

func (s *telegramUploader) UploadAll(images []image.Image) ([]imageSvc.Identifier, error) {
	batchSize := 1
	if len(images) == 0 {
		return []imageSvc.Identifier{}, nil
	}
	if len(images) == 1 {
		id, err := s.Upload(images[0])
		log.Info("upload image done", zap.String("file_id", string(id)))
		return []imageSvc.Identifier{id}, err
	}
	ids := make([]imageSvc.Identifier, len(images))
	for i := 0; i < len(images); i += batchSize {
		sup := min(i+batchSize, len(images))
		if sup-i == 1 {
			var err error
			ids[i], err = s.Upload(images[i])
			if err != nil {
				return nil, err
			}
		} else {
			files, err := buildFiles(images[i:sup])
			if err != nil {
				return nil, errors.Wrap(err, "can't buildFiles images and names")
			}
			config := bot.SendMediaGroupConfig{
				ChatID: strconv.FormatInt(s.storeChannelID, 10),
				File:   files,
			}
			messages, err := s.bot.SendMediaGroup(config)
			if err != nil {
				return nil, errors.Wrap(err, "can't SendMediaGroup")
			}
			for j := i; j < sup; j++ {
				photo := *messages[j-i].Photo
				ids[j] = imageSvc.Identifier(photo[len(photo)-1].FileID)
			}
		}
		log.Info("upload part of images done", zap.Int("from", i), zap.Int("to", sup))
	}
	var stringIDs []string
	for _, id := range ids {
		stringIDs = append(stringIDs, string(id))
	}
	log.Info("upload multiple images done", zap.Strings("file_ids", stringIDs))
	return ids, nil
}

func buildFiles(images []image.Image) ([]bot.FileConfig, error) {
	ret := make([]bot.FileConfig, 0, len(images))
	for i := 0; i < len(images); i++ {
		uuid4, err := uuid.NewRandom()
		if err != nil {
			return nil, errors.Wrap(err, "can't generate uuid4")
		}
		name := uuid4.String()
		buf := bytes.NewBuffer(nil)
		err = png.Encode(buf, images[i])
		if err != nil {
			return nil, errors.Wrap(err, "can't encode image")
		}
		ret = append(ret, &bot.PhotoConfig{
			Name: name,
			Data: buf.Bytes(),
		})
	}
	return ret, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
