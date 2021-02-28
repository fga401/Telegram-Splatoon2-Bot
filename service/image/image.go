package image

import "image"

// Identifier of image.
type Identifier string

// Service manages download and upload of images.
type Service interface {
	Uploader
	Downloader
}

// Downloader manages download of images.
type Downloader interface {
	Download(url string) (image.Image, error)
	DownloadAll(urls []string) ([]image.Image, error)
}

// Uploader manages download of images.
type Uploader interface {
	Upload(img image.Image) (Identifier, error)
	UploadAll(imgs []image.Image) ([]Identifier, error)
}

type serviceImpl struct {
	Uploader
	Downloader
}

// NewService returns a new image Service.
func NewService(uploader Uploader, downloader Downloader) Service {
	return &serviceImpl{uploader, downloader}
}
