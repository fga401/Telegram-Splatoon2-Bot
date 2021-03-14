package battle

import (
	"math"
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
	keyDefeat = "defeat"

	textKeyVictoryEmoji     = "✅"
	textKeyDefeatEmoji      = "❌"
	textKeyTimeTemplate     = "2006-01-02 15:04:05"
	textKeyBoldBattleResult = `*[/%s] [%s]*
- Time: %s
- Mode: %s - %s
- Stage: %s 
- Count: %s %s %s
- Weapon: %s
- K(A)/D/SP: *%d(%d)/%d/%d*
`
	textKeyBattleResult = `\[/%s] \[%s] 
- Time: %s
- Mode: %s - %s
- Stage: %s 
- Count: %s %s %s
- Weapon: %s
- K(A)/D/SP: *%d(%d)/%d/%d*
`
)

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
	var myTeamCount, otherTeamCount float64
	var myTeamCountString, otherTeamCountString string
	switch battle.Type() {
	case nintendo.BattleResultTypeEnum.Regular:
		battle := battle.(*nintendo.RegularBattleResult)
		myTeamCount = float64(battle.MyTeamPercentage)
		otherTeamCount = float64(battle.OtherTeamPercentage)
		myTeamCountString = strconv.FormatFloat(myTeamCount, 'f', 1, 64)
		otherTeamCountString = strconv.FormatFloat(otherTeamCount, 'f', 1, 64)
	case nintendo.BattleResultTypeEnum.Gachi:
		battle := battle.(*nintendo.GachiBattleResult)
		myTeamCount = float64(battle.MyTeamCount)
		otherTeamCount = float64(battle.OtherTeamCount)
		myTeamCountString = strconv.Itoa(int(myTeamCount))
		otherTeamCountString = strconv.Itoa(int(otherTeamCount))
	case nintendo.BattleResultTypeEnum.League:
		battle := battle.(*nintendo.LeagueBattleResult)
		myTeamCount = float64(battle.MyTeamCount)
		otherTeamCount = float64(battle.OtherTeamCount)
		myTeamCountString = strconv.Itoa(int(myTeamCount))
		otherTeamCountString = strconv.Itoa(int(otherTeamCount))
	case nintendo.BattleResultTypeEnum.Festival:
		battle := battle.(*nintendo.FesBattleResult)
		myTeamCount = float64(battle.MyTeamPercentage)
		otherTeamCount = float64(battle.OtherTeamPercentage)
		myTeamCountString = strconv.FormatFloat(myTeamCount, 'f', 1, 64)
		otherTeamCountString = strconv.FormatFloat(otherTeamCount, 'f', 1, 64)
	}
	textKey := textKeyBattleResult
	if emphasis {
		textKey = textKeyBoldBattleResult
	}
	emoji := textKeyVictoryEmoji
	if battle.Metadata().MyTeamResult.Key == keyDefeat {
		emoji = textKeyDefeatEmoji
	}
	ret := printer.Sprintf(textKey,
		printer.Sprintf(battle.Metadata().BattleNumber),
		printer.Sprintf(battle.Metadata().MyTeamResult.Name)+" "+emoji,
		util.Time.LocalTime(battle.Metadata().StartTime, timezone.Minute()).Format(template),
		printer.Sprintf(battle.Metadata().GameMode.Name), printer.Sprintf(battle.Metadata().Rule.Name),
		printer.Sprintf(battle.Metadata().Stage.Name),
		myTeamCountString, formatCount(myTeamCount, otherTeamCount), otherTeamCountString,
		printer.Sprintf(battle.Metadata().PlayerResult.Player.Weapon.Name),
		battle.Metadata().PlayerResult.KillCount+battle.Metadata().PlayerResult.AssistCount, battle.Metadata().PlayerResult.AssistCount, battle.Metadata().PlayerResult.DeathCount, battle.Metadata().PlayerResult.SpecialCount,
	)
	return ret
}

func formatCount(myCount, otherCount float64) string {
	myPct := myCount / (myCount + otherCount)
	otherPct := otherCount / (myCount + otherCount)
	mySeg := int(round(myPct, 0.1) * 10)
	otherSeg := int(round(otherPct, 0.1) * 10)
	return strings.Repeat("=", mySeg) + ">/<" + strings.Repeat("≈", otherSeg)
}

func round(x, unit float64) float64 {
	return math.Round(x/unit) * unit
}
