package service

import (
	botapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"sort"
	"strings"
	log "telegram-splatoon2-bot/logger"
	"telegram-splatoon2-bot/nintendo"
	"time"
)

var salmonSchedules *nintendo.SalmonSchedules
var furtherSalmonScheduleImageID string
var laterSalmonScheduleImageID string

func startSalmonJobScheduler() {
	// todo: debug: will be failed after first user register
	go func() {
		//first attempt
		err := updateSalmonSchedules()
		if err != nil {
			log.Error("can't update salmon schedules", zap.Error(err))
			return
		}
		// update periodically
		for {
			now := time.Now()
			nextUpdateTime := getSplatoonNextUpdateTime(now)
			task := time.After(nextUpdateTime.Sub(now))
			select {
			case <-task:
				err := updateSalmonSchedules()
				if err != nil {
					log.Error("can't update salmon schedules")
				}
			}
		}
	}()
}

func updateSalmonSchedules() error {
	for k, _ := range admins {
		err := updateSalmonSchedulesWithUid(k)
		if err == nil {
			return nil
		}
	}
	log.Warn("can't update salmon schedules by admin")
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
	furtherImg, err:= concatSalmonScheduleImage(&salmonSchedules.Details[0])
	if err != nil {
		return errors.Wrap(err, "can't prepare image")
	}
	laterImg, err:= concatSalmonScheduleImage(&salmonSchedules.Details[1])
	if err != nil {
		return errors.Wrap(err, "can't prepare image")
	}
	furtherImgID, err := uploadImage(furtherImg, "further")
	laterImgID, err := uploadImage(laterImg, "later")
	furtherSalmonScheduleImageID = furtherImgID
	laterSalmonScheduleImageID = laterImgID
	return nil
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
	// todo: show time to start / to end
	keys := []I18nKeys{
		{salmonSchedulesFutureTextKey, nil},
		{salmonSchedulesNextTextKey, []interface{}{nintendo.Host + schedules.Details[0].Stage.Image}},
		{salmonSchedulesOpenTextKey, []interface{}{nintendo.Host + schedules.Details[1].Stage.Image}},
	}
	for _, s := range schedules.Schedules {
		startTime := time.Unix(s.StartTime, 0).Format("01-02 15:04")
		endTime := time.Unix(s.EndTime, 0).Format("01-02 15:04")
		keys = append(keys, I18nKeys{
			salmonSchedulesScheduleTextKey,
			[]interface{}{startTime, endTime}})
	}
	for _, s := range schedules.Details {
		startTime := time.Unix(s.StartTime, 0).Format("01-02 15:04")
		endTime := time.Unix(s.EndTime, 0).Format("01-02 15:04")
		keys = append(keys, I18nKeys{
			salmonSchedulesDetailTextKey,
			[]interface{}{startTime, endTime, s.Stage.Name,
				s.Weapons[0].Weapon.Name,
				s.Weapons[1].Weapon.Name,
				s.Weapons[2].Weapon.Name,
				s.Weapons[3].Weapon.Name}})
	}
	texts := getI18nText(runtime.Language, user, keys...)
	futureText := strings.Join(texts[3:len(schedules.Schedules) + 3], "") + texts[0]
	futureMsg := botapi.NewMessage(update.Message.Chat.ID, futureText)
	futureMsg.ParseMode = "Markdown"
	furtherText := texts[3+len(schedules.Schedules)] + texts[1]
	furtherMsg := botapi.NewPhotoShare(update.Message.Chat.ID, furtherSalmonScheduleImageID)
	furtherMsg.Caption= furtherText
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
