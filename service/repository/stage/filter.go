package stage

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"telegram-splatoon2-bot/common/util"
)

type GameModeName string

const (
	GameModeRegular GameModeName = "r"
	GameModeGachi   GameModeName = "g"
	GameModeLeague  GameModeName = "l"
)

type PrimaryFilter struct {
	orderByName []GameModeName
}

func (filter PrimaryFilter) Filter(schedules *content, secondaryFilters []SecondaryFilter, proposedN int) []CompositeSchedule {
	order := make([][]CompositeSchedule, 0)
	for _, name := range filter.orderByName {
		switch name {
		case GameModeLeague:
			order = append(order, schedules.leagueSchedules)
		case GameModeGachi:
			order = append(order, schedules.gachiSchedules)
		case GameModeRegular:
			order = append(order, schedules.regularSchedules)
		}
	}
	ret := make([]CompositeSchedule, 0)
	for i := 0; i < 12; i++ {
		for j := range order {
			s := order[j][i]
			keep := true
			for _, f := range secondaryFilters {
				if f.Filter(s) == false {
					keep = false
					break
				}
			}
			if keep {
				ret = append(ret, s)
			}
		}
		if len(ret) >= proposedN {
			break
		}
	}
	return ret
}

func NewPrimaryFilter(text string) (PrimaryFilter, error) {
	existed := make(map[GameModeName]struct{})
	ret := PrimaryFilter{}
	for _, c := range text {
		switch c {
		case 'l':
			if _, found := existed[GameModeLeague]; !found {
				ret.orderByName = append(ret.orderByName, GameModeLeague)
				existed[GameModeLeague] = struct{}{}
			}
		case 'r':
			if _, found := existed[GameModeRegular]; !found {
				ret.orderByName = append(ret.orderByName, GameModeRegular)
				existed[GameModeRegular] = struct{}{}
			}
		case 'g':
			if _, found := existed[GameModeGachi]; !found {
				ret.orderByName = append(ret.orderByName, GameModeGachi)
				existed[GameModeGachi] = struct{}{}
			}
		default:
			return ret, errors.Errorf("unknown character")
		}
	}
	return ret, nil
}

type SecondaryFilter interface {
	Filter(stage CompositeSchedule) bool
}

type RuleSecondaryFilter struct {
	allowRules map[string]struct{}
}

func (filter RuleSecondaryFilter) Filter(stage CompositeSchedule) bool {
	_, found := filter.allowRules[stage.Schedule.Rule.Key]
	return found
}
func NewRuleSecondaryFilter(text string) RuleSecondaryFilter {
	filter := RuleSecondaryFilter{allowRules: make(map[string]struct{})}
	for _, c := range text {
		switch c {
		case 'z':
			filter.allowRules["splat_zones"] = struct{}{}
		case 't':
			filter.allowRules["tower_control"] = struct{}{}
		case 'c':
			filter.allowRules["clam_blitz"] = struct{}{}
		case 'r':
			filter.allowRules["rainmaker"] = struct{}{}
		}
	}
	return filter
}

type TimeSecondaryFilter struct {
	begin, end int64
}

func (filter TimeSecondaryFilter) Filter(stage CompositeSchedule) bool {
	if filter.begin <= stage.Schedule.StartTime && stage.Schedule.StartTime < filter.end {
		return true
	}
	if filter.begin < stage.Schedule.EndTime && stage.Schedule.EndTime <= filter.end {
		return true
	}
	return false
}
func NewBetweenHourSecondaryFilter(begin string, end string, offset int) TimeSecondaryFilter {
	expectedBeginHour, err := strconv.Atoi(begin)
	if err != nil {
		expectedBeginHour = time.Now().Hour()
	}
	expectedEndHour, err := strconv.Atoi(end)
	if err != nil {
		expectedEndHour = (expectedBeginHour + 23) % 24
	}

	now := time.Now()
	userCurrentHour := util.Time.LocalTime(now.Unix(), offset).Hour()
	beginHourOffset := expectedBeginHour - userCurrentHour
	endHourOffset := expectedEndHour - userCurrentHour
	if beginHourOffset < 0 {
		beginHourOffset += 24
	}
	if endHourOffset < 0 {
		endHourOffset += 24
	}
	if endHourOffset-beginHourOffset < 0 {
		endHourOffset += 24
	}
	beginTime := now.Truncate(time.Hour).Add(time.Hour * time.Duration(beginHourOffset)).Unix()
	endTime := now.Truncate(time.Hour).Add(time.Hour * time.Duration(endHourOffset)).Unix()

	return TimeSecondaryFilter{begin: beginTime, end: endTime}
}
func NewNextNSecondaryFilter(text string) TimeSecondaryFilter {
	n, err := strconv.Atoi(text)
	if err != nil {
		n = 1
	}
	now := time.Now()
	beginTime := util.Time.SplatoonNextUpdateTime(now).Add(time.Hour * time.Duration(-2)).Unix()
	endTime := util.Time.SplatoonNextUpdateTime(now).Add(time.Hour * time.Duration(n*2-2)).Unix()
	return TimeSecondaryFilter{begin: beginTime, end: endTime}
}

type stageScheduleHelper struct {
	primaryFilterRegExp           *regexp.Regexp
	ruleSecondaryFilterRegExp     *regexp.Regexp
	nextNSecondFilterRegExp       *regexp.Regexp
	betweenHourSecondFilterRegExp *regexp.Regexp
}

var StageScheduleHelper = stageScheduleHelper{
	primaryFilterRegExp:           regexp.MustCompile(`^(?P<primary>[lgrLGR]+)$`),
	ruleSecondaryFilterRegExp:     regexp.MustCompile(`^(?P<primary>[czrtCZRT]+)$`),
	nextNSecondFilterRegExp:       regexp.MustCompile(`(?P<n>^\d+$)`),
	betweenHourSecondFilterRegExp: regexp.MustCompile(`^[bB](?P<begin>\d+)-(?P<end>\d+)$`),
}

func NewSecondaryFilter(text string, offset int) (SecondaryFilter, error) {
	text = strings.ToLower(text)
	if args := StageScheduleHelper.ruleSecondaryFilterRegExp.FindStringSubmatch(text); len(args) != 0 {
		return NewRuleSecondaryFilter(args[1]), nil
	}
	if args := StageScheduleHelper.nextNSecondFilterRegExp.FindStringSubmatch(text); len(args) != 0 {
		return NewNextNSecondaryFilter(args[1]), nil
	}
	if args := StageScheduleHelper.betweenHourSecondFilterRegExp.FindStringSubmatch(text); len(args) != 0 {
		return NewBetweenHourSecondaryFilter(args[1], args[2], offset), nil
	}
	return nil, errors.Errorf("unknown supported filter format")
}

func (stageScheduleHelper) firstIndexOfSecondaryFilterParam(text string) int {
	for i, c := range text {
		if c != 'l' && c != 'r' && c != 'g' {
			return i
		}
	}
	return len(text)
}
