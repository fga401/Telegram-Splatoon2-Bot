package nintendo

import (
	json "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	log "telegram-splatoon2-bot/common/log"
	"telegram-splatoon2-bot/service/language"
	"telegram-splatoon2-bot/service/timezone"
)

// GetSalmonSchedules returns SalmonSchedules and error.
// If error is caused by cookies expiration, it will return a ErrIKSMExpired
func (svc *impl) GetSalmonSchedules(iksm string, timezone timezone.Timezone, language language.Language) (SalmonSchedules, error) {
	reqUrl := "https://app.splatoon2.nintendo.net/api/coop_schedules"
	respJson, err := svc.getSplatoon2RestfulJson(reqUrl, iksm, timezone.Minute(), language.IETF())
	if isCookiesExpired(respJson) {
		return SalmonSchedules{}, &ErrIKSMExpired{iksm}
	}
	log.Debug("get salmon schedules", zap.ByteString("salmon_schedules", respJson))
	salmonSchedules := SalmonSchedules{}
	err = json.Unmarshal(respJson, &salmonSchedules)
	if err != nil || salmonSchedules.Details == nil || salmonSchedules.Schedules == nil {
		return SalmonSchedules{}, errors.Wrap(err, "can't parse json to SalmonSchedules")
	}
	return salmonSchedules,  nil
}
