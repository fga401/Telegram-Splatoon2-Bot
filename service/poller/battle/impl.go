package battle

import (
	"time"

	"github.com/pkg/errors"
	"telegram-splatoon2-bot/common/queue"
	"telegram-splatoon2-bot/common/util"
	"telegram-splatoon2-bot/service/language"
	"telegram-splatoon2-bot/service/nintendo"
	"telegram-splatoon2-bot/service/repository/stage"
	"telegram-splatoon2-bot/service/user"
	"telegram-splatoon2-bot/telegram/bot"
)

type impl struct {
	bot         bot.Bot
	nintendoSvc nintendo.Service
	userSvc     user.Service
	repository  stage.Repository

	currentMinBattleTime time.Duration
	refreshTime          time.Duration
	maxIdleTime          time.Duration
	maxWorker            int32
	minBattleTime        MinBattleTime

	runningTasks map[user.ID]*statistics
	resultQueue  queue.Queue
	toFetchQueue chan user.ID
	restartQueue queue.Queue
	refreshQueue queue.Queue
	startChan    chan user.ID
	stopChan     chan user.ID
	outChan      chan Result
	outQueue     queue.Queue
}

// New returns a battle poller object.
func New(
	bot bot.Bot,
	repository stage.Repository,
	nintendoSvc nintendo.Service,
	userSvc user.Service,
	config Config,
) Service {
	svc := &impl{
		bot:         bot,
		repository:  repository,
		nintendoSvc: nintendoSvc,
		userSvc:     userSvc,

		refreshTime:   config.RefreshmentTime,
		maxWorker:     config.MaxWorker,
		maxIdleTime:   config.MaxIdleTime,
		minBattleTime: config.MinBattleTime,

		runningTasks: make(map[user.ID]*statistics),
		resultQueue:  queue.New(),
		toFetchQueue: make(chan user.ID),
		restartQueue: queue.New(),
		refreshQueue: queue.New(),
		startChan:    make(chan user.ID),
		stopChan:     make(chan user.ID),
		outChan:      make(chan Result),
		outQueue:     queue.New(),
	}
	go svc.dispatchRoutine()
	go svc.fetchingRoutine()
	go svc.statisticsManagementRoutine()
	go svc.returnRoutine()
	return svc
}

func (svc *impl) Results() <-chan Result {
	return svc.outChan
}

func (svc *impl) Start(id user.ID) {
	go func() {
		svc.startChan <- id
	}()
}

func (svc *impl) Stop(id user.ID) {
	go func() {
		svc.stopChan <- id
	}()
}

func (svc *impl) fetchingRoutine() {
	if svc.maxWorker > 0 {
		for i := int32(0); i < svc.maxWorker; i++ {
			go func() {
				for id := range svc.toFetchQueue {
					result := svc.fetch(id)
					svc.resultQueue.EnqueueChan() <- result
				}
			}()
		}
	} else {
		for id := range svc.toFetchQueue {
			go func(id user.ID) {
				result := svc.fetch(id)
				svc.resultQueue.EnqueueChan() <- result
			}(id)
		}
	}
}

func (svc *impl) dispatchRoutine() {
	restartTimer := time.NewTimer(0)
	refreshTimer := time.NewTimer(0)
	<-restartTimer.C
	<-refreshTimer.C
	var restartID, refreshID user.ID = 0, 0
	dequeueOrWaitTimer := func(q queue.Queue, id user.ID) <-chan interface{} {
		if id == 0 {
			return q.DequeueChan()
		}
		return nil
	}
	for {
		select {
		case <-restartTimer.C:
			svc.toFetchQueue <- restartID
			restartID = 0
		case <-refreshTimer.C:
			svc.toFetchQueue <- refreshID
			refreshID = 0
		case taskRaw := <-dequeueOrWaitTimer(svc.restartQueue, restartID):
			task := taskRaw.(task)
			restartID = task.UserID
			restartTimer.Reset(time.Until(task.UpdateTime.Add(svc.currentMinBattleTime)))
		case taskRaw := <-dequeueOrWaitTimer(svc.refreshQueue, refreshID):
			task := taskRaw.(task)
			refreshID = task.UserID
			refreshTimer.Reset(time.Until(task.UpdateTime.Add(svc.refreshTime)))
		}
	}
}

func (svc *impl) statisticsManagementRoutine() {
	for {
		select {
		case id := <-svc.startChan:
			if _, ok := svc.runningTasks[id]; ok {
				return
			}
			svc.start(id)
		case id := <-svc.stopChan:
			svc.stop(id)
		case resultRaw := <-svc.resultQueue.DequeueChan():
			result := resultRaw.(Result)
			if stat, ok := svc.runningTasks[result.UserID]; ok {
				if svc.doCancel(stat, result) {
					svc.stop(result.UserID)
					result.Error = &ErrCanceledPolling{Reason: CancelReasonEnum.NoNewBattles}
					svc.outQueue.EnqueueChan() <- result
					continue
				}
				if isValidResult(result) && isDifferentFromLastBattle(stat, result) {
					stat.LastBattle = result.Battles[0]
					updateTime := time.Unix(stat.LastBattle.EndTime(), 0)
					if stat.LastBattle == nil {
						updateTime = time.Now()
					}
					svc.restartQueue.EnqueueChan() <- task{
						UserID:     result.UserID,
						UpdateTime: updateTime,
					}
					svc.outQueue.EnqueueChan() <- result
					continue
				}
				svc.refreshQueue.EnqueueChan() <- task{
					UserID:     result.UserID,
					UpdateTime: time.Now(),
				}
			}
		}
	}
}

