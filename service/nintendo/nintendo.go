package nintendo

import (
	"telegram-splatoon2-bot/service/language"
	"telegram-splatoon2-bot/service/timezone"
)

// AccountMetadata stores useful info fetched by session token.
type AccountMetadata struct {
	// IKSM of Splatoon2
	IKSM string
	// AccountName is the name of Nintendo account
	AccountName string
	// UserName is the name of Switch user
	UserName string
}

// Service manages all transactions about Nintendo.
type Service interface {
	// NewProofKey generates a new proof key.
	NewProofKey() ([]byte, error)
	// NewLoginLink generates a new login link by proof key.
	NewLoginLink(proofKey []byte) (string, error)
	// GetSessionToken fetches session token by proof key and user-input link.
	GetSessionToken(link string, proofKey []byte, language language.Language) (string, error)
	// GetSessionToken fetches AccountMetadata by session token.
	GetAccountMetadata(sessionToken string, language language.Language) (AccountMetadata, error)

	// GetSalmonSchedules fetches current salmon schedules.
	GetSalmonSchedules(iksm string, timezone timezone.Timezone, language language.Language) (SalmonSchedules, error)
	// GetSalmonSchedules fetches current stage schedules.
	GetStageSchedules(iksm string, timezone timezone.Timezone, language language.Language) (StageSchedules, error)
}
