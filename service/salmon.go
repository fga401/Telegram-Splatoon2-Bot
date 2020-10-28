package service

import (
	botapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"image"
	"os"
	"sort"
	"strconv"
	"strings"
	log "telegram-splatoon2-bot/logger"
	"telegram-splatoon2-bot/nintendo"
	"time"
)

var (
	salmonScheduleRepo *SalmonScheduleRepo
)

type SalmonDumpling struct {
	stageFileName  string
	weaponFileName string
	stages         map[string]*nintendo.SalmonStage
	weapons        map[string]*nintendo.SalmonWeaponWrapper
}

func NewSalmonDumpling() *SalmonDumpling {
	return &SalmonDumpling{
		stageFileName:  viper.GetString("service.salmon.stageFileName"),
		weaponFileName: viper.GetString("service.salmon.weaponFileName"),
		stages:         make(map[string]*nintendo.SalmonStage),
		weapons:        make(map[string]*nintendo.SalmonWeaponWrapper),
	}
}

func (d *SalmonDumpling) Save() error {
	err := DumpingHelper.marshalToFile(d.stageFileName, d.stages)
	if err != nil {
		return errors.Wrap(err, "can't save salmon stage")
	}
	err = DumpingHelper.marshalToFile(d.weaponFileName, d.weapons)
	if err != nil {
		return errors.Wrap(err, "can't save salmon weapon")
	}
	return nil
}

func (d *SalmonDumpling) Load() error {
	if _, err := os.Stat(d.stageFileName); err == nil {
		if err := DumpingHelper.unmarshalFromFile(d.stageFileName, &d.stages); err != nil {
			return errors.Wrap(err, "can't load salmon stage")
		}
	} else {
		log.Warn("can't open salmon stage file", zap.Error(err))
	}
	if _, err := os.Stat(d.weaponFileName); err == nil {
		if err := DumpingHelper.unmarshalFromFile(d.weaponFileName, &d.weapons); err != nil {
			return errors.Wrap(err, "can't load salmon weapons")
		}
	} else {
		log.Warn("can't open salmon weapon file", zap.Error(err))
	}
	return nil
}

