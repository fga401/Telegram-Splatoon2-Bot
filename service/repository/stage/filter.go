package stage

import (
	"time"

	"telegram-splatoon2-bot/common/enum"
	"telegram-splatoon2-bot/common/util"
	"telegram-splatoon2-bot/service/timezone"
)

// Mode of stage.
type Mode enum.Enum

type modeEnum struct {
	Regular Mode
	Gachi   Mode
	League  Mode
}

// ModeEnum lists all Mode.
var ModeEnum = enum.Assign(&modeEnum{}).(*modeEnum)

const (
	ruleSplatZones   = "splat_zones"
	ruleTowerControl = "tower_control"
	ruleClamBlitz    = "clam_blitz"
	ruleRainmaker    = "rainmaker"
)

// Filter applies primaryFilter and secondaryFilters to content.
// The the exceeding part would be removed.
func (c *content) Filter(primaryFilter PrimaryFilter, secondaryFilters []SecondaryFilter, limit int) []WrappedSchedule {
	order := make([][]WrappedSchedule, 0)
	for _, name := range primaryFilter.orderByName {
		switch name {
		case ModeEnum.League:
			order = append(order, c.LeagueSchedules)
		case ModeEnum.Gachi:
			order = append(order, c.GachiSchedules)
		case ModeEnum.Regular:
			order = append(order, c.RegularSchedules)
		}
	}
	ret := make([]WrappedSchedule, 0)
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
		if len(ret) >= limit {
			break
		}
	}
	return ret
}

// PrimaryFilter determines the availability and order of Mode.
type PrimaryFilter struct {
	orderByName []Mode
}

// NewPrimaryFilter returns a PrimaryFilter.
func NewPrimaryFilter(modes []Mode) PrimaryFilter {
	existed := make(map[Mode]struct{})
	ret := PrimaryFilter{}
	for _, mode := range modes {
		switch mode {
		case ModeEnum.League:
			if _, found := existed[ModeEnum.League]; !found {
				ret.orderByName = append(ret.orderByName, ModeEnum.League)
				existed[ModeEnum.League] = struct{}{}
			}
		case ModeEnum.Regular:
			if _, found := existed[ModeEnum.Regular]; !found {
				ret.orderByName = append(ret.orderByName, ModeEnum.Regular)
				existed[ModeEnum.Regular] = struct{}{}
			}
		case ModeEnum.Gachi:
			if _, found := existed[ModeEnum.Gachi]; !found {
				ret.orderByName = append(ret.orderByName, ModeEnum.Gachi)
				existed[ModeEnum.Gachi] = struct{}{}
			}
		}
	}
	return ret
}

// SecondaryFilter filter schedules by other factor.
type SecondaryFilter interface {
	Filter(stage WrappedSchedule) bool
}

// RuleSecondaryFilter filter schedules by Rule.
type RuleSecondaryFilter struct {
	allowRules map[string]struct{}
}

// Filter applies RuleSecondaryFilter.
func (filter RuleSecondaryFilter) Filter(stage WrappedSchedule) bool {
	_, found := filter.allowRules[stage.Schedule.Rule.Key]
	return found
}

// NewRuleSecondaryFilter returns a RuleSecondaryFilter.
func NewRuleSecondaryFilter(zone, tower, clam, rainmaker bool) RuleSecondaryFilter {
	filter := RuleSecondaryFilter{allowRules: make(map[string]struct{})}
	if zone {
		filter.allowRules[ruleSplatZones] = struct{}{}
	}
	if tower {
		filter.allowRules[ruleTowerControl] = struct{}{}
	}
	if clam {
		filter.allowRules[ruleClamBlitz] = struct{}{}
	}
	if rainmaker {
		filter.allowRules[ruleRainmaker] = struct{}{}
	}
	return filter
}

// TimeSecondaryFilter keeps stages from begin (hour) to end (hour) in user timezone.
type TimeSecondaryFilter struct {
	begin, end int64
}

// Filter applies TimeSecondaryFilter.
func (filter TimeSecondaryFilter) Filter(stage WrappedSchedule) bool {
	if filter.begin <= stage.Schedule.StartTime && stage.Schedule.StartTime < filter.end {
		return true
	}
	if filter.begin < stage.Schedule.EndTime && stage.Schedule.EndTime <= filter.end {
		return true
	}
	return false
}

// NewBetweenHourSecondaryFilter returns a TimeSecondaryFilter.
func NewBetweenHourSecondaryFilter(beginHour int, endHour int, timezone timezone.Timezone) TimeSecondaryFilter {
	now := time.Now()
	userCurrentHour := util.Time.LocalTime(now.Unix(), timezone.Minute()).Hour()
	beginHourOffset := beginHour - userCurrentHour
	endHourOffset := endHour - userCurrentHour
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

// NewNextNSecondaryFilter returns a TimeSecondaryFilter that keeps the next n stages.
func NewNextNSecondaryFilter(n int) TimeSecondaryFilter {
	now := time.Now()
	beginTime := util.Time.SplatoonNextUpdateTime(now).Add(time.Hour * time.Duration(-2)).Unix()
	endTime := util.Time.SplatoonNextUpdateTime(now).Add(time.Hour * time.Duration(n*2-2)).Unix()
	return TimeSecondaryFilter{begin: beginTime, end: endTime}
}
