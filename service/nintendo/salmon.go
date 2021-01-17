package nintendo

import (
	json "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	log "telegram-splatoon2-bot/common/log"
)

// GetSalmonSchedules returns SalmonSchedules and error.
// If error is caused by cookies expiration, it will return a ExpirationError
func GetSalmonSchedules(iksm string, timezone int, acceptLang string) (*SalmonSchedules, error) {
	reqUrl := "https://app.splatoon2.nintendo.net/api/coop_schedules"
	respJson, err := getSplatoon2RestfulJson(reqUrl, iksm, timezone, acceptLang)
	if isCookiesExpired(respJson) {
		return nil, &ExpirationError{iksm}
	}
	log.Debug("get salmon schedules", zap.ByteString("salmon_schedules", respJson))
	salmonSchedules := &SalmonSchedules{}
	err = json.Unmarshal(respJson, salmonSchedules)
	if err != nil || salmonSchedules.Details == nil || salmonSchedules.Schedules == nil {
		return nil, errors.Wrap(err, "can't parse json to SalmonSchedules")
	}
	return salmonSchedules,  nil
}
