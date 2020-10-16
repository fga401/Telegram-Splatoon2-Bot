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

func tryStartJobSchedulers() {
	if !salmonScheduleRepo.HasInit() {
		log.Info("start salmon job scheduler")
		startJobScheduler(salmonScheduleRepo)
	}
	// todo: stage schedule repo
	//if !stageScheduleRepo.HasInit() {
	//	log.Info("start stage job scheduler")
	//	startJobScheduler(stageScheduleRepo)
	//}
}

var updateInterval = int64(2 * time.Hour.Seconds())

func getSplatoonNextUpdateTime(t time.Time) time.Time {
	nowTimestamp := t.Unix()
	nextTimestamp := (nowTimestamp/updateInterval + 1) * updateInterval
	nextTimestamp += updateDelayInSecond
	return time.Unix(nextTimestamp, 0)
}

func startJobScheduler(repo Repo) {
	name := repo.RepoName()
	go func() {
		//first attempt
		err := repo.Update()
		if err != nil {
			log.Error(name+": can't update", zap.Error(err))
			return
		}
		// update periodically
		nextUpdateTime := getSplatoonNextUpdateTime(time.Now())
		log.Info(name+": update successfully. start periodical task.", zap.Time("next_update_time", nextUpdateTime))
		for {
			task := time.After(time.Until(nextUpdateTime))
			select {
			case <-task:
				err := repo.Update()
				if err != nil {
					nextUpdateTime = time.Now().Add(updateFailureRetryInterval)
					log.Error(name+": can't update", zap.Time("next_update_time", nextUpdateTime), zap.Error(err))
				} else {
					nextUpdateTime = getSplatoonNextUpdateTime(time.Now())
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

func marshalToFile(fileName string, obj interface{}) error {
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

func unmarshalFromFile(fileName string, obj interface{}) error {
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
