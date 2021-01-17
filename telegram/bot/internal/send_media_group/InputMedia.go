package send_media_group

import (
	"encoding/json"
)

//easyjson:json
type InputMediaPhoto struct {
	Type      string  `json:"type"`
	Media     string  `json:"media"`
	Caption   *string `json:"caption,omitempty"`
	ParseMode *string `json:"parse_mode,omitempty"`
}

//easyjson:json
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

//easyjson:json
type InputMediaDocument struct {
	Type                        string  `json:"type"`
	Media                       string  `json:"media"`
	Thumb                       *string `json:"thumb,omitempty"`
	Caption                     *string `json:"caption,omitempty"`
	ParseMode                   *string `json:"parse_mode,omitempty"`
	DisableContentTypeDetection *bool   `json:"disable_content_type_detection,omitempty"`
}

//easyjson:json
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

type InputMedia interface {
	json.Marshaler
}
