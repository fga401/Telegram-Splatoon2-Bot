package sendmediagroup

// BasicFileConfig sets up the base info of a file to upload.
type BasicFileConfig struct {
	// Name, the file name.
	Name string
	// Data, the file data.
	Data []byte
}

// FileConfig sets up the file to upload.
type FileConfig interface {
	// InputMediaConfig converts a FileConfig to InputMedia
	InputMediaConfig() InputMedia
	// File returns the BasicFileConfig of FileConfig
	File() BasicFileConfig
}


// Config sets up the SendMediaGroup request.
// More info: https://core.telegram.org/bots/api#sendmediagroup
type Config struct {
	ChatID                   string
	File                     []FileConfig
	DisableNotification      *bool
	ReplyToMessageID         *int64
	AllowSendingWithoutReply *bool
}

// PhotoConfig sets up the Photo to upload.
// More info: https://core.telegram.org/bots/api#inputmediaphoto
type PhotoConfig struct {
	Name      string
	Data      []byte
	Caption   *string
	ParseMode *string
}

// InputMediaConfig converts a PhotoConfig to InputMedia
func (c PhotoConfig) InputMediaConfig() InputMedia {
	return InputMediaPhoto{
		Type:      "photo",
		Media:     "attach://" + c.Name,
		Caption:   c.Caption,
		ParseMode: c.ParseMode,
	}
}

// File returns the BasicFileConfig of PhotoConfig
func (c PhotoConfig) File() BasicFileConfig {
	return BasicFileConfig{
		Name: c.Name,
		Data: c.Data,
	}}

// VideoConfig sets up the Photo to upload.
// More info: https://core.telegram.org/bots/api#inputmediavideo
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

// InputMediaConfig converts a VideoConfig to InputMedia
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

// File returns the BasicFileConfig of VideoConfig
func (c VideoConfig) File() BasicFileConfig {
	return BasicFileConfig{
		Name: c.Name,
		Data: c.Data,
	}
}

// DocumentConfig sets up the Photo to upload.
// More info: https://core.telegram.org/bots/api#inputmediadocument
type DocumentConfig struct {
	Name                        string
	Data                        []byte
	ThumbName                   *string
	Caption                     *string
	ParseMode                   *string
	DisableContentTypeDetection *bool
}

// InputMediaConfig converts a DocumentConfig to InputMedia
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

// File returns the BasicFileConfig of DocumentConfig
func (c DocumentConfig) File() BasicFileConfig {
	return BasicFileConfig{
		Name: c.Name,
		Data: c.Data,
	}}

// AudioConfig sets up the Photo to upload.
// More info: https://core.telegram.org/bots/api#inputmediaaudio
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

// InputMediaConfig converts a AudioConfig to InputMedia
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

// File returns the BasicFileConfig of AudioConfig
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
