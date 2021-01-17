package user

import (
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"telegram-splatoon2-bot/common/log"
	"telegram-splatoon2-bot/service/language"
	"telegram-splatoon2-bot/service/nintendo"
	"telegram-splatoon2-bot/service/user/serializer"
)

func (svc *serviceImpl) NewLoginLink(uid ID) (string, error) {
	key := serializer.FromID(uid)
	if value := svc.proofKeyCache.Get(key); value != nil {
		// todo: login in link is existed, metrics
		log.Warn("login in link is existed", zap.Any("user_id", uid))
	}
	proofKey, err := svc.nintendoSvc.NewProofKey()
	if err != nil {
		return "", errors.Wrap(err, "can't generate proof key")
	}
	svc.proofKeyCache.SetExpiration(key, proofKey, svc.proofKeyCacheExpiration)
	link, err := svc.nintendoSvc.NewLoginLink(proofKey)
	if err != nil {
		return "", errors.Wrap(err, "can't generate login link")
	}
	return link, nil
}

func (svc *serviceImpl) AddAccount(uid ID, link string) error {
	key := serializer.FromID(uid)
	accounts, err := svc.ListAccounts(uid)
	if err != nil {
		return errors.Wrap(err, "can't fetch accounts")
	}
	svc.accountCache.Del(key)
	log.Debug("accounts cache delete", zap.Any("user_id", uid))

	proofKey := svc.proofKeyCache.Get(key)
	svc.proofKeyCache.Del(key)
	if proofKey == nil {
		return errors.New("expired proof key")
	}
	sessionToken, err := svc.nintendoSvc.GetSessionToken(link, proofKey, language.English)
	if err != nil {
		return errors.Wrap(err, "can't fetch session token")
	}
	nintendoAccount, err := svc.nintendoSvc.GetAccountMetadata(sessionToken, language.English)
	if err != nil {
		return errors.Wrap(err, "can't fetch nintendo account metadata")
	}
	account := Account{
		UserID:       uid,
		SessionToken: sessionToken,
		Tag:          formatTag(nintendoAccount),
	}
	if len(accounts) == 0 {
		svc.statusCache.Del(key)
		err = svc.db.InsertAndSwitchAccount(account, nintendoAccount.IKSM)
		log.Debug("status cache delete", zap.Any("user_id", uid))
		if err != nil {
			return errors.Wrap(err, "can't add new account and switch to it in database")
		}
	} else {
		err = svc.db.InsertAccount(account)
		if err != nil {
			return errors.Wrap(err, "can't add account to database")
		}
	}
	return nil
}

func formatTag(nintendoAccount nintendo.AccountMetadata) string {
	return nintendoAccount.AccountName + ":" + nintendoAccount.UserName
}

func (svc *serviceImpl) DeleteAccount(uid ID, tag string) error {
	status, err := svc.GetStatus(uid)
	if err != nil {
		return errors.Wrap(err, "can't get status to check the current account")
	}
	accounts, err := svc.ListAccounts(uid)
	if err != nil {
		return errors.Wrap(err, "can't fetch accounts")
	}
	account := getAccountFromAccounts(accounts, tag)
	svc.accountCache.Del(serializer.FromID(uid))
	log.Debug("accounts cache delete", zap.Any("user_id", uid))
	if account.SessionToken == status.SessionToken {
		key := serializer.FromID(uid)
		svc.statusCache.Del(key)
		log.Debug("status cache delete", zap.Any("user_id", uid))
		sessionToken := emptySessionToken
		iksm := emptyIKSM
		if len(accounts) > 0 {
			sessionToken = accounts[0].SessionToken
			nintendoAccount, err := svc.nintendoSvc.GetAccountMetadata(sessionToken, language.English)
			if err != nil {
				return errors.Wrap(err, "can't fetch nintendo account metadata")
			}
			iksm = nintendoAccount.IKSM
		}
		err := svc.db.DeleteAndSwitchAccount(uid, tag, sessionToken, iksm)
		if err != nil {
			return errors.Wrap(err, "can't delete the account and switch to a new account in database")
		}
	} else {
		err = svc.db.DeleteAccount(uid, tag)
		if err != nil {
			return errors.Wrap(err, "can't delete the account in database")
		}
	}
	return nil
}

func getAccountFromAccounts(accounts []Account, tag string) Account {
	for _, account := range accounts {
		if account.Tag == tag {
			return account
		}
	}
	return Account{}
}

func isValidAccount(account Account) bool {
	return account.UserID != 0
}

func (svc *serviceImpl) GetAccount(uid ID, tag string) (Account, error) {
	key := serializer.FromID(uid)
	if value := svc.accountCache.Get(key); value != nil {
		accounts := serializer.ToAccounts(value)
		account := getAccountFromAccounts(accounts, tag)
		if isValidAccount(account) {
			return account, nil
		}
	}
	// todo: metrics
	log.Debug("accounts cache miss", zap.Any("user_id", uid))
	account, err := svc.db.SelectAccount(uid, tag)
	if err != nil {
		return account, errors.Wrap(err, "can't fetch account from database")
	}
	return account, nil
}

func (svc *serviceImpl) SwitchAccount(uid ID, tag string) error {
	account, err := svc.GetAccount(uid, tag)
	if err != nil {
		return errors.Wrap(err, "can't find the selected account")
	}
	nintendoAccount, err := svc.nintendoSvc.GetAccountMetadata(account.SessionToken, language.English)
	if err != nil {
		return errors.Wrap(err, "can't fetch nintendo account metadata")
	}
	key := serializer.FromID(uid)
	svc.statusCache.Del(key)
	log.Debug("status cache delete", zap.Any("user_id", uid))
	return svc.db.SwitchAccount(uid, account.SessionToken, nintendoAccount.IKSM)
}

func (svc *serviceImpl) ListAccounts(uid ID) ([]Account, error) {
	key := serializer.FromID(uid)
	if value := svc.accountCache.Get(key); value != nil {
		return serializer.ToAccounts(value), nil
	}
	// todo: metrics
	log.Debug("accounts cache miss", zap.Any("user_id", uid))

	accounts, err := svc.db.SelectAccounts(uid)
	if err != nil {
		return nil, errors.Wrap(err, "can't load accounts from database")
	}

	svc.accountCache.SetExpiration(key, serializer.FromAccounts(accounts), svc.accountsCacheExpiration)
	log.Debug("accounts cache set", zap.Any("user_id", uid), zap.Time("time", time.Now()))
	return accounts, nil
}
