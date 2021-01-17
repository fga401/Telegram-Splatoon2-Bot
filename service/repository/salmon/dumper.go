package salmon

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
	log "telegram-splatoon2-bot/common/log"
	nintendo2 "telegram-splatoon2-bot/service/nintendo"
	"telegram-splatoon2-bot/service/repository/utils"
)

//easyjson:json
type stageCollection map[string]*nintendo2.SalmonStage

//easyjson:json
type weaponCollection map[string]*nintendo2.SalmonWeaponType

type DumperImpl struct {
	config  DumperConfig
	stages  stageCollection
	weapons weaponCollection
}

func NewDumper(config DumperConfig) *DumperImpl {
	return &DumperImpl{
		config:  config,
		stages:  make(stageCollection),
		weapons: make(weaponCollection),
	}
}

func (d *DumperImpl) Save() error {
	err := utils.MarshalToFile(d.config.StageFile, d.stages)
	if err != nil {
		return errors.Wrap(err, "can't dump salmon stage")
	}
	err = utils.MarshalToFile(d.config.WeaponFile, d.weapons)
	if err != nil {
		return errors.Wrap(err, "can't dump salmon weapon")
	}
	return nil
}

func (d *DumperImpl) Load() error {
	err := utils.UnmarshalFromFile(d.config.StageFile, &d.stages)
	if err != nil {
		log.Warn("can't load salmon stage", zap.Error(err))
	}
	err = utils.UnmarshalFromFile(d.config.WeaponFile, &d.weapons)
	if err != nil {
		log.Warn("can't load salmon weapons", zap.Error(err))
	}
	return nil
}

func (d *DumperImpl) Update(obj interface{}) error {
	schedules, ok := obj.(*nintendo2.SalmonSchedules)
	if !ok {
		return errors.Errorf("unknown input type")
	}
	for _, detail := range schedules.Details {
		d.stages[detail.Stage.Name] = detail.Stage
		for _, weapon := range detail.Weapons {
			if weapon.Weapon != nil {
				d.weapons[weapon.Weapon.Name] = weapon
			} else if weapon.SpecialWeapon != nil {
				d.weapons[weapon.SpecialWeapon.Name] = weapon
			}
		}
	}
	return nil
}
