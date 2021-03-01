package repository

import (
	"sync"
	"time"

	"go.uber.org/zap"
	"telegram-splatoon2-bot/common/log"
)

type ManagerConfig struct {
	Delay time.Duration
}

type Manager interface {
	Start()
}

type managerImpl struct {
	delay time.Duration
	repos []Repository
	once  sync.Once
}

func NewManager(config ManagerConfig, repos ...Repository) Manager {
	return &managerImpl{
		repos: repos,
		delay: config.Delay,
	}
}

func (m *managerImpl) Start() {
	m.once.Do(func() {
		for _, repo := range m.repos {
			m.start(repo)
		}
	})
}

func (m *managerImpl) start(repo Repository) {
	go func() {
		name := repo.Name()
		for {
			nextUpdateTime := repo.NextUpdateTime()
			task := time.After(time.Until(nextUpdateTime))
			select {
			case <-task:
				err := repo.Update()
				if err != nil {
					log.Error("can't update Repository: "+name, zap.Time("next_update_time", nextUpdateTime), zap.Error(err))
					<-time.Tick(m.delay)
				} else {
					log.Info("update Repository: "+name, zap.Time("next_update_time", nextUpdateTime))
				}
			}
		}
	}()
}
