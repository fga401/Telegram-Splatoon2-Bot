package battle

import (
	"time"

	"telegram-splatoon2-bot/service/nintendo"
	"telegram-splatoon2-bot/service/poller"
	"telegram-splatoon2-bot/service/user"
)

// Result is the result fetched by poller.
type Result struct {
	UserID  user.ID
	Battles []nintendo.BattleResult
	Detail  nintendo.DetailedBattleResult
	Error   error
}

type task struct {
	UserID     user.ID
	UpdateTime time.Time
}

type statistics struct {
	LastBattle nintendo.BattleResult
	CreateTime time.Time
}

// Service wrapper poller.Poller with battle Result.
type Service interface {
	poller.Poller
	Results() <-chan Result
}
