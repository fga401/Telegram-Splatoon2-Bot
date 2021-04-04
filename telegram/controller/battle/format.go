package battle

import (
	"bytes"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	botApi "github.com/go-telegram-bot-api/telegram-bot-api"
	"golang.org/x/text/message"
	"telegram-splatoon2-bot/common/util"
	"telegram-splatoon2-bot/service/nintendo"
	"telegram-splatoon2-bot/service/timezone"
	botMessage "telegram-splatoon2-bot/telegram/controller/internal/message"
)

const (
	escapeChar = '`'

	textKeyVictoryEmoji  = `✅`
	textKeyDefeatEmoji   = `❌`
	textKeyRegularEmoji  = `🟢`
	textKeyGachiEmoji    = `🟠`
	textKeyLeagueEmoji   = `🔴`
	textKeyPrivateEmoji  = `🟣`
	textKeyFestivalEmoji = `🟡`

	textKeyTimeTemplate     = "2006-01-02 15:04:05"
	textKeyBoldBattleResult = `*[ /%s ] [ %s ]*
- Time: %s
- Mode: %s - %s
- Stage: %s
- Count: %s
- Weapon: %s
- K(A)/D/SP: *%d(%d)/%d/%d*
`
	textKeyBattleResult = `\[ /%s ] \[ %s ] 
- Time: %s
- Mode: %s - %s
- Stage: %s
- Count: %s
- Weapon: %s
- K(A)/D/SP: *%d(%d)/%d/%d*
`
	textKeyBattleDetailResult = `*[ /%s Detail ] [ %s ]*
- Start Time: %s
- End Time: %s
- Mode: %s - %s
- Stage: %s 
- Count: %s
*[ My Team ]*:
%s
*[ Other Team ]*:
%s
`
	textKeyPlayerResult = "    *[ %s ]*    `%s`\n        - Weapon: %s\n        - K(A)/D/SP: *%d(%d)/%d/%d*\n        - Point: %dp"
)

func formatDetailedBattleResults(printer *message.Printer, battle nintendo.DetailedBattleResult, timezone timezone.Timezone) string {
	template := printer.Sprintf(textKeyTimeTemplate)
	myPlayerResult := []nintendo.PlayerResult{battle.Metadata().PlayerResult}
	myTeamPlayerResults := battle.MyTeamPlayerResults()
	sort.Slice(myTeamPlayerResults, func(i, j int) bool {
		return myTeamPlayerResults[i].KillCount+myTeamPlayerResults[i].AssistCount > myTeamPlayerResults[j].KillCount+myTeamPlayerResults[j].AssistCount
	})
	myTeamPlayerResults = append(myPlayerResult, myTeamPlayerResults...)
	otherTeamPlayerResults := battle.OtherTeamPlayerResults()
	sort.Slice(otherTeamPlayerResults, func(i, j int) bool {
		return otherTeamPlayerResults[i].KillCount+otherTeamPlayerResults[i].AssistCount > otherTeamPlayerResults[j].KillCount+otherTeamPlayerResults[j].AssistCount
	})
	ret := printer.Sprintf(textKeyBattleDetailResult,
		printer.Sprintf(encodeBattleNumberCommand(battle.Metadata().BattleNumber)), formatTeamResult(printer, battle),
		util.Time.LocalTime(battle.Metadata().StartTime, timezone.Minute()).Format(template),
		util.Time.LocalTime(battle.EndTime(), timezone.Minute()).Format(template),
		formatMode(printer, battle), printer.Sprintf(battle.Metadata().Rule.Name),
		printer.Sprintf(battle.Metadata().Stage.Name),
		formatTeamCountBanner(battle),
		formatPlayerResults(printer, myTeamPlayerResults),
		formatPlayerResults(printer, otherTeamPlayerResults),
	)
	return ret
}

func formatPlayerResults(printer *message.Printer, results []nintendo.PlayerResult) string {
	texts := make([]string, 0, 4)
	for _, r := range results {
		text := printer.Sprintf(textKeyPlayerResult,
			r.Player.Udemae.Name, escapeNickName(r.Player.Nickname), printer.Sprintf(r.Player.Weapon.Name),
			r.KillCount+r.AssistCount, r.AssistCount, r.DeathCount, r.SpecialCount,
			r.GamePaintPoint,
		)
		texts = append(texts, text)
	}
	return strings.Join(texts, "\n")
}

func (ctrl *battleCtrl) formatBattleResults(printer *message.Printer, update botApi.Update, battles []nintendo.BattleResult, timezone timezone.Timezone, emphasis []bool) []botApi.Chattable {
	texts := make([]string, 0, ctrl.maxResultsPerMessage)
	ret := make([]botApi.Chattable, 0)
	for i := len(battles) - 1; i >= 0; i-- {
		battle := battles[i]
		texts = append(texts, formatBattleResult(printer, battle, timezone, emphasis[i]))
		if len(texts) == ctrl.maxResultsPerMessage {
			text := strings.Join(texts, "\n")
			ret = append(ret, botMessage.NewByUpdate(update, text, nil))
			texts = texts[:0]
		}
	}
	if len(texts) > 0 {
		text := strings.Join(texts, "\n")
		ret = append(ret, botMessage.NewByUpdate(update, text, nil))
	}
	return ret
}

