package repository

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	botApi "github.com/go-telegram-bot-api/telegram-bot-api"
	"golang.org/x/text/message"
	"telegram-splatoon2-bot/common/util"
	"telegram-splatoon2-bot/service/repository/stage"
	"telegram-splatoon2-bot/service/timezone"
	userSvc "telegram-splatoon2-bot/service/user"
	"telegram-splatoon2-bot/telegram/controller/internal/adapter"
	botMessage "telegram-splatoon2-bot/telegram/controller/internal/message"
)

var (
	primaryFilterRegExp           = regexp.MustCompile(`^(?P<primary>[lgrLGR]+)$`)
	ruleSecondaryFilterRegExp     = regexp.MustCompile(`^(?P<primary>[czrtCZRT]+)$`)
	nextNSecondFilterRegExp       = regexp.MustCompile(`^(?P<n>\d+)$`)
	betweenHourSecondFilterRegExp = regexp.MustCompile(`^[bB](?P<begin>\d+)[-_](?P<end>\d+)$`)
)

func (ctrl *repositoryCtrl) stage(update botApi.Update, argManager adapter.Manager, args ...interface{}) error {
	statusArgIdx := argManager.Index(ctrl.statusAdapter)[0]
	status := args[statusArgIdx].(userSvc.Status)

	filterArgs := update.Message.CommandArguments()
	primaryArg, SecondaryArgs := splitFilterArgs(filterArgs)
	primaryFilter, err := parsePrimaryFilterArgs(primaryArg)
	if err != nil {
		msg := getStageSchedulesWrongArgsMessage(ctrl.languageSvc.Printer(status.Language), update)
		_, err := ctrl.bot.Send(msg)
		return err
	}
	secondaryFilters := make([]stage.SecondaryFilter, 0)
	for _, arg := range SecondaryArgs {
		f, err := parseSecondFilterArgs(arg, status.Timezone)
		if err != nil {
			msg := getStageSchedulesWrongArgsMessage(ctrl.languageSvc.Printer(status.Language), update)
			_, err := ctrl.bot.Send(msg)
			return err
		}
		secondaryFilters = append(secondaryFilters, f)
	}

	content := ctrl.stageRepo.Content(primaryFilter, secondaryFilters, limit)
	if content == nil {
		msg := getStageSchedulesNoReadyMessage(ctrl.languageSvc.Printer(status.Language), update)
		_, err := ctrl.bot.Send(msg)
		return err
	}
	if len(content) > ctrl.limit {
		content = content[:ctrl.limit]
		msg := getStageSchedulesOverLimitMessage(ctrl.languageSvc.Printer(status.Language), update)
		_, err := ctrl.bot.Send(msg)
		return err
	}
	msgs := getStageSchedulesMessages(ctrl.languageSvc.Printer(status.Language), update, content, status.Timezone)
	for _, msg := range msgs {
		_, err := ctrl.bot.Send(msg)
		if err != nil {
			return err
		}
	}
	return nil
}

// splitFilterArgs splits text into primary filter args and secondary filter args.
// If primary filter args or secondary filter args are empty, default value will be filled.
func splitFilterArgs(text string) (string, []string) {
	var args = []string{"lgr", "1"} // default filter
	if text != "" {
		args = strings.Split(text, " ")
	}
	// primary filter not fount
	if !isPrimaryFilterArgs(args[0]) {
		primaryFilterArg := "lgr"
		idx := firstIndexOfSecondaryFilterParam(args[0])
		if idx > 0 {
			primaryFilterArg = args[0][:idx]
			args[0] = args[0][idx:]
		}
		args = append([]string{primaryFilterArg}, args...) // add primary filter
	}
	if len(args) == 1 {
		args = append(args, "2")
	}
	return args[0], args[1:]
}

func isPrimaryFilterArgs(text string) bool {
	return len(primaryFilterRegExp.FindStringSubmatch(text)) != 0
}

func firstIndexOfSecondaryFilterParam(text string) int {
	for i, c := range text {
		if c != 'l' && c != 'r' && c != 'g' {
			return i
		}
	}
	return len(text)
}

