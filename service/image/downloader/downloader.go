package downloader

import (
	"image"
	"net/http"
	"os"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	log "telegram-splatoon2-bot/common/log"
	"telegram-splatoon2-bot/common/proxyclient"
	"telegram-splatoon2-bot/common/util"
	imageSvc "telegram-splatoon2-bot/service/image"
	"telegram-splatoon2-bot/service/todo"
)

type impl struct {
	client *http.Client
}

func NewDownloader(config Config) imageSvc.Downloader {
	return &impl{
		client: proxyclient.New(config.proxy),
	}
}

func (s *impl) Download(url string) (image.Image, error) {
	switch {
	case strings.HasPrefix(url, "file://"):
		return s.downloadFromFile(url)
	case strings.HasPrefix(url, "https://"):
		return s.downloadFromNet(url)
	case strings.HasPrefix(url, "http://"):
		return s.downloadFromNet(url)
	default:
		return nil, errors.Errorf("unknown url scheme")
	}
}

func (s *impl) DownloadAll(urls []string) ([]image.Image, error) {
	imgs := make([]image.Image, len(urls))
	errChan := make(chan error, len(urls))
	for i, url := range urls {
		go func(idx int, u string) {
			img, err := s.Download(u)
			imgs[idx] = img
			errChan <- err
		}(i, url)
	}
	for range urls {
		err := <-errChan
		if err != nil {
			return nil, errors.Wrap(err, "can't downloader images")
		}
	}
	return imgs, nil
}

func (s *impl) downloadFromNet(url string) (image.Image, error) {
	var resp *http.Response
	err := util.Retry(func() error {
		var err error
		resp, err = s.client.Get(url)
		return err
	}, todo.RetryTimes)
	if err != nil {
		return nil, errors.Wrap(err, "can't get resp")
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "can't decode image")
	}
	log.Info("get image from http(s) url", zap.String("url", url))
	return img, nil
}

func (s *impl) downloadFromFile(url string) (image.Image, error) {
	filePath := url[7:]
	imgFile, err := os.Open(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "can't open image file")
	}
	defer func() {
		_ = imgFile.Close()
	}()
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, errors.Wrap(err, "can't decode image")
	}
	log.Info("get image from local file", zap.String("url", url))
	return img, nil
}
