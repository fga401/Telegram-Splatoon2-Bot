package service

import (
	"bytes"
	"crypto/tls"
	botapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/nfnt/resize"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"image"
	"image/draw"
	"image/png"
	"net/http"
	proxy2 "telegram-splatoon2-bot/common/proxy"
	log "telegram-splatoon2-bot/logger"
	"telegram-splatoon2-bot/nintendo"
)

var client *http.Client

func InitImageClient() {
	// disable http 2
	useProxy := viper.GetBool("nintendo.useProxy")
	proxy := proxy2.GetProxy()
	if viper.InConfig("nintendo.proxyUrl") {
		proxy = proxy2.GetProxyWithUrl(viper.GetString("nintendo.proxyUrl"))
	}
	if !useProxy {
		proxy = nil
	}
	client = &http.Client{
		Transport: &http.Transport{
			TLSNextProto: make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
			Proxy:        proxy,
		},
	}
}

func downloadImage(url string) (image.Image, error) {
	var resp *http.Response
	err := retry(func() error {
		var err error
		resp, err = client.Get(url)
		return err
	}, retryTimes)
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
	log.Info("get image from url", zap.String("url", url))
	return img, nil
}

func downloadImages(urls []string) ([]image.Image, error) {
	imgs := make([]image.Image, len(urls))
	errChan := make(chan error, len(urls))
	for i, url := range urls {
		go func(idx int, u string) {
			img, err := downloadImage(u)
			if err != nil {
				errChan <- err
			}
			imgs[idx] = img
			errChan <- err
		}(i, url)
	}
	for range urls {
		err := <-errChan
		if err != nil {
			return nil, errors.Wrap(err, "can't download images")
		}
	}
	return imgs, nil
}

func uploadImage(img image.Image, name string) (string, error) {
	buf := bytes.NewBuffer(nil)
	err := png.Encode(buf, img)
	if err != nil {
		return "", errors.Wrap(err, "can't encode image")
	}
	msg := botapi.NewPhotoUpload(storeChannelID, botapi.FileBytes{Name: name + ".png", Bytes: buf.Bytes()})
	msg.Caption = "test"
	respMsg, err := sendWithRetryAndResponse(bot, msg)
	if err != nil {
		return "", errors.Wrap(err, "no response photo")
	}
	photo := *respMsg.Photo
	return photo[0].FileID, nil
}

func concatSalmonScheduleImage(detail *nintendo.SalmonScheduleDetail) (image.Image, error) {
	urls := []string{
		nintendo.Host + detail.Stage.Image,
		nintendo.Host + detail.Weapons[0].Weapon.Image,
		nintendo.Host + detail.Weapons[1].Weapon.Image,
		nintendo.Host + detail.Weapons[2].Weapon.Image,
		nintendo.Host + detail.Weapons[3].Weapon.Image,
	}
	imgs, err := downloadImages(urls)
	if err != nil {
		return nil, err
	}
	stage := imgs[0]
	weapons := imgs[1:5]
	width := stage.Bounds().Dx()
	qtrWidth := width / 4
	height := stage.Bounds().Dy() + qtrWidth
	// resize
	for i, img := range weapons {
		weapons[i] = resize.Resize(uint(qtrWidth), uint(qtrWidth), img, resize.Lanczos3)
	}
	// prepare canvas
	r := image.Rectangle{Min: image.Point{}, Max: image.Point{X: width, Y: height}}
	rgba := image.NewRGBA(r)
	draw.Draw(rgba,
		image.Rectangle{Min: image.Point{Y: qtrWidth}, Max: image.Point{X: width, Y: height}},
		stage, image.Point{}, draw.Src)
	for i, img := range weapons {
		draw.Draw(rgba,
			image.Rectangle{Min: image.Point{X: i * qtrWidth}, Max: image.Point{X: (i + 1) * qtrWidth, Y: qtrWidth}},
			img, image.Point{}, draw.Src)
	}
	return rgba, nil
}
