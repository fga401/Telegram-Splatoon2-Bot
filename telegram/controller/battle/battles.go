package battle

import (
	botApi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/text/message"
	"telegram-splatoon2-bot/common/log"
	"telegram-splatoon2-bot/service/language"
	"telegram-splatoon2-bot/service/nintendo"
	"telegram-splatoon2-bot/service/timezone"
	userSvc "telegram-splatoon2-bot/service/user"
	"telegram-splatoon2-bot/telegram/controller/internal/adapter"
	botMessage "telegram-splatoon2-bot/telegram/controller/internal/message"
)

func (ctrl *battleCtrl) battlePolling(update botApi.Update, argManager adapter.Manager, args ...interface{}) error {
	statusArgIdx := argManager.Index(ctrl.statusAdapter)[0]
	status := args[statusArgIdx].(userSvc.Status)
	start := false
	if _, ok := ctrl.pollingChats[status.UserID]; !ok {
		start = true
		ctrl.startPolling(status.UserID, update.Message.Chat.ID)
	} else {
		ctrl.stopPolling(status.UserID)
	}
	printer := ctrl.languageSvc.Printer(status.Language)
	msg := getBattlePollingMessage(printer, update, start)
	_, err := ctrl.bot.Send(msg)
	return err
}

const (
	textKeyBattlePollingStart = "Start polling battles."
	textKeyBattlePollingStop  = "Stop polling battles."
)

func getBattlePollingMessage(printer *message.Printer, update botApi.Update, start bool) botApi.Chattable {
	textKey := textKeyBattlePollingStart
	if !start {
		textKey = textKeyBattlePollingStop
	}
	text := printer.Sprintf(textKey)
	return botMessage.NewByUpdate(update, text, nil)
}

func (ctrl *battleCtrl) battleAll(update botApi.Update, argManager adapter.Manager, args ...interface{}) error {
	statusArgIdx := argManager.Index(ctrl.statusAdapter)[0]
	status := args[statusArgIdx].(userSvc.Status)
	battles, err := ctrl.nintendoSvc.GetAllBattleResults(status.IKSM, status.Timezone, language.English)
	printer := ctrl.languageSvc.Printer(status.Language)
	if errors.Is(err, &nintendo.ErrIKSMExpired{}) {
		msg := botMessage.UpdatingToken(printer, update)
		var resp *botApi.Message
		resp, err = ctrl.bot.Send(msg)
		if err != nil {
			log.Warn("can't send UpdateToken message")
		}
		status, err = ctrl.userSvc.UpdateStatusIKSM(status.UserID)
		if err != nil {
			msg := botMessage.InternalError(printer, resp)
			_, _ = ctrl.bot.Send(msg)
			return errors.Wrap(err, "can't update IKSM when fetching user's battles")
		}
		battles, err = ctrl.nintendoSvc.GetAllBattleResults(status.IKSM, status.Timezone, language.English)
		_, _ = ctrl.bot.Send(botApi.NewDeleteMessage(resp.Chat.ID, resp.MessageID))
	}
	if err != nil {
		return errors.Wrap(err, "can't fetches user's battles")
	}
	msgs := ctrl.getAllBattlesMessage(printer, update, battles, status.Timezone)
	for _, msg := range msgs {
		_, err = ctrl.bot.Send(msg)
	}
	lastBattleNumber := battles.Results[0].Metadata().BattleNumber
	_, err = ctrl.userSvc.UpdateStatusLastBattle(status.UserID, lastBattleNumber)
	if err != nil {
		log.Warn("can't update lastBattleNumber", zap.Int64("user_id", int64(status.UserID)), zap.Error(err))
	}
	return err
}

const (
	textKeyAllBattlesMessage = `- Use /battle\_summary to show summary.
- Use /battle\_last to show last battles.`
)

func (ctrl *battleCtrl) getAllBattlesMessage(printer *message.Printer, update botApi.Update, battles nintendo.BattleResults, timezone timezone.Timezone) []botApi.Chattable {
	emphasis := make([]bool, len(battles.Results))
	for i := range emphasis {
		emphasis[i] = true
	}
	msgs := ctrl.formatBattleResults(printer, update, battles.Results, timezone, emphasis)
	text := printer.Sprintf(textKeyAllBattlesMessage)
	msg := botMessage.NewByUpdate(update, text, nil)
	msgs = append(msgs, msg)
	return msgs
}

