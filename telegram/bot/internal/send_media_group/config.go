package send_media_group

type BasicFileConfig struct {
	Name      string
	Data      []byte
}

type FileConfig interface {
	InputMediaConfig() InputMedia
	File() BasicFileConfig
}

type Config struct {
	ChatID                   string
	File                     []FileConfig
	DisableNotification      *bool
	ReplyToMessageID         *int64
	AllowSendingWithoutReply *bool
}

type PhotoConfig struct {
	Name      string
	Data      []byte
	Caption   *string
	ParseMode *string
}

func (c PhotoConfig) InputMediaConfig() InputMedia {
	return InputMediaPhoto{
		Type:      "photo",
		Media:     "attach://" + c.Name,
		Caption:   c.Caption,
		ParseMode: c.ParseMode,
	}
}

func (c PhotoConfig) File() BasicFileConfig {
	return BasicFileConfig{
		Name: c.Name,
		Data: c.Data,
	}}

type VideoConfig struct {
	Name              string
	Data              []byte
	ThumbName         *string
	Caption           *string
	ParseMode         *string
	Width             *int
	Height            *int
	Duration          *int
	SupportsStreaming *bool
}

func (c VideoConfig) InputMediaConfig() InputMedia {
	return InputMediaVideo{
		Type:              "video",
		Media:             "attach://" + c.Name,
		Thumb:             parseThumb(c.ThumbName),
		Caption:           c.Caption,
		ParseMode:         c.ParseMode,
		Width:             c.Width,
		Height:            c.Height,
		Duration:          c.Duration,
		SupportsStreaming: c.SupportsStreaming,
	}
}

func (c VideoConfig) File() BasicFileConfig {
	return BasicFileConfig{
		Name: c.Name,
		Data: c.Data,
	}
}

type DocumentConfig struct {
	Name                        string
	Data                        []byte
	ThumbName                   *string
	Caption                     *string
	ParseMode                   *string
	DisableContentTypeDetection *bool
}

func (c DocumentConfig) InputMediaConfig() InputMedia {
	return InputMediaDocument{
		Type:                        "document",
		Media:                       "attach://" + c.Name,
		Thumb:                       parseThumb(c.ThumbName),
		Caption:                     c.Caption,
		ParseMode:                   c.ParseMode,
		DisableContentTypeDetection: c.DisableContentTypeDetection,
	}
}

func (c DocumentConfig) File() BasicFileConfig {
	return BasicFileConfig{
		Name: c.Name,
		Data: c.Data,
	}}

type AudioConfig struct {
	Name      string
	Data      []byte
	ThumbName *string
	Caption   *string
	ParseMode *string
	Duration  *int
	Performer *string
	Title     *string
}

func (c AudioConfig) InputMediaConfig() InputMedia {
	return InputMediaAudio{
		Type:      "audio",
		Media:     "attach://" + c.Name,
		Thumb:     parseThumb(c.ThumbName),
		Caption:   c.Caption,
		ParseMode: c.ParseMode,
		Duration:  c.Duration,
		Performer: c.Performer,
		Title:     c.Title,
	}
}

func (c AudioConfig) File() BasicFileConfig {
	return BasicFileConfig{
		Name: c.Name,
		Data: c.Data,
	}}

func parseThumb(thumb *string) *string {
	if thumb == nil {
		return nil
	}
	temp := "attach://" + *thumb
	return &temp
}
