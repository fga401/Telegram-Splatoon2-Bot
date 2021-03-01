package nintendo

import (
	json "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	log "telegram-splatoon2-bot/common/log"
	"telegram-splatoon2-bot/service/language"
	"telegram-splatoon2-bot/service/timezone"
)

// GetStageSchedules returns StageSchedules and error
// If error is caused by cookies expiration, it will return a ErrIKSMExpired
func (svc *impl) GetStageSchedules(iksm string, timezone timezone.Timezone, language language.Language) (StageSchedules, error) {
	reqUrl := "https://app.splatoon2.nintendo.net/api/schedules"
	respJson, err := svc.getSplatoon2RestfulJson(reqUrl, iksm, timezone.Minute(), language.IETF())
	if err != nil {
		return StageSchedules{}, errors.Wrap(err, "can't get splatoon2 restful response")
	}
	if isCookiesExpired(respJson) {
		return StageSchedules{}, &ErrIKSMExpired{iksm}
	}
	log.Debug("get stage schedules", zap.ByteString("stage_schedules", respJson))
	stageSchedules := StageSchedules{}
	err = json.Unmarshal(respJson, &stageSchedules)
	if err != nil || stageSchedules.Regular == nil || stageSchedules.Gachi == nil || stageSchedules.League == nil {
		return StageSchedules{}, errors.Wrap(err, "can't parse json to StageSchedules")
	}
	return stageSchedules, nil
}
