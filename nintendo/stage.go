package nintendo

import (
	json "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	log "telegram-splatoon2-bot/logger"
)

// GetStageSchedules returns StageSchedules and error
// If error is caused by cookies expiration, it will return a ExpirationError
func GetStageSchedules(iksm string, timezone int, acceptLang string) (*StageSchedules, error) {
	reqUrl := "https://app.splatoon2.nintendo.net/api/schedules"
	respJson, err := getSplatoon2RestfulJson(reqUrl, iksm, timezone, acceptLang)
	if err != nil {
		return nil, errors.Wrap(err, "can't get splatoon2 restful response")
	}
	if isCookiesExpired(respJson) {
		return nil, &ExpirationError{iksm}
	}
	log.Debug("get stage schedules", zap.ByteString("stage_schedules", respJson))
	stageSchedules := &StageSchedules{}
	err = json.Unmarshal(respJson, stageSchedules)
	if err != nil || stageSchedules.Regular == nil || stageSchedules.Gachi == nil || stageSchedules.League == nil {
		return nil, errors.Wrap(err, "can't parse json to StageSchedules")
	}
	return stageSchedules, nil
}
