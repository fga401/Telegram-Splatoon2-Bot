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
	"os"
	"strconv"
	"strings"
	"telegram-splatoon2-bot/botutil"
	proxy2 "telegram-splatoon2-bot/common/proxy"
	log "telegram-splatoon2-bot/logger"
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

type imageHelperPrivate struct {}
type imageHelper struct {
	privateMethod imageHelperPrivate
}

var (
	ImageHelper imageHelper
)

func (imageHelperPrivate)downloadImageFromNet(url string) (image.Image, error) {
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
	log.Info("get image from http(s) url", zap.String("url", url))
	return img, nil
}

func (imageHelperPrivate)downloadImageFromFile(url string) (image.Image, error) {
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

func (helper imageHelper)downloadImage(url string) (image.Image, error) {
	switch {
	case strings.HasPrefix(url, "file://"):
		return helper.privateMethod.downloadImageFromFile(url)
	case strings.HasPrefix(url, "https://"):
		return helper.privateMethod.downloadImageFromNet(url)
	case strings.HasPrefix(url, "http://"):
		return helper.privateMethod.downloadImageFromNet(url)
	default:
		return nil, errors.Errorf("")
	}

}

func (helper imageHelper)downloadImages(urls []string) ([]image.Image, error) {
	imgs := make([]image.Image, len(urls))
	errChan := make(chan error, len(urls))
	for i, url := range urls {
		go func(idx int, u string) {
			img, err := helper.downloadImage(u)
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

func (imageHelper)uploadImage(img image.Image, name string) (string, error) {
	buf := bytes.NewBuffer(nil)
	err := png.Encode(buf, img)
	if err != nil {
		return "", errors.Wrap(err, "can't encode image")
	}
	msg := botapi.NewPhotoUpload(storeChannelID, botapi.FileBytes{Name: name + ".png", Bytes: buf.Bytes()})
	msg.Caption = name
	respMsg, err := sendWithRetryAndResponse(bot, msg)
	if err != nil {
		return "", errors.Wrap(err, "no response photo")
	}
	photo := *respMsg.Photo
	return photo[len(photo)-1].FileID, nil
}

type FileItem struct {
	Key      string //image_content
	FileName string //test.jpg
	Content  []byte //[]byte
}

func (imageHelperPrivate)zipImageAndName(imgs []image.Image, names []string) ([]interface{}, error) {
	ret := make([]interface{}, len(imgs))
	for i := 0; i < len(imgs); i++ {
		buf := bytes.NewBuffer(nil)
		err := png.Encode(buf, imgs[i])
		if err != nil {
			return nil, errors.Wrap(err, "can't encode image")
		}
		ret[i] = &botutil.InputMediaPhoto{
			Name: names[i],
			Data: buf.Bytes(),
		}
	}
	return ret, nil
}

func (helper imageHelper)uploadImages(imgs []image.Image, names []string) ([]string, error) {
	if len(imgs) != len(names) {
		return nil, errors.Errorf("numbers of image and name are misatched")
	}
	if len(imgs) == 0 {
		return []string{}, nil
	}
	if len(imgs) == 1 {
		id, err := helper.uploadImage(imgs[0], names[0])
		return []string{id}, err
	}
	ids := make([]string, len(imgs))
	for i := 0; i < len(imgs); i += 10 {
		sup := MinInt(i+10, len(imgs))
		if sup-i == 1 {
			var err error
			ids[i], err = helper.uploadImage(imgs[i], names[i])
			if err != nil {
				return nil, err
			}
		} else {
			input, err := helper.privateMethod.zipImageAndName(imgs[i:sup], names[i:sup])
			if err != nil {
				return nil, err
			}
			req, err := botutil.NewUploadMediaRequest(strconv.FormatInt(storeChannelID, 10), input...)
			if err != nil {
				return nil, errors.Wrap(err, "can't generate new request")
			}
			var messages []*botapi.Message
			err = retry(func() error {
				var err error
				messages, err = botutil.DoUploadMedia(bot, req)
				return err
			}, retryTimes)
			if err != nil {
				return nil, errors.Wrap(err, "can't get response message")
			}
			for j := i; j < sup; j++ {
				photo := *messages[j-i].Photo
				ids[j] = photo[len(photo)-1].FileID
			}
		}
	}
	log.Info("upload multiple images done", zap.Strings("file_ids", ids))
	return ids, nil
}

// concatSalmonScheduleImage parameter:
//   imgs[0]: stage; imgs[1:5]: weapons
func concatSalmonScheduleImage(imgs []image.Image) image.Image {
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
	//draw
	draw.Draw(rgba,
		image.Rectangle{Min: image.Point{Y: qtrWidth}, Max: image.Point{X: width, Y: height}},
		stage, image.Point{}, draw.Src)
	for i, img := range weapons {
		draw.Draw(rgba,
			image.Rectangle{Min: image.Point{X: i * qtrWidth}, Max: image.Point{X: (i + 1) * qtrWidth, Y: qtrWidth}},
			img, image.Point{}, draw.Src)
	}
	return rgba
}

func MinInt(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func concatStageScheduleImage(imgA, imgB image.Image) image.Image {
	width := MinInt(imgA.Bounds().Dx(), imgB.Bounds().Dx())
	halfHeight := MinInt(imgA.Bounds().Dy(), imgB.Bounds().Dy())
	height := halfHeight * 2
	// resize
	imgA = resize.Resize(uint(width), uint(halfHeight), imgA, resize.Lanczos3)
	imgB = resize.Resize(uint(width), uint(halfHeight), imgB, resize.Lanczos3)
	// prepare canvas
	r := image.Rectangle{Min: image.Point{}, Max: image.Point{X: width, Y: height}}
	rgba := image.NewRGBA(r)
	//draw
	draw.Draw(rgba,
		image.Rectangle{Min: image.Point{}, Max: image.Point{X: width, Y: height / 2}},
		imgB, image.Point{}, draw.Src)
	draw.Draw(rgba,
		image.Rectangle{Min: image.Point{Y: height / 2}, Max: image.Point{X: width, Y: height}},
		imgA, image.Point{}, draw.Src)
	return rgba
}
