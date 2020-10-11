package service

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"sync"
	log "telegram-splatoon2-bot/logger"
)

type UserSet map[int64]struct{}
type SyncUserSet struct {
	set UserSet
	mtx sync.RWMutex
}

func NewUserSet(set UserSet) *SyncUserSet {
	return &SyncUserSet{
		set: set,
		mtx: sync.RWMutex{},
	}
}

func (set *SyncUserSet) Existed(uid int64) bool {
	set.mtx.RLock()
	defer set.mtx.RUnlock()
	_, found := set.set[uid]
	return found
}

func (set *SyncUserSet) Add(uid int64) {
	set.mtx.Lock()
	defer set.mtx.Unlock()
	set.set[uid] = struct{}{}
}

func (set *SyncUserSet) Del(uid int64) {
	set.mtx.Lock()
	defer set.mtx.Unlock()
	delete(set.set, uid)
}

func (set *SyncUserSet) Range(f func(int64) bool){
	set.mtx.RLock()
	for uid, _ := range set.set{
		set.mtx.RUnlock()
		continued := f(uid)
		set.mtx.RLock()
		if !continued {
			break
		}
	}
	set.mtx.RUnlock()
}

var (
	admins *SyncUserSet
	// allowPolling map[int64]struct{}
	// isBlock map[int64]struct{}
)

func loadUsers() {
	adminsList, err := UserTable.LoadAdmin()
	if err != nil {
		panic(errors.Wrap(err, "can't load admins"))
	}
	adminsMap := make(UserSet)
	for _, id := range adminsList {
		adminsMap[id] = struct{}{}
	}
	log.Info("loaded admins list", zap.Int64s("admins", adminsList))
	admins = NewUserSet(adminsMap)
	// todo: load block and polling
}
