package service

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"image"
	"os"
	"sort"
	"strconv"
	log "telegram-splatoon2-bot/logger"
	"telegram-splatoon2-bot/nintendo"
	"time"
)

var (
	stageScheduleRepo *StageScheduleRepo
)

type StageDumpling struct {
	stageFileName string
	stages        map[string]*nintendo.Stage
}

func NewStageDumpling() *StageDumpling {
	return &StageDumpling{
		stageFileName: viper.GetString("service.stage.stageFileName"),
		stages:        make(map[string]*nintendo.Stage),
	}
}

func (d *StageDumpling) Save() error {
	err := marshalToFile(d.stageFileName, d.stages)
	if err != nil {
		return errors.Wrap(err, "can't save stage")
	}
	return nil
}

func (d *StageDumpling) Load() error {
	if _, err := os.Stat(d.stageFileName); err == nil {
		if err := unmarshalFromFile(d.stageFileName, &d.stages); err != nil {
			return errors.Wrap(err, "can't load stage stage")
		}
	} else {
		log.Warn("can't open stage file", zap.Error(err))
	}
	return nil
}

func (d *StageDumpling) Update(src interface{}) error {
	stageSchedules, ok := src.(*nintendo.StageSchedules)
	if !ok {
		return errors.Errorf("unknown input type")
	}
	for _, stage := range stageSchedules.Regular {
		d.stages[stage.StageA.Name] = stage.StageA
		d.stages[stage.StageB.Name] = stage.StageB
	}
	for _, stage := range stageSchedules.Gachi {
		d.stages[stage.StageA.Name] = stage.StageA
		d.stages[stage.StageB.Name] = stage.StageB
	}
	for _, stage := range stageSchedules.League {
		d.stages[stage.StageA.Name] = stage.StageA
		d.stages[stage.StageB.Name] = stage.StageB
	}
	return nil
}

type GameMode string

const (
	ModeLeague  GameMode = "league"
	ModeGachi   GameMode = "gachi"
	ModeRegular GameMode = "regular"
)

type StageScheduleRepo struct {
	schedules *nintendo.StageSchedules
	imageIDs  map[GameMode][]string

	admins       *SyncUserSet
	stageDumping *StageDumpling
}

func NewStageScheduleRepo(admins *SyncUserSet) (*StageScheduleRepo, error) {
	stageDumping := NewStageDumpling()
	err := stageDumping.Load()
	if err != nil {
		return nil, errors.Wrap(err, "can't load stage dumping files")
	}
	return &StageScheduleRepo{
		admins:        admins,
		stageDumping: stageDumping,
	}, nil
}

func (repo *StageScheduleRepo) HasInit() bool {
	return repo.schedules != nil
}

func (repo *StageScheduleRepo) RepoName() string {
	return "StageScheduleRepo"
}

func (repo *StageScheduleRepo) Update() error {
	var err error
	repo.admins.Range(func(uid int64) (continued bool) {
		err = repo.updateByUid(uid)
		return err != nil
	})
	if err == nil {
		return nil
	}
	log.Warn("can't update stage schedules by admin", zap.Error(err))
	runtime, err := RuntimeTable.GetFirstRuntime()
	if err != nil {
		return errors.Wrap(err, "can't get first runtime object")
	}
	err = repo.updateByUid(runtime.Uid)
	if err != nil {
		return errors.Wrap(err, "can't update stage schedules by other user")
	}
	return nil
}

func (repo *StageScheduleRepo) updateByUid(uid int64) error {
	wrapper := func(iksm string, timezone int, acceptLang string, _ ...interface{}) (interface{}, error) {
		return nintendo.GetStageSchedules(iksm, timezone, acceptLang)
	}
	result, err := fetchResourceWithUpdate(uid, wrapper)
	if err != nil {
		return errors.Wrap(err, "can't fetch runtime")
	}
	schedules := result.(*nintendo.StageSchedules)

	repo.sortSchedules(schedules)
	repo.populateFields(schedules)
	err = repo.stageDumping.Update(schedules)
	if err != nil {
		log.Warn("can't update stage dumping file", zap.Error(err))
	}
	err = repo.stageDumping.Save()
	if err != nil {
		log.Warn("can't update stage dumping file", zap.Error(err))
	} else {
		log.Info("dumped stages to files")
	}

	err = repo.uploadSchedulesImages(schedules)
	if err != nil {
		return errors.Wrap(err, "can't upload stage schedules images")
	}
	repo.schedules = schedules
	return nil
}

