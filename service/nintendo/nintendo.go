package nintendo

import "telegram-splatoon2-bot/service/language"

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
}
