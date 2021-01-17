package stage

import (
	"github.com/pkg/errors"
	nintendo2 "telegram-splatoon2-bot/service/nintendo"
	"telegram-splatoon2-bot/service/repository/utils"
)

//easyjson:json
type stageCollection map[string]*nintendo2.Stage

type DumperImpl struct {
	config DumperConfig
	stages stageCollection
}

func NewDumper(config DumperConfig) *DumperImpl {
	return &DumperImpl{
		config: config,
		stages: make(stageCollection),
	}
}

func (d *DumperImpl) Save() error {
	err := utils.MarshalToFile(d.config.StageFile, d.stages)
	if err != nil {
		return errors.Wrap(err, "can't dump stage")
	}
	return nil
}

func (d *DumperImpl) Load() error {
	err := utils.UnmarshalFromFile(d.config.StageFile, &d.stages)
	if err != nil {
		return errors.Wrap(err, "can't load stage")
	}
	return nil
}

func (d *DumperImpl) Update(src interface{}) error {
	stageSchedules, ok := src.(*nintendo2.StageSchedules)
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
