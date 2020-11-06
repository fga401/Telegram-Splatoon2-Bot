package service

import (
	json "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"io/ioutil"
	log "telegram-splatoon2-bot/logger"
	"time"
)

type Repo interface {
	HasInit() bool
	RepoName() string
	Update() error
}

type scheduler struct{}

var Scheduler scheduler

func (s scheduler) tryStart() {
	if salmonScheduleRepo != nil && !salmonScheduleRepo.HasInit() {
		log.Info("start salmon job scheduler")
		s.start(salmonScheduleRepo)
	}
	if stageScheduleRepo != nil && !stageScheduleRepo.HasInit() {
		log.Info("start stage job scheduler")
		s.start(stageScheduleRepo)
	}
}

func (scheduler) start(repo Repo) {
	delay := time.Duration(updateDelayInSecond) * time.Second
	name := repo.RepoName()
	go func() {
		nextUpdateTime := time.Now()
		for {
			task := time.After(time.Until(nextUpdateTime))
			select {
			case <-task:
				err := repo.Update()
				if err != nil {
					nextUpdateTime = time.Now().Add(updateFailureRetryInterval)
					log.Error(name+": can't update", zap.Time("next_update_time", nextUpdateTime), zap.Error(err))
				} else {
					nextUpdateTime = TimeHelper.getSplatoonNextUpdateTime(time.Now()).Add(delay)
					log.Info(name+": update successfully. set next update task", zap.Time("next_update_time", nextUpdateTime))
				}
			}
		}
	}()
}

type Dumping interface {
	Update(src interface{}) error
	Load() error
	Save() error
}

type dumpingHelper struct{}

var DumpingHelper dumpingHelper

func (dumpingHelper) marshalToFile(fileName string, obj interface{}) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return errors.Wrap(err, "can't marshal object")
	}
	err = ioutil.WriteFile(fileName, data, 0644)

	if err != nil {
		return errors.Wrap(err, "can't write object to file:"+fileName)
	}
	return nil
}

func (dumpingHelper) unmarshalFromFile(fileName string, obj interface{}) error {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return errors.Wrap(err, "can't read object to file:"+fileName)
	}
	err = json.Unmarshal(data, obj)
	if err != nil {
		return errors.Wrap(err, "can't unmarshal object")
	}
	return nil
}