func (repo *StageScheduleRepo) sortSchedules(stageSchedules *nintendo.StageSchedules) {
	// sort by start time in descending order
	sort.Slice(stageSchedules.League, func(i, j int) bool {
		return stageSchedules.League[i].StartTime > stageSchedules.League[j].StartTime
	})
	sort.Slice(stageSchedules.League, func(i, j int) bool {
		return stageSchedules.Gachi[i].StartTime > stageSchedules.Gachi[j].StartTime
	})
	sort.Slice(stageSchedules.League, func(i, j int) bool {
		return stageSchedules.Regular[i].StartTime > stageSchedules.Regular[j].StartTime
	})
}

func (repo *StageScheduleRepo) populateFields(stageSchedules *nintendo.StageSchedules) {
	for _, stage := range stageSchedules.Regular {
		stage.StageA.Image = nintendo.Host + stage.StageA.Image
		stage.StageB.Image = nintendo.Host + stage.StageB.Image
	}
	for _, stage := range stageSchedules.Gachi {
		stage.StageA.Image = nintendo.Host + stage.StageA.Image
		stage.StageB.Image = nintendo.Host + stage.StageB.Image
	}
	for _, stage := range stageSchedules.League {
		stage.StageA.Image = nintendo.Host + stage.StageA.Image
		stage.StageB.Image = nintendo.Host + stage.StageB.Image
	}
}

func (repo *StageScheduleRepo) uploadSchedulesImages(stageSchedules *nintendo.StageSchedules) error {
	urls := make([]string, 0, 2*12*3)
	for _, stage := range stageSchedules.Regular {
		urls = append(urls, stage.StageA.Image, stage.StageB.Image)
	}
	for _, stage := range stageSchedules.Gachi {
		urls = append(urls, stage.StageA.Image, stage.StageB.Image)
	}
	for _, stage := range stageSchedules.League {
		urls = append(urls, stage.StageA.Image, stage.StageB.Image)
	}
	imgs, err := downloadImages(urls)
	if err != nil {
		return err
	}
	now := strconv.FormatInt(time.Now().Round(time.Hour).Unix(), 10)
	concatedImgs := make([]image.Image, 0, 12*3)
	concatedImgNames := make([]string, 0, 12*3)
	regularImgs := imgs[0:24]
	gachiImgs := imgs[24:48]
	leagueImgs := imgs[48:72]
	for i := 0; i < len(regularImgs); i += 2 {
		img := concatStageScheduleImage(imgs[i], imgs[i+1])
		concatedImgs = append(concatedImgs, img)
		concatedImgNames = append(concatedImgNames, "regelar_"+strconv.FormatInt(stageSchedules.Regular[i/2].StartTime, 10)+"_"+now)
	}
	for i := 0; i < len(gachiImgs); i += 2 {
		img := concatStageScheduleImage(imgs[i], imgs[i+1])
		concatedImgs = append(concatedImgs, img)
		concatedImgNames = append(concatedImgNames, "gachi_"+strconv.FormatInt(stageSchedules.Gachi[i/2].StartTime, 10)+"_"+now)
	}
	for i := 0; i < len(leagueImgs); i += 2 {
		img := concatStageScheduleImage(imgs[i], imgs[i+1])
		concatedImgs = append(concatedImgs, img)
		concatedImgNames = append(concatedImgNames, "league_"+strconv.FormatInt(stageSchedules.League[i/2].StartTime, 10)+"_"+now)
	}
	ids, err := uploadImages(concatedImgs, concatedImgNames)
	if err != nil {
		return errors.Wrap(err, "can't upload images")
	}
	repo.imageIDs = map[GameMode][]string{
		ModeRegular: ids[0:12],
		ModeGachi:   ids[12:24],
		ModeLeague:  ids[24:36],
	}
	return nil
}
