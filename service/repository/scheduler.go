package repository

type Repository interface {
	HasInit() bool
	Name() string
	Update() error
}

type scheduler struct{

}

var Scheduler scheduler

func (s scheduler) tryStart() {
	//if salmon.salmonScheduleRepo != nil && !salmon.salmonScheduleRepo.HasInit() {
	//	log.Info("start salmon job scheduler")
	//	s.start(salmon.salmonScheduleRepo)
	//}
	//if stage.stageScheduleRepo != nil && !stage.stageScheduleRepo.HasInit() {
	//	log.Info("start stage job scheduler")
	//	s.start(stage.stageScheduleRepo)
	//}
}

func (scheduler) start(repo Repository) {
	//delay := time.Duration(service.updateDelayInSecond) * time.Second
	//name := repo.Name()
	//go func() {
	//	nextUpdateTime := time.Now()
	//	for {
	//		task := time.After(time.Until(nextUpdateTime))
	//		select {
	//		case <-task:
	//			err := repo.Update()
	//			if err != nil {
	//				nextUpdateTime = time.Now().Add(service.updateFailureRetryInterval)
	//				log.Error(name+": can't update", zap.Time("next_update_time", nextUpdateTime), zap.Error(err))
	//			} else {
	//				nextUpdateTime = service.TimeHelper.getSplatoonNextUpdateTime(time.Now()).Add(delay)
	//				log.Info(name+": update successfully. set next update task", zap.Time("next_update_time", nextUpdateTime))
	//			}
	//		}
	//	}
	//}()
}
