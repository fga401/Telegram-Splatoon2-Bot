package service

import (
	botapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"sort"
	"strconv"
	"strings"
	log "telegram-splatoon2-bot/logger"
	"telegram-splatoon2-bot/nintendo"
	"time"
)

var salmonSchedules *nintendo.SalmonSchedules
var furtherSalmonScheduleImageID string
var laterSalmonScheduleImageID string

func startSalmonJobScheduler() {
	go func() {
		//first attempt
		err := updateSalmonSchedules()
		if err != nil {
			log.Error("can't update salmon schedules", zap.Error(err))
			return
		}
		// update periodically
		nextUpdateTime := getSplatoonNextUpdateTime(time.Now())
		log.Info("update salmon schedules successfully. start periodical task.", zap.Time("next_update_time", nextUpdateTime))
		for {
			task := time.After(time.Until(nextUpdateTime))
			select {
			case <-task:
				err := updateSalmonSchedules()
				if err != nil {
					nextUpdateTime = time.Now().Add(updateFailureRetryInterval)
					log.Error("can't update salmon schedules", zap.Time("next_update_time", nextUpdateTime), zap.Error(err))
				} else {
					nextUpdateTime = getSplatoonNextUpdateTime(time.Now())
					log.Info("update salmon schedules successfully. set next update task", zap.Time("next_update_time", nextUpdateTime))
				}
			}
		}
	}()
}