func formatBattleResult(printer *message.Printer, battle nintendo.BattleResult, timezone timezone.Timezone, emphasis bool) string {
	template := printer.Sprintf(textKeyTimeTemplate)
	textKey := textKeyBattleResult
	if emphasis {
		textKey = textKeyBoldBattleResult
	}
	ret := printer.Sprintf(textKey,
		printer.Sprintf(encodeBattleNumberCommand(battle.Metadata().BattleNumber)), formatTeamResult(printer, battle),
		util.Time.LocalTime(battle.Metadata().StartTime, timezone.Minute()).Format(template),
		formatMode(printer, battle), printer.Sprintf(battle.Metadata().Rule.Name),
		printer.Sprintf(battle.Metadata().Stage.Name),
		formatTeamCountBanner(battle),
		printer.Sprintf(battle.Metadata().PlayerResult.Player.Weapon.Name),
		battle.Metadata().PlayerResult.KillCount+battle.Metadata().PlayerResult.AssistCount, battle.Metadata().PlayerResult.AssistCount, battle.Metadata().PlayerResult.DeathCount, battle.Metadata().PlayerResult.SpecialCount,
	)
	return ret
}

func formatMode(printer *message.Printer, battle nintendo.BattleResult) string {
	return printer.Sprintf(battle.Metadata().GameMode.Name) + " " + modeEmoji(battle.Metadata().GameMode.Key)
}

func formatTeamResult(printer *message.Printer, battle nintendo.BattleResult) string {
	emoji := textKeyVictoryEmoji
	if battle.Metadata().MyTeamResult.Key == nintendo.KeyDefeat {
		emoji = textKeyDefeatEmoji
	}
	return printer.Sprintf(battle.Metadata().MyTeamResult.Name) + " " + emoji
}

func modeEmoji(key string) string {
	switch key {
	case nintendo.KeyRegular:
		return textKeyRegularEmoji
	case nintendo.KeyGachi:
		return textKeyGachiEmoji
	case nintendo.KeyLeaguePair:
		return textKeyLeagueEmoji
	case nintendo.KeyLeagueTeam:
		return textKeyLeagueEmoji
	case nintendo.KeyPrivate:
		return textKeyPrivateEmoji
	case nintendo.KeyFestivalSolo:
		return textKeyFestivalEmoji
	case nintendo.KeyFestivalTeam:
		return textKeyFestivalEmoji
	default:
		return ""
	}
}

func formatTeamCountBanner(battleRaw nintendo.BattleResult) string {
	var myTeamCount, otherTeamCount float64
	var myTeamCountString, otherTeamCountString string
	switch battleRaw.Type() {
	case nintendo.BattleResultTypeEnum.Regular:
		battle, ok := battleRaw.(*nintendo.RegularBattleResult)
		if !ok {
			battle = &battleRaw.(*nintendo.DetailedRegularBattleResult).RegularBattleResult
		}
		myTeamCount = float64(battle.MyTeamPercentage)
		otherTeamCount = float64(battle.OtherTeamPercentage)
		myTeamCountString = strconv.FormatFloat(myTeamCount, 'f', 1, 64)
		otherTeamCountString = strconv.FormatFloat(otherTeamCount, 'f', 1, 64)
	case nintendo.BattleResultTypeEnum.Gachi:
		battle, ok := battleRaw.(*nintendo.GachiBattleResult)
		if !ok {
			battle = &battleRaw.(*nintendo.DetailedGachiBattleResult).GachiBattleResult
		}
		myTeamCount = float64(battle.MyTeamCount)
		otherTeamCount = float64(battle.OtherTeamCount)
		myTeamCountString = strconv.Itoa(int(myTeamCount))
		otherTeamCountString = strconv.Itoa(int(otherTeamCount))
	case nintendo.BattleResultTypeEnum.League:
		battle, ok := battleRaw.(*nintendo.LeagueBattleResult)
		if !ok {
			battle = &battleRaw.(*nintendo.DetailedLeagueBattleResult).LeagueBattleResult
		}
		myTeamCount = float64(battle.MyTeamCount)
		otherTeamCount = float64(battle.OtherTeamCount)
		myTeamCountString = strconv.Itoa(int(myTeamCount))
		otherTeamCountString = strconv.Itoa(int(otherTeamCount))
	case nintendo.BattleResultTypeEnum.Festival:
		battle, ok := battleRaw.(*nintendo.FesBattleResult)
		if !ok {
			battle = &battleRaw.(*nintendo.DetailedFesBattleResult).FesBattleResult
		}
		myTeamCount = float64(battle.MyTeamPercentage)
		otherTeamCount = float64(battle.OtherTeamPercentage)
		myTeamCountString = strconv.FormatFloat(myTeamCount, 'f', 1, 64)
		otherTeamCountString = strconv.FormatFloat(otherTeamCount, 'f', 1, 64)
	}
	return fmt.Sprintf("%s %s %s", myTeamCountString, formatBanner(myTeamCount, otherTeamCount), otherTeamCountString)
}

func formatBanner(myCount, otherCount float64) string {
	myPct := myCount / (myCount + otherCount)
	otherPct := otherCount / (myCount + otherCount)
	mySeg := int(round(myPct, 0.1) * 10)
	otherSeg := int(round(otherPct, 0.1) * 10)
	return strings.Repeat("=", mySeg) + ">/<" + strings.Repeat("≈", otherSeg)
}

func escapeNickName(nickName string) string {
	buf := new(bytes.Buffer)
	i := strings.IndexByte(nickName, escapeChar)
	for ; i != -1; i = strings.IndexByte(nickName, escapeChar) {
		buf.WriteString(nickName[:i])
		buf.WriteString("`\\``")
		nickName = nickName[i+1:]
	}
	buf.WriteString(nickName)
	return buf.String()
}

func round(x, unit float64) float64 {
	return math.Round(x/unit) * unit
}
