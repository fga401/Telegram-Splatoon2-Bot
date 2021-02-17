package user

import (
	"time"

	"go.uber.org/zap"
	log "telegram-splatoon2-bot/common/log"
	"telegram-splatoon2-bot/driver/cache"
	"telegram-splatoon2-bot/service/language"
	"telegram-splatoon2-bot/service/nintendo"
	"telegram-splatoon2-bot/service/timezone"
	"telegram-splatoon2-bot/service/user/database"
	"telegram-splatoon2-bot/service/user/serializer"
)

type defaultPermission struct {
	Admins       map[ID]struct{}
	MaxAccount   int32
	AllowPolling bool
	Timezone     timezone.Timezone
	Language     language.Language
	IsBlock      bool
}

type serviceImpl struct {
	nintendoSvc   nintendo.Service
	db            database.Service
	adminsCache   cache.IterableCache
	statusCache   cache.Cache
	accountCache  cache.Cache
	proofKeyCache cache.Cache

	defaultPermission       defaultPermission
	accountsCacheExpiration time.Duration
	proofKeyCacheExpiration time.Duration
}

func (svc *serviceImpl) Admins() []ID {
	ret := make([]ID, 0)
	svc.adminsCache.Range(func(key []byte, _ []byte) bool {
		id := serializer.ToID(key)
		ret = append(ret, id)
		return true
	})
	return ret
}

func (svc *serviceImpl) Existed(uid ID) (bool, error) {
	return svc.db.Existed(uid)
}

func New(
	db database.Service,
	adminsCache cache.IterableCache,
	statusCache cache.Cache,
	accountCache cache.Cache,
	proofKeyCache cache.Cache,
	nintendoSvc nintendo.Service,
	config Config,
) Service {
	set := make(map[ID]struct{})
	for _, uid := range config.DefaultPermission.Admins {
		set[uid] = struct{}{}
	}
	svc := &serviceImpl{
		db:            db,
		adminsCache:   adminsCache,
		statusCache:   statusCache,
		accountCache:  accountCache,
		proofKeyCache: proofKeyCache,
		nintendoSvc:   nintendoSvc,

		defaultPermission: defaultPermission{
			Admins:       set,
			MaxAccount:   config.DefaultPermission.MaxAccount,
			AllowPolling: config.DefaultPermission.AllowPolling,
			Timezone:     config.DefaultPermission.Timezone,
			Language:     config.DefaultPermission.Language,
			IsBlock:      config.DefaultPermission.IsBlock,
		},
		accountsCacheExpiration: config.AccountsCacheExpiration,
		proofKeyCacheExpiration: config.ProofKeyCacheExpiration,
	}
	svc.init()
	return svc
}

func (svc *serviceImpl) init() {
	adminIDs, _ := svc.db.Admins()
	for _, adminID := range adminIDs {
		key := serializer.FromID(adminID)
		svc.adminsCache.Set(key, nil)
	}
	validAdmins := svc.Admins()
	log.Info("admins have been loaded.", zap.Any("IDs", validAdmins))
}
