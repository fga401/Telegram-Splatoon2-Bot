package bot

import (
	sendMediaGroup "telegram-splatoon2-bot/telegram/bot/internal/sendmediagroup"
)

// SendMediaGroupConfig sets up the SendMediaGroup request.
// More info: https://core.telegram.org/bots/api#sendmediagroup
type SendMediaGroupConfig = sendMediaGroup.Config
// FileConfig sets up the file to upload.
type FileConfig = sendMediaGroup.FileConfig
// PhotoConfig sets up the Photo to upload.
// More info: https://core.telegram.org/bots/api#inputmediaphoto
type PhotoConfig = sendMediaGroup.PhotoConfig
// DocumentConfig sets up the Photo to upload.
// More info: https://core.telegram.org/bots/api#inputmediadocument
type DocumentConfig = sendMediaGroup.DocumentConfig
// VideoConfig sets up the Photo to upload.
// More info: https://core.telegram.org/bots/api#inputmediavideo
type VideoConfig = sendMediaGroup.VideoConfig
// AudioConfig sets up the Photo to upload.
// More info: https://core.telegram.org/bots/api#inputmediaaudio
type AudioConfig = sendMediaGroup.AudioConfig
