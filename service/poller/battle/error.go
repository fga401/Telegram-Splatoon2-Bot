package battle

import (
	"telegram-splatoon2-bot/common/enum"
)

// ErrCanceledPolling wraps the error causing polling canceled.
type ErrCanceledPolling struct {
	Reason CancelReason
}

func (err *ErrCanceledPolling) Error() string {
	return errStringMap[err.Reason]
}

// Is checks if an error is ErrIKSMExpired.
func (err *ErrCanceledPolling) Is(e error) bool {
	_, ok := e.(*ErrCanceledPolling)
	return ok
}

// CancelReason identifies the reason of canceling polling.
type CancelReason enum.Enum
type cancelReasonEnum struct {
	NoNewBattles CancelReason
}

var (
	// CancelReasonEnum lists all available CancelReason.
	CancelReasonEnum = enum.Assign(&cancelReasonEnum{}).(*cancelReasonEnum)

	errStringMap     = map[CancelReason]string{
		CancelReasonEnum.NoNewBattles: "no new battles for a long time",
	}
)
