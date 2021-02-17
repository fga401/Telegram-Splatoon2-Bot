package image

import "image"

type Identifier string

type Service interface {
	Uploader
	Downloader
}

type Downloader interface {
	Download(url string) (image.Image, error)
	DownloadAll(urls []string) ([]image.Image, error)
}


type Uploader interface {
	Upload(img image.Image) (Identifier, error)
	UploadAll(imgs []image.Image) ([]Identifier, error)
}

type serviceImpl struct {
	Uploader
	Downloader
}

func NewService(uploader Uploader, downloader Downloader) Service {
	return &serviceImpl{uploader, downloader}
}