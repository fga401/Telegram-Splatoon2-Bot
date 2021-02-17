package repository

import (
	"strings"
	"time"

	botApi "github.com/go-telegram-bot-api/telegram-bot-api"
	"golang.org/x/text/message"
	"telegram-splatoon2-bot/common/util"
	"telegram-splatoon2-bot/service/repository/salmon"
	"telegram-splatoon2-bot/service/timezone"
	userSvc "telegram-splatoon2-bot/service/user"
	"telegram-splatoon2-bot/telegram/controller/internal/adapter"
	botMessage "telegram-splatoon2-bot/telegram/controller/internal/message"
)

func (ctrl *repositoryCtrl) salmon(update botApi.Update, argManager adapter.Manager, args ...interface{}) error {
	statusArgIdx := argManager.Index(ctrl.statusAdapter)[0]
	status := args[statusArgIdx].(userSvc.Status)
	content := ctrl.salmonRepo.Content()
	if content == nil {
		msg := getSalmonSchedulesNoReadyMessage(ctrl.languageSvc.Printer(status.Language), update)
		_, err := ctrl.bot.Send(msg)
		return err
	}
	msgs:= getSalmonSchedulesMessages(ctrl.languageSvc.Printer(status.Language), update, content, status.Timezone)
	for _, msg := range msgs {
		_, err := ctrl.bot.Send(msg)
		if err != nil {
			return err
		}
	}
	return nil
}

const (
	textKeySalmonSchedulesNoReady   = "Salmon schedules have not been ready yet."

	textKeySalmonSchedulesFutureTag = "#Future"
	textKeySalmonSchedulesNextTag   = "#Next"
	textKeySalmonSchedulesOpenTag   = "#Open: *Will be over in %dh %dm!*"
	textKeySalmonSchedulesSoonTag   = "#Soon: *Will start in %dh %dm!*"
	textKeySalmonSchedulesSchedule  = "*Time*: `%s ~ %s`\n"
	textKeySalmonSchedulesDetail    = "*Time*: `%s ~ %s`\n*Stage*: %s\n*Weapons*:\n- %s\n- %s\n- %s\n- %s\n"
)

func getSalmonSchedulesNoReadyMessage(printer *message.Printer, update botApi.Update) botApi.Chattable {
	text := printer.Sprintf(textKeySalmonSchedulesNoReady)
	return botMessage.NewByUpdate(update, text, nil)
}

func getSalmonSchedulesMessages(printer *message.Printer, update botApi.Update, content *salmon.Content, timezone timezone.Timezone) []botApi.Chattable {
	timeTemplate := printer.Sprintf(timeTemplateTextKey)
	// future msg
	var texts []string
	for _, s := range content.Schedules.Schedules[:len(content.Schedules.Schedules)-2] {
		startTime := util.Time.LocalTime(s.StartTime, timezone.Minute()).Format(timeTemplate)
		endTime := util.Time.LocalTime(s.EndTime, timezone.Minute()).Format(timeTemplate)
		texts = append(texts, printer.Sprintf(textKeySalmonSchedulesSchedule, startTime, endTime))
	}
	tag := printer.Sprintf(textKeySalmonSchedulesFutureTag)
	text := strings.Join(texts, "") + tag
	futureMsg := botMessage.NewByUpdate(update, text, nil)
	// further detail msg
	s := content.Schedules.Details[salmon.SchedulesIdx.Further]
	startTime := util.Time.LocalTime(s.StartTime, timezone.Minute()).Format(timeTemplate)
	endTime := util.Time.LocalTime(s.EndTime, timezone.Minute()).Format(timeTemplate)
	text = printer.Sprintf(textKeySalmonSchedulesDetail,
		startTime, endTime, s.Stage.Name,
		s.Weapons[0].Weapon.Name,
		s.Weapons[1].Weapon.Name,
		s.Weapons[2].Weapon.Name,
		s.Weapons[3].Weapon.Name,
	)
	tag = printer.Sprintf(textKeySalmonSchedulesNextTag)
	text = text + tag
	furtherMsg := botApi.NewPhotoShare(update.Message.Chat.ID, string(content.ImageIDs[salmon.SchedulesIdx.Further]))
	furtherMsg.Caption = text
	furtherMsg.ParseMode = "Markdown"
	// latest detail msg
	s = content.Schedules.Details[salmon.SchedulesIdx.Latest]
	startTime = util.Time.LocalTime(s.StartTime, timezone.Minute()).Format(timeTemplate)
	endTime = util.Time.LocalTime(s.EndTime, timezone.Minute()).Format(timeTemplate)
	text = printer.Sprintf(textKeySalmonSchedulesDetail,
		startTime, endTime, s.Stage.Name,
		s.Weapons[0].Weapon.Name,
		s.Weapons[1].Weapon.Name,
		s.Weapons[2].Weapon.Name,
		s.Weapons[3].Weapon.Name,
	)
	now := time.Now().Unix()
	if now > s.StartTime {
		remainingTime := time.Until(time.Unix(s.EndTime, 0))
		h, m := getHourAndMinute(remainingTime)
		tag = printer.Sprintf(textKeySalmonSchedulesOpenTag, h, m)
	} else {
		remainingTime := time.Until(time.Unix(s.StartTime, 0))
		h, m := getHourAndMinute(remainingTime)
		tag = printer.Sprintf(textKeySalmonSchedulesSoonTag, h, m)
	}
	text = text + tag
	latestMsg := botApi.NewPhotoShare(update.Message.Chat.ID, string(content.ImageIDs[salmon.SchedulesIdx.Latest]))
	latestMsg.Caption = text
	latestMsg.ParseMode = "Markdown"

	return []botApi.Chattable{futureMsg, furtherMsg, latestMsg}
}

func getHourAndMinute(ts time.Duration) (int64, int64) {
	ts = ts.Round(time.Minute)
	hour := ts / time.Hour
	ts -= hour * time.Hour
	minute := ts / time.Minute
	return int64(hour), int64(minute)
}
