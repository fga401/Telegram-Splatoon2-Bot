package user

import (
	"github.com/pkg/errors"
	"telegram-splatoon2-bot/service/user/internal/serializer"
)

const (
	emptyIKSM         = "0000000000000000000000000000000000000000"
	emptySessionToken = ""
)

func (svc *serviceImpl) Register(uid ID, username string) error {
	_, isAdmin := svc.defaultPermission.Admins[uid]
	permission := Permission{
		UserID:       uid,
		IsBlock:      svc.defaultPermission.IsBlock,
		MaxAccount:   svc.defaultPermission.MaxAccount,
		IsAdmin:      isAdmin,
		AllowPolling: svc.defaultPermission.AllowPolling,
	}
	status := Status{
		UserID:       uid,
		SessionToken: emptySessionToken,
		IKSM:         emptyIKSM,
		Language:     svc.defaultPermission.Language,
		Timezone:     svc.defaultPermission.Timezone,
	}
	user := User{
		UserID:   uid,
		UserName: username,
	}
	err := svc.db.Register(user, permission, status)
	if err != nil {
		return errors.Wrap(err, "can't insert user and status to db")
	}

	key := serializer.FromID(uid)
	value := serializer.FromStatus(status)
	svc.statusCache.Set(key, value)
	if isAdmin {
		svc.adminsCache.Set(key, nil)
	}
	return nil
}