func isValidResult(result Result) bool {
	return len(result.Battles) > 0 && result.Error == nil
}

func isDifferentFromLastBattle(stat *statistics, result Result) bool {
	return stat.LastBattle == nil || (stat.LastBattle != nil && stat.LastBattle.Metadata().BattleNumber != result.Battles[0].Metadata().BattleNumber)
}

func (svc *impl) stop(id user.ID) {
	delete(svc.runningTasks, id)
}

func (svc *impl) start(id user.ID) {
	svc.runningTasks[id] = &statistics{
		LastBattle: nil,
		CreateTime: time.Now(),
	}
	svc.refreshQueue.EnqueueChan() <- task{
		UserID:     id,
		UpdateTime: time.Now(),
	}
}

func (svc *impl) doCancel(stat *statistics, result Result) bool {
	lastUpdateTime := stat.CreateTime
	if stat.LastBattle != nil {
		lastUpdateTime = time.Unix(stat.LastBattle.EndTime(), 0)
	}
	return !isValidResult(result) && time.Since(lastUpdateTime) > svc.maxIdleTime
}

func (svc *impl) fetch(id user.ID) Result {
	status, err := svc.userSvc.GetStatus(id)
	if err != nil {
		return Result{
			UserID: id,
			Error:  err,
		}
	}
	battles, err := svc.nintendoSvc.GetLatestBattleResults(status.LastBattle, 0, status.IKSM, status.Timezone, language.English)
	if errors.Is(err, &nintendo.ErrIKSMExpired{}) {
		status, err = svc.userSvc.UpdateStatusIKSM(status.UserID)
		if err != nil {
			return Result{
				UserID: id,
				Error:  err,
			}
		}
		battles, err = svc.nintendoSvc.GetLatestBattleResults(status.LastBattle, 0, status.IKSM, status.Timezone, language.English)
	}
	if err != nil {
		return Result{
			UserID: id,
			Error:  err,
		}
	}
	var detail nintendo.DetailedBattleResult
	if len(battles) == 1 {
		detail, err = svc.nintendoSvc.GetDetailedBattleResults(battles[0].Metadata().BattleNumber, status.IKSM, status.Timezone, language.English)
	}
	return Result{
		UserID:  id,
		Battles: battles,
		Detail:  detail,
		Error:   err,
	}
}

func (svc *impl) returnRoutine() {
	for result := range svc.outQueue.DequeueChan() {
		svc.outChan <- result.(Result)
	}
}

func (svc *impl) updateMinBattleTimeRoutine() {
	ticker := time.NewTicker(0)
	primary := stage.NewPrimaryFilter([]stage.Mode{stage.ModeEnum.Gachi, stage.ModeEnum.League})
	secondary := []stage.SecondaryFilter{stage.NewNextNSecondaryFilter(2)}
	for range ticker.C {
		// todo: Private Battle?
		schedules := svc.repository.Content(primary, secondary, 4) // Gachi[0], League[0], Gachi[1], League[1]
		if len(schedules) < 4 {
			svc.currentMinBattleTime = minDuration(svc.minBattleTime.Waiting, svc.minBattleTime.Zone, svc.minBattleTime.Clam, svc.minBattleTime.Tower, svc.minBattleTime.Rainmaker) + svc.minBattleTime.Waiting
		} else if schedules[0].Schedule.EndTime > time.Now().Unix() {
			svc.currentMinBattleTime = minDuration(svc.ruleToDuration(schedules[0].Schedule.Rule.Key), svc.ruleToDuration(schedules[1].Schedule.Rule.Key)) + svc.minBattleTime.Waiting
		} else {
			svc.currentMinBattleTime = minDuration(svc.ruleToDuration(schedules[2].Schedule.Rule.Key), svc.ruleToDuration(schedules[3].Schedule.Rule.Key)) + svc.minBattleTime.Waiting
		}
		ticker.Reset(time.Until(util.Time.SplatoonNextUpdateTime(time.Now())))
	}
}

func minDuration(durations ...time.Duration) time.Duration {
	if len(durations) == 0 {
		return 0
	}
	min := durations[0]
	for i := 1; i < len(durations); i++ {
		if durations[i] < min {
			min = durations[i]
		}
	}
	return min
}

func (svc *impl) ruleToDuration(rule string) time.Duration {
	switch rule {
	case nintendo.KeyTurfWar:
		return 3 * time.Minute
	case nintendo.KeyClamBlitz:
		return svc.minBattleTime.Clam
	case nintendo.KeyTowerControl:
		return svc.minBattleTime.Tower
	case nintendo.KeySplatZones:
		return svc.minBattleTime.Zone
	case nintendo.KeyRainmaker:
		return svc.minBattleTime.Rainmaker
	default:
		return svc.minBattleTime.Waiting
	}
}
