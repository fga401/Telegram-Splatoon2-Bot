package downloader

import (
	"context"
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
)

type impl struct {
	client     *http.Client
	retryTimes int
}

func NewDownloader(config Config) imageSvc.Downloader {
	return &impl{
		client:     proxyclient.New(config.Proxy),
		retryTimes: config.RetryTimes,
	}
}

func (s *impl) download(ctx context.Context, url string) (image.Image, error) {
	switch {
	case strings.HasPrefix(url, "file://"):
		return s.downloadFromFile(url)
	case strings.HasPrefix(url, "https://"):
		return s.downloadFromNet(ctx, url)
	case strings.HasPrefix(url, "http://"):
		return s.downloadFromNet(ctx, url)
	default:
		return nil, errors.Errorf("unknown url scheme")
	}
}

func (s *impl) Download(url string) (image.Image, error) {
	return s.download(context.Background(), url)
}

func (s *impl) DownloadAll(urls []string) ([]image.Image, error) {
	imgs := make([]image.Image, len(urls))
	errChan := make(chan error, len(urls))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for i, url := range urls {
		go func(idx int, u string) {
			img, err := s.download(ctx, u)
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

func (s *impl) downloadFromNet(ctx context.Context, url string) (image.Image, error) {
	var resp *http.Response
	err := util.Retry(func() error {
		var err error
		req, err := http.NewRequestWithContext(ctx,"GET", url, nil)
		if err != nil {
			return errors.Wrap(err, "can't make request")
		}
		resp, err = s.client.Do(req)
		return err
	}, s.retryTimes)
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