func (d *SalmonDumpling) Update(src interface{}) error {
	schedules, ok := src.(*nintendo.SalmonSchedules)
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

type SalmonScheduleRepo struct {
	schedules      *nintendo.SalmonSchedules
	furtherImageID string
	laterImageID   string

	admins        *SyncUserSet
	salmonDumping *SalmonDumpling
}

func NewSalmonScheduleRepo(admins *SyncUserSet) (*SalmonScheduleRepo, error) {
	salmonDumping := NewSalmonDumpling()
	err := salmonDumping.Load()
	if err != nil {
		return nil, errors.Wrap(err, "can't load salmon dumping files")
	}
	return &SalmonScheduleRepo{
		admins:         admins,
		salmonDumping:  salmonDumping,
	}, nil
}

func (repo *SalmonScheduleRepo) HasInit() bool {
	return repo.schedules != nil
}

func (repo *SalmonScheduleRepo) RepoName() string {
	return "SalmonScheduleRepo"
}

func (repo *SalmonScheduleRepo) Update() error {
	var err error
	repo.admins.Range(func(uid int64) (continued bool) {
		err = repo.updateByUid(uid)
		return err != nil
	})
	if err == nil {
		return nil
	}
	log.Warn("can't update salmon schedules by admin", zap.Error(err))
	runtime, err := RuntimeTable.GetFirstRuntime()
	if err != nil {
		return errors.Wrap(err, "can't get first runtime object")
	}
	err = repo.updateByUid(runtime.Uid)
	if err != nil {
		return errors.Wrap(err, "can't update salmon schedules by other user")
	}
	return nil
}

func (repo *SalmonScheduleRepo) updateByUid(uid int64) error {
	wrapper := func(iksm string, timezone int, acceptLang string, _ ...interface{}) (interface{}, error) {
		return nintendo.GetSalmonSchedules(iksm, timezone, acceptLang)
	}
	result, err := FetchResourceWithUpdate(uid, wrapper)
	if err != nil {
		return errors.Wrap(err, "can't fetch runtime")
	}
	schedules := result.(*nintendo.SalmonSchedules)

	repo.sortSchedules(schedules)
	repo.populateFields(schedules)
	err = repo.salmonDumping.Update(schedules)
	if err != nil {
		log.Warn("can't update salmon dumping file", zap.Error(err))
	}
	err = repo.salmonDumping.Save()
	if err != nil {
		log.Warn("can't update salmon dumping file", zap.Error(err))
	} else {
		log.Info("dumped salmon stages and weapons to files")
	}

	err = repo.uploadSchedulesImages(schedules)
	if err != nil {
		return errors.Wrap(err, "can't upload salmon schedules images")
	}
	repo.schedules = schedules
	return nil
}

func (repo *SalmonScheduleRepo) sortSchedules(salmonSchedules *nintendo.SalmonSchedules) {
	// sort by start time in descending order
	sort.Slice(salmonSchedules.Details, func(i, j int) bool {
		return salmonSchedules.Details[i].StartTime > salmonSchedules.Details[j].StartTime
	})
	sort.Slice(salmonSchedules.Schedules, func(i, j int) bool {
		return salmonSchedules.Schedules[i].StartTime > salmonSchedules.Schedules[j].StartTime
	})
}

// populateFields fills Weapon's fields by SpecialWeapon, if Weapon is nil
func (repo *SalmonScheduleRepo) populateFields(salmonSchedules *nintendo.SalmonSchedules) {
	for _, detail := range salmonSchedules.Details {
		for _, weapon := range detail.Weapons {
			if weapon.Weapon != nil {
				weapon.Weapon.Image = nintendo.Host + weapon.Weapon.Image
				weapon.Weapon.Thumbnail = nintendo.Host + weapon.Weapon.Thumbnail
			}
			if weapon.SpecialWeapon != nil {
				weapon.Weapon = &nintendo.SalmonWeapon{
					ID:   weapon.ID,
					Name: weapon.SpecialWeapon.Name,
					// default images are too ugly
					//Image:     nintendo.Host + weapon.SpecialWeapon.Image,
					//Thumbnail: nintendo.Host + weapon.SpecialWeapon.Image,
				}
				if weapon.SpecialWeapon.Name == "Random" {
					weapon.Weapon.Image = "file://./service/resources/salmon_random_weapon_green.png"
					weapon.Weapon.Thumbnail = "file://./service/resources/salmon_random_weapon_green.png"
				} else {
					weapon.Weapon.Image = "file://./service/resources/salmon_random_weapon_yellow.png"
					weapon.Weapon.Thumbnail = "file://./service/resources/salmon_random_weapon_yellow.png"
				}
				weapon.SpecialWeapon = nil
			}
		}
		detail.Stage.Image = nintendo.Host + detail.Stage.Image
	}
}

func (repo *SalmonScheduleRepo) uploadSchedulesImages(salmonSchedules *nintendo.SalmonSchedules) error {
	urls := make([]string, 0, 10)
	for _, detail := range salmonSchedules.Details {
		urls = append(urls, detail.Stage.Image)
		for _, weapon := range detail.Weapons {
			if weapon.Weapon != nil {
				urls = append(urls, weapon.Weapon.Image)
			} else if weapon.SpecialWeapon != nil{
				urls = append(urls, weapon.SpecialWeapon.Image)
			} else {
				return errors.Errorf("no image found")
			}
		}
	}
	imgs, err := ImageHelper.downloadImages(urls)
	if err != nil {
		return err
	}
	now := strconv.FormatInt(time.Now().Round(time.Hour).Unix(), 10)
	furtherImg := concatSalmonScheduleImage(imgs[0:5])
	laterImg := concatSalmonScheduleImage(imgs[5:10])
	ids, err := ImageHelper.uploadImages(
		[]image.Image{furtherImg, laterImg},
		[]string{"further_salmon_schedule_"+now, "later_salmon_schedule_"+now})
	if err != nil {
		return errors.Wrap(err, "can't upload images")
	}
	repo.furtherImageID = ids[0]
	repo.laterImageID = ids[1]
	return nil
}

func QuerySalmonSchedules(update *botapi.Update) error {
	user := update.Message.From
	schedules := salmonScheduleRepo.schedules
	if schedules == nil {
		return errors.Errorf("no cached schedules")
	}
	runtime, err := FetchRuntime(int64(user.ID))
	if err != nil {
		return errors.Wrap(err, "can't fetch runtime")
	}
	now := time.Now().Unix()
	startTime := schedules.Details[1].StartTime
	endTime := schedules.Details[1].EndTime
	var textKey string
	var remainingTime time.Duration
	if now > startTime {
		textKey = salmonSchedulesOpenTextKey
		remainingTime = time.Until(time.Unix(endTime, 0))
	} else {
		textKey = salmonSchedulesSoonTextKey
		remainingTime = time.Until(time.Unix(startTime, 0))
	}
	remainingTime = remainingTime.Round(time.Minute)
	hour := remainingTime / time.Hour
	remainingTime -= hour * time.Hour
	minute := remainingTime / time.Minute

	timeTemplate := getI18nText(runtime.Language, user, NewI18nKey(TimeTemplateTextKey))[0]
	keys := []I18nKeys{
		NewI18nKey(salmonSchedulesFutureTextKey),
		NewI18nKey(salmonSchedulesNextTextKey),
		NewI18nKey(textKey, hour, minute),
	}
	future := schedules.Schedules[:len(schedules.Schedules) - 2]
	for _, s := range future {
		startTime := TimeHelper.getLocalTime(s.StartTime, runtime.Timezone).Format(timeTemplate)
		endTime := TimeHelper.getLocalTime(s.EndTime, runtime.Timezone).Format(timeTemplate)
		keys = append(keys, NewI18nKey(salmonSchedulesScheduleTextKey, startTime, endTime))
	}
	for _, s := range schedules.Details {
		startTime := TimeHelper.getLocalTime(s.StartTime, runtime.Timezone).Format(timeTemplate)
		endTime := TimeHelper.getLocalTime(s.EndTime, runtime.Timezone).Format(timeTemplate)
		keys = append(keys, NewI18nKey(salmonSchedulesDetailTextKey,
			startTime, endTime, s.Stage.Name,
			s.Weapons[0].Weapon.Name,
			s.Weapons[1].Weapon.Name,
			s.Weapons[2].Weapon.Name,
			s.Weapons[3].Weapon.Name))
	}
	texts := getI18nText(runtime.Language, user, keys...)
	futureText := strings.Join(texts[3:len(future)+3], "") + texts[0]
	futureMsg := botapi.NewMessage(update.Message.Chat.ID, futureText)
	futureMsg.ParseMode = "Markdown"
	furtherText := texts[3+len(future)] + texts[1]
	furtherMsg := botapi.NewPhotoShare(update.Message.Chat.ID, salmonScheduleRepo.furtherImageID)
	furtherMsg.Caption = furtherText
	furtherMsg.ParseMode = "Markdown"
	laterText := texts[4+len(future)] + texts[2]
	laterMsg := botapi.NewPhotoShare(update.Message.Chat.ID, salmonScheduleRepo.laterImageID)
	laterMsg.Caption = laterText
	laterMsg.ParseMode = "Markdown"
	err = sendWithRetry(bot, futureMsg)
	if err != nil {
		return err
	}
	err = sendWithRetry(bot, furtherMsg)
	if err != nil {
		return err
	}
	err = sendWithRetry(bot, laterMsg)
	if err != nil {
		return err
	}
	return nil
}
