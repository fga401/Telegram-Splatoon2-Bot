package user

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
	log "telegram-splatoon2-bot/common/log"
	"telegram-splatoon2-bot/service/language"
	"telegram-splatoon2-bot/service/timezone"
	"telegram-splatoon2-bot/service/user/internal/serializer"
)

func (svc *serviceImpl) GetStatus(uid ID) (Status, error) {
	// fetch from cache
	key := serializer.FromID(uid)
	if value := svc.statusCache.Get(key); value != nil {
		return serializer.ToStatus(value), nil
	}
	// todo: metrics
	log.Debug("status cache miss", zap.Any("user_id", uid))

	// fetch from database
	status, err := svc.db.SelectStatus(uid)
	if err != nil {
		return status, errors.Wrap(err, "can't load status from database")
	}

	// set cache
	svc.statusCache.Set(key, serializer.FromStatus(status))
	log.Debug("status cache set", zap.Any("user_id", uid))
	return status, nil
}

func (svc *serviceImpl) UpdateStatusIKSM(uid ID) (Status, error) {
	status, err := svc.GetStatus(uid)
	if err != nil {
		return Status{}, errors.Wrap(err, "can't fetch status")
	}
	metadata, err := svc.nintendoSvc.GetAccountMetadata(status.SessionToken, language.English)
	if err != nil {
		return Status{}, errors.Wrap(err, "can't get account metadata")
	}
	err = svc.db.UpdateStatusIKSM(uid, metadata.IKSM)
	if err != nil {
		return Status{}, errors.Wrap(err, "can't update status IKSM in database")
	}
	svc.statusCache.Del(serializer.FromID(uid))
	log.Debug("status cache delete", zap.Any("user_id", uid))
	return svc.GetStatus(uid)
}

func (svc *serviceImpl) UpdateStatusTimezone(uid ID, timezone timezone.Timezone) (Status, error) {
	err := svc.db.UpdateStatusTimezone(uid, timezone)
	if err != nil {
		return Status{}, errors.Wrap(err, "can't update status timezone in database")
	}
	svc.statusCache.Del(serializer.FromID(uid))
	log.Debug("status cache delete", zap.Any("user_id", uid))
	return svc.GetStatus(uid)
}

func (svc *serviceImpl) UpdateStatusLanguage(uid ID, language language.Language) (Status, error) {
	err := svc.db.UpdateStatusLanguage(uid, language)
	if err != nil {
		return Status{}, errors.Wrap(err, "can't update status language in database")
	}
	svc.statusCache.Del(serializer.FromID(uid))
	log.Debug("status cache delete", zap.Any("user_id", uid))
	return svc.GetStatus(uid)
}

func (svc *serviceImpl) UpdateStatusLastBattle(uid ID, lastBattle string) (Status, error) {
	err := svc.db.UpdateStatusLastBattle(uid, lastBattle)
	if err != nil {
		return Status{}, errors.Wrap(err, "can't update status lastBattle in database")
	}
	svc.statusCache.Del(serializer.FromID(uid))
	log.Debug("status cache delete", zap.Any("user_id", uid))
	return svc.GetStatus(uid)
}

func (svc *serviceImpl) UpdateStatusLastSalmon(uid ID, lastSalmon string) (Status, error) {
	err := svc.db.UpdateStatusLastSalmon(uid, lastSalmon)
	if err != nil {
		return Status{}, errors.Wrap(err, "can't update status lastSalmon in database")
	}
	svc.statusCache.Del(serializer.FromID(uid))
	log.Debug("status cache delete", zap.Any("user_id", uid))
	return svc.GetStatus(uid)
}
