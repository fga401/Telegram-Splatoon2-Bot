package nintendo

import (
	"telegram-splatoon2-bot/service/language"
	"telegram-splatoon2-bot/service/timezone"
)

type AccountMetadata struct {
	IKSM        string // Splatoon2 IKSM
	AccountName string // Name of Nintendo account
	UserName    string // Name of Switch user
}

type Service interface {
	NewProofKey() ([]byte, error)
	NewLoginLink(proofKey []byte) (string, error)
	GetSessionToken(link string, proofKey []byte, language language.Language) (string, error)
	GetAccountMetadata(sessionToken string, language language.Language) (AccountMetadata, error)

	GetSalmonSchedules(iksm string, timezone timezone.Timezone, language language.Language) (SalmonSchedules, error)
	GetStageSchedules(iksm string, timezone timezone.Timezone, language language.Language) (StageSchedules, error)
}
