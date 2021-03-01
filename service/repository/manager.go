package repository

import (
	"sync"
	"time"

	"go.uber.org/zap"
	"telegram-splatoon2-bot/common/log"
)

// ManagerConfig sets up a Manager.
type ManagerConfig struct {
	Delay time.Duration
}

// Manager manages all Repository independently.
type Manager interface {
	// Start automatically calls Repository.Update for each Repository and sets a ticker according to Repository.NextUpdateTime.
	Start()
}

type managerImpl struct {
	delay time.Duration
	repos []Repository
	once  sync.Once
}

// NewManager returns a new Manager.
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