func updateSalmonSchedules() error {
	var err error
	admins.Range(func(uid int64) (continued bool) {
		err = updateSalmonSchedulesWithUid(uid)
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
	err = updateSalmonSchedulesWithUid(runtime.Uid)
	if err != nil {
		return errors.Wrap(err, "can't update salmon schedules by other user")
	}
	return nil
}

func updateSalmonSchedulesWithUid(uid int64) error {
	runtime, err := fetchRuntime(uid)
	if err != nil {
		return errors.Wrap(err, "can't fetch runtime")
	}

	var result *nintendo.SalmonSchedules
	var expired bool
	err = retry(func() error {
		result, expired, err = nintendo.GetSalmonSchedules(runtime.IKSM, runtime.Timezone, runtime.Language)
		return err
	}, retryTimes)
	if err != nil {
		return errors.Wrap(err, "can't get salmon schedules")
	}

	if expired {
		// todo: add metric
		iksm, err := updateCookies(runtime)
		if err != nil {
			return errors.Wrap(err, "cookie expired and can't update it")
		}
		runtime.IKSM = iksm
		err = retry(func() error {
			result, expired, err = nintendo.GetSalmonSchedules(runtime.IKSM, runtime.Timezone, runtime.Language)
			return err
		}, retryTimes)
		if err != nil {
			return errors.Wrap(err, "can't get salmon schedules")
		}
	}
	if expired {
		return errors.Errorf("invalid cookie")
	}

	sortSalmonSchedules(result)
	AddPrefixAndEmptyField(result)

	err = uploadSalmonSchedulesImages(result)
	if err != nil {
		return errors.Wrap(err, "can't upload salmon schedules images")
	}
	salmonSchedules = result
	return nil
}

func sortSalmonSchedules(salmonSchedules *nintendo.SalmonSchedules) {
	// sort by start time in descending order
	sort.Slice(salmonSchedules.Details, func(i, j int) bool {
		return salmonSchedules.Details[i].StartTime > salmonSchedules.Details[j].StartTime
	})
	sort.Slice(salmonSchedules.Schedules, func(i, j int) bool {
		return salmonSchedules.Schedules[i].StartTime > salmonSchedules.Schedules[j].StartTime
	})
}

func uploadSalmonSchedulesImages(salmonSchedules *nintendo.SalmonSchedules) error {
	urls := []string{
		salmonSchedules.Details[0].Stage.Image,
		salmonSchedules.Details[0].Weapons[0].Weapon.Image,
		salmonSchedules.Details[0].Weapons[1].Weapon.Image,
		salmonSchedules.Details[0].Weapons[2].Weapon.Image,
		salmonSchedules.Details[0].Weapons[3].Weapon.Image,
		salmonSchedules.Details[1].Stage.Image,
		salmonSchedules.Details[1].Weapons[0].Weapon.Image,
		salmonSchedules.Details[1].Weapons[1].Weapon.Image,
		salmonSchedules.Details[1].Weapons[2].Weapon.Image,
		salmonSchedules.Details[1].Weapons[3].Weapon.Image,
	}
	imgs, err := downloadImages(urls)
	if err != nil {
		return err
	}
	now := strconv.FormatInt(time.Now().Round(time.Hour).Unix(), 10)
	furtherImg := concatSalmonScheduleImage(imgs[0:5])
	laterImg:= concatSalmonScheduleImage(imgs[5:10])
	furtherImgID, err := uploadImage(furtherImg, "further_salmon_schedule_"+now)
	if err != nil {
		return errors.Wrap(err, "can't upload further detail image")
	}
	laterImgID, err := uploadImage(laterImg, "later_salmon_schedule_"+now)
	if err != nil {
		return errors.Wrap(err, "can't upload later detail image")
	}
	furtherSalmonScheduleImageID = furtherImgID
	laterSalmonScheduleImageID = laterImgID
	return nil
}

func initRandomWeapon() {
	randomWeapon.ID = "-1"
	randomWeapon.Name = "Random"
	randomWeapon.Image = "file://./service/resources/salmon_random_weapon_green.png"
	randomWeapon.Thumbnail = "file://./service/resources/salmon_random_weapon_green.png"
}

func AddPrefixAndEmptyField(salmonSchedules *nintendo.SalmonSchedules) {
	for i, detail := range salmonSchedules.Details {
		for j, weapon := range detail.Weapons {
			// todo: distinguish grizzco weapons and normal weapons
			if id, err := strconv.Atoi(weapon.ID); err==nil && id < 0 {
				salmonSchedules.Details[i].Weapons[j].Weapon = randomWeapon
			} else {
				salmonSchedules.Details[i].Weapons[j].Weapon.Image = nintendo.Host +  salmonSchedules.Details[i].Weapons[j].Weapon.Image
				salmonSchedules.Details[i].Weapons[j].Weapon.Thumbnail = nintendo.Host +  salmonSchedules.Details[i].Weapons[j].Weapon.Thumbnail
			}
		}
		salmonSchedules.Details[i].Stage.Image = nintendo.Host + salmonSchedules.Details[i].Stage.Image
	}
}

func QuerySalmonSchedules(update *botapi.Update) error {
	user := update.Message.From
	schedules := salmonSchedules
	if schedules == nil {
		return errors.Errorf("no cached salmonSchedules")
	}
	runtime, err := fetchRuntime(int64(user.ID))
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
	for _, s := range schedules.Schedules {
		startTime := getLocalTime(s.StartTime, runtime.Timezone).Format(timeTemplate)
		endTime := getLocalTime(s.EndTime, runtime.Timezone).Format(timeTemplate)
		keys = append(keys, NewI18nKey(salmonSchedulesScheduleTextKey, startTime, endTime))
	}
	for _, s := range schedules.Details {
		startTime := getLocalTime(s.StartTime, runtime.Timezone).Format(timeTemplate)
		endTime := getLocalTime(s.EndTime, runtime.Timezone).Format(timeTemplate)
		keys = append(keys, NewI18nKey(salmonSchedulesDetailTextKey,
			startTime, endTime, s.Stage.Name,
			s.Weapons[0].Weapon.Name,
			s.Weapons[1].Weapon.Name,
			s.Weapons[2].Weapon.Name,
			s.Weapons[3].Weapon.Name))
	}
	texts := getI18nText(runtime.Language, user, keys...)
	futureText := strings.Join(texts[3:len(schedules.Schedules)+3], "") + texts[0]
	futureMsg := botapi.NewMessage(update.Message.Chat.ID, futureText)
	futureMsg.ParseMode = "Markdown"
	furtherText := texts[3+len(schedules.Schedules)] + texts[1]
	furtherMsg := botapi.NewPhotoShare(update.Message.Chat.ID, furtherSalmonScheduleImageID)
	furtherMsg.Caption = furtherText
	furtherMsg.ParseMode = "Markdown"
	laterText := texts[4+len(schedules.Schedules)] + texts[2]
	laterMsg := botapi.NewPhotoShare(update.Message.Chat.ID, laterSalmonScheduleImageID)
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