func (ctrl *battleCtrl) battleLast(update botApi.Update, argManager adapter.Manager, args ...interface{}) error {
	statusArgIdx := argManager.Index(ctrl.statusAdapter)[0]
	status := args[statusArgIdx].(userSvc.Status)
	battles, err := ctrl.nintendoSvc.GetLatestBattleResults(status.LastBattle, ctrl.minLastResults, status.IKSM, status.Timezone, language.English)
	printer := ctrl.languageSvc.Printer(status.Language)
	if errors.Is(err, &nintendo.ErrIKSMExpired{}) {
		msg := botMessage.UpdatingToken(printer, update)
		var resp *botApi.Message
		resp, err = ctrl.bot.Send(msg)
		if err != nil {
			log.Warn("can't send UpdateToken message")
		}
		status, err = ctrl.userSvc.UpdateStatusIKSM(status.UserID)
		if err != nil {
			msg := botMessage.InternalError(printer, resp)
			_, _ = ctrl.bot.Send(msg)
			return errors.Wrap(err, "can't update IKSM when fetching user's last battles")
		}
		battles, err = ctrl.nintendoSvc.GetLatestBattleResults(status.LastBattle, ctrl.minLastResults, status.IKSM, status.Timezone, language.English)
		_, _ = ctrl.bot.Send(botApi.NewDeleteMessage(resp.Chat.ID, resp.MessageID))
	}
	if err != nil {
		return errors.Wrap(err, "can't fetches user's last battles")
	}
	msgs := ctrl.getLastBattlesMessage(printer, update, status.LastBattle, battles, status.Timezone)
	for _, msg := range msgs {
		_, err = ctrl.bot.Send(msg)
	}
	lastBattleNumber := battles[0].Metadata().BattleNumber
	_, err = ctrl.userSvc.UpdateStatusLastBattle(status.UserID, lastBattleNumber)
	if err != nil {
		log.Warn("can't update lastBattleNumber", zap.Int64("user_id", int64(status.UserID)), zap.Error(err))
	}
	return err
}

const (
	textKeyLastBattlesMessage = `- Use /battle\_all to show last 50 battles.
- Use /battle\_last to show last battles.`
)

func (ctrl *battleCtrl) getLastBattlesMessage(printer *message.Printer, update botApi.Update, lastBattleNumber string, battles []nintendo.BattleResult, timezone timezone.Timezone) []botApi.Chattable {
	emphasis := make([]bool, len(battles))
	for i := range emphasis {
		if battles[i].Metadata().BattleNumber == lastBattleNumber {
			break
		}
		emphasis[i] = true
	}
	msgs := ctrl.formatBattleResults(printer, update, battles, timezone, emphasis)
	text := printer.Sprintf(textKeyLastBattlesMessage)
	msg := botMessage.NewByUpdate(update, text, nil)
	msgs = append(msgs, msg)
	return msgs
}

func (ctrl *battleCtrl) battleSummary(update botApi.Update, argManager adapter.Manager, args ...interface{}) error {
	statusArgIdx := argManager.Index(ctrl.statusAdapter)[0]
	status := args[statusArgIdx].(userSvc.Status)
	summary, err := ctrl.nintendoSvc.GetBattleSummary(status.IKSM, status.Timezone, language.English)
	printer := ctrl.languageSvc.Printer(status.Language)
	if errors.Is(err, &nintendo.ErrIKSMExpired{}) {
		msg := botMessage.UpdatingToken(printer, update)
		var resp *botApi.Message
		resp, err = ctrl.bot.Send(msg)
		if err != nil {
			log.Warn("can't send UpdateToken message")
		}
		status, err = ctrl.userSvc.UpdateStatusIKSM(status.UserID)
		if err != nil {
			msg := botMessage.InternalError(printer, resp)
			_, _ = ctrl.bot.Send(msg)
			return errors.Wrap(err, "can't update IKSM when fetching user's last battles")
		}
		summary, err = ctrl.nintendoSvc.GetBattleSummary(status.IKSM, status.Timezone, language.English)
		_, _ = ctrl.bot.Send(botApi.NewDeleteMessage(resp.Chat.ID, resp.MessageID))
	}
	if err != nil {
		return errors.Wrap(err, "can't fetches user's last battles")
	}
	msg := getBattleSummaryMessage(printer, update, summary)
	_, err = ctrl.bot.Send(msg)
	return err
}

const (
	textKeyBattleSummary = `*Battle Summary*
- Victory/Defeat: *%d / %d*
- Victory Rate: *%.2f*
- Average K(A)/D/SP: *%.1f (%.1f) / %.1f / %.1f*`
)

func getBattleSummaryMessage(printer *message.Printer, update botApi.Update, summary nintendo.BattleSummary) botApi.Chattable {
	text := printer.Sprintf(textKeyBattleSummary,
		summary.VictoryCount, summary.DefeatCount,
		summary.VictoryRate,
		summary.KillCountAverage+summary.AssistCountAverage, summary.AssistCountAverage, summary.DeathCountAverage, summary.SpecialCountAverage,
	)
	return botMessage.NewByUpdate(update, text, nil)
}
