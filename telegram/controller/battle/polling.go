package battle

import (
	"strings"

	botApi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/text/message"
	"telegram-splatoon2-bot/common/log"
	"telegram-splatoon2-bot/service/nintendo"
	battlePoller "telegram-splatoon2-bot/service/poller/battle"
	"telegram-splatoon2-bot/service/timezone"
	botMessage "telegram-splatoon2-bot/telegram/controller/internal/message"
)

func (ctrl *battleCtrl) getChatID(id UserID) (int64, bool) {
	ctrl.pollingMutex.RLock()
	defer ctrl.pollingMutex.RUnlock()
	chatID, ok := ctrl.pollingChats[id]
	return chatID, ok
}

func (ctrl *battleCtrl)stopPolling(id UserID){
	ctrl.pollingMutex.Lock()
	defer ctrl.pollingMutex.Unlock()
	ctrl.battlePoller.Stop(id)
	delete(ctrl.pollingChats, id)
}

func (ctrl *battleCtrl)startPolling(userID UserID, chatID int64){
	ctrl.pollingMutex.Lock()
	defer ctrl.pollingMutex.Unlock()
	ctrl.battlePoller.Start(userID)
	ctrl.pollingChats[userID] = chatID
}

func (ctrl *battleCtrl) sendPolledResult(result battlePoller.Result) {
	if result.Error != nil {
		if errors.Is(result.Error, &battlePoller.ErrCanceledPolling{}) {
			log.Debug("polling stopped by internal error", zap.Int64("user_id", int64(result.UserID)), zap.Error(result.Error))
			if chatID, ok := ctrl.getChatID(result.UserID); ok {
				status, err := ctrl.userSvc.GetStatus(result.UserID)
				if err != nil {
					log.Error("can't fetch status when polling battles.", zap.Int64("user_id", int64(result.UserID)), zap.Error(err))
					return
				}
				ctrl.stopPolling(status.UserID)
				printer := ctrl.languageSvc.Printer(status.Language)
				msg := getBattlePollingCancellationMessage(printer, chatID, result.Error.(*battlePoller.ErrCanceledPolling))
				_, _ = ctrl.bot.Send(msg)
			}
		} else {
			log.Warn("invalid result", zap.Int64("user_id", int64(result.UserID)), zap.Error(result.Error))
		}
		return
	}
	if len(result.Battles) == 0 {
		return
	}
	if chatID, ok := ctrl.getChatID(result.UserID); ok {
		status, err := ctrl.userSvc.GetStatus(result.UserID)
		if err != nil {
			log.Error("can't fetch status when polling battles.", zap.Int64("user_id", int64(result.UserID)), zap.Error(err))
			return
		}
		_, err = ctrl.userSvc.UpdateStatusLastBattle(result.UserID, result.Battles[0].Metadata().BattleNumber)
		if err != nil {
			log.Warn("can't update last battle number when polling battles.", zap.Int64("user_id", int64(result.UserID)), zap.Error(err))
		}
		printer := ctrl.languageSvc.Printer(status.Language)
		messages := ctrl.formatBattleResultsByChatID(printer, chatID, result.Battles, status.Timezone)
		for _, msg := range messages {
			_, _ = ctrl.bot.Send(msg)
		}
	}
}

func (ctrl *battleCtrl) pollingRoutine() {
	if ctrl.pollingMaxWorker > 0 {
		for i := int32(0); i < ctrl.pollingMaxWorker; i++ {
			go func() {
				for result := range ctrl.battlePoller.Results() {
					ctrl.sendPolledResult(result)
				}
			}()
		}
	} else {
		for result := range ctrl.battlePoller.Results() {
			go ctrl.sendPolledResult(result)
		}
	}
}

func (ctrl *battleCtrl) formatBattleResultsByChatID(printer *message.Printer, chatID int64, battles []nintendo.BattleResult, timezone timezone.Timezone) []botApi.Chattable {
	texts := make([]string, 0, ctrl.maxResultsPerMessage)
	ret := make([]botApi.Chattable, 0)
	for i := len(battles) - 1; i >= 0; i-- {
		battle := battles[i]
		texts = append(texts, formatBattleResult(printer, battle, timezone, true))
		if len(texts) == ctrl.maxResultsPerMessage {
			text := strings.Join(texts, "\n")
			ret = append(ret, botMessage.NewByChatID(chatID, text, nil))
			texts = texts[:0]
		}
	}
	if len(texts) > 0 {
		text := strings.Join(texts, "\n")
		ret = append(ret, botMessage.NewByChatID(chatID, text, nil))
	}
	return ret
}

const (
	textKeyBattlePollingCancellation = "Polling has been stopped. %s"
	textKeyBattlePollingCancellationReasonNoNewBattles = "No new battles for a long time"
)

func getBattlePollingCancellationMessage(printer *message.Printer, chatID int64, cause *battlePoller.ErrCanceledPolling) botApi.Chattable {
	var reasonTextKey string
	switch cause.Reason {
	case battlePoller.CancelReasonEnum.NoNewBattles:
		reasonTextKey = textKeyBattlePollingCancellationReasonNoNewBattles
	}
	return botMessage.NewByChatID(chatID, printer.Sprintf(textKeyBattlePollingCancellation, printer.Sprintf(reasonTextKey)), nil)
}
