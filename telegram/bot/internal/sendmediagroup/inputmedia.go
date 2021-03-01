package sendmediagroup

import (
	"encoding/json"
)

// InputMediaPhoto sets up a photo JSON to upload.
type InputMediaPhoto struct {
	Type      string  `json:"type"`
	Media     string  `json:"media"`
	Caption   *string `json:"caption,omitempty"`
	ParseMode *string `json:"parse_mode,omitempty"`
}

// InputMediaVideo sets up a video JSON to upload.
type InputMediaVideo struct {
	Type              string  `json:"type"`
	Media             string  `json:"media"`
	Thumb             *string `json:"thumb,omitempty"`
	Caption           *string `json:"caption,omitempty"`
	ParseMode         *string `json:"parse_mode,omitempty"`
	Width             *int    `json:"width,omitempty"`
	Height            *int    `json:"height,omitempty"`
	Duration          *int    `json:"duration,omitempty"`
	SupportsStreaming *bool   `json:"supports_streaming,omitempty"`
}

// InputMediaDocument sets up a document JSON to upload.
type InputMediaDocument struct {
	Type                        string  `json:"type"`
	Media                       string  `json:"media"`
	Thumb                       *string `json:"thumb,omitempty"`
	Caption                     *string `json:"caption,omitempty"`
	ParseMode                   *string `json:"parse_mode,omitempty"`
	DisableContentTypeDetection *bool   `json:"disable_content_type_detection,omitempty"`
}

// InputMediaAudio sets up a audio JSON to upload.
type InputMediaAudio struct {
	Type      string  `json:"type"`
	Media     string  `json:"media"`
	Thumb     *string `json:"thumb,omitempty"`
	Caption   *string `json:"caption,omitempty"`
	ParseMode *string `json:"parse_mode,omitempty"`
	Duration  *int    `json:"duration,omitempty"`
	Performer *string `json:"performer,omitempty"`
	Title     *string `json:"title,omitempty"`
}

// InputMedia sets up a file JSON to upload.
type InputMedia interface {
	json.Marshaler
}