func parsePrimaryFilterArgs(text string) (stage.PrimaryFilter, error) {
	text = strings.ToLower(text)
	var modes []stage.Mode
	for _, c := range text {
		switch c {
		case 'l':
			modes = append(modes, stage.ModeEnum.League)
		case 'g':
			modes = append(modes, stage.ModeEnum.Gachi)
		case 'r':
			modes = append(modes, stage.ModeEnum.Regular)
		default:
			return stage.PrimaryFilter{}, errors.New("wrong primary filter args")
		}
	}
	return stage.NewPrimaryFilter(modes), nil
}

func parseSecondFilterArgs(text string, timezone timezone.Timezone) (stage.SecondaryFilter, error) {
	text = strings.ToLower(text)
	if args := ruleSecondaryFilterRegExp.FindStringSubmatch(text); len(args) != 0 {
		var zone, tower, clam, rainmaker bool
		for _, c := range args[1] {
			switch c {
			case 'z':
				zone = true
			case 't':
				tower = true
			case 'c':
				clam = true
			case 'r':
				rainmaker = true
			default:
				return nil, errors.New("unknown secondary filter args")
			}
		}
		return stage.NewRuleSecondaryFilter(zone, tower, clam, rainmaker), nil
	}
	if args := nextNSecondFilterRegExp.FindStringSubmatch(text); len(args) != 0 {
		n, err := strconv.Atoi(args[1])
		if err != nil {
			return nil, errors.New("unknown secondary filter args")
		}
		return stage.NewNextNSecondaryFilter(n), nil
	}
	if args := betweenHourSecondFilterRegExp.FindStringSubmatch(text); len(args) != 0 {
		begin, err := strconv.Atoi(args[1])
		if err != nil {
			return nil, errors.New("unknown secondary filter args")
		}
		end, err := strconv.Atoi(args[2])
		if err != nil {
			return nil, errors.New("unknown secondary filter args")
		}
		return stage.NewBetweenHourSecondaryFilter(begin, end, timezone), nil
	}
	return nil, errors.New("unknown secondary filter args")
}

const (
	textKeyStageSchedulesNoReady   = "Stage schedules have not been ready yet."
	textKeyStageSchedulesWrongArgs = `Wrong arguments. Please use /help\_stages to get help.`
	textKeyStageSchedulesOverLimit = "_Note: your query returns too many results, and some of them have been omitted to avoid reaching telegram rate limit._"
	textKeyStageSchedulesDetail    = "*Time*:\n`%s ~ %s`\n*Mode*: %s\n*Rule*: %s\n*Stage*:\n- %s\n- %s\n#%s  #%s"
)

func getStageSchedulesNoReadyMessage(printer *message.Printer, update botApi.Update) botApi.Chattable {
	text := printer.Sprintf(textKeyStageSchedulesNoReady)
	return botMessage.NewByUpdate(update, text, nil)
}

func getStageSchedulesWrongArgsMessage(printer *message.Printer, update botApi.Update) botApi.Chattable {
	text := printer.Sprintf(textKeyStageSchedulesWrongArgs)
	return botMessage.NewByUpdate(update, text, nil)
}

func getStageSchedulesOverLimitMessage(printer *message.Printer, update botApi.Update) botApi.Chattable {
	text := printer.Sprintf(textKeyStageSchedulesOverLimit)
	return botMessage.NewByUpdate(update, text, nil)
}

func getStageSchedulesMessages(printer *message.Printer, update botApi.Update, content []stage.WrappedSchedule, timezone timezone.Timezone) []botApi.Chattable {
	timeTemplate := printer.Sprintf(timeTemplateTextKey)
	var ret []botApi.Chattable
	for _, s := range content {
		msg := botApi.NewPhotoShare(update.Message.Chat.ID, string(s.ImageID))
		startTime := util.Time.LocalTime(s.Schedule.StartTime, timezone.Minute()).Format(timeTemplate)
		endTime := util.Time.LocalTime(s.Schedule.EndTime, timezone.Minute()).Format(timeTemplate)
		text := printer.Sprintf(textKeyStageSchedulesDetail,
			startTime, endTime,
			s.Schedule.GameMode.Name, s.Schedule.Rule.Name,
			s.Schedule.StageB.Name, s.Schedule.StageA.Name,
			strings.Replace(s.Schedule.GameMode.Name, " ", `\_`, -1),
			strings.Replace(s.Schedule.Rule.Name, " ", `\_`, -1),
		)
		msg.Caption = text
		msg.ParseMode = "Markdown"
		ret = append(ret, msg)
	}
	return ret
}
