package nintendo

import (
	"strconv"

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
	reqURL := "https://app.splatoon2.nintendo.net/api/coop_schedules"
	respJSON, err := svc.getSplatoon2RestfulJSON(reqURL, iksm, timezone.Minute(), language.IETF())
	if isCookiesExpired(respJSON) {
		return SalmonSchedules{}, &ErrIKSMExpired{iksm}
	}
	log.Debug("get salmon schedules", zap.ByteString("salmon_schedules", respJSON))
	salmonSchedules := SalmonSchedules{}
	err = json.Unmarshal(respJSON, &salmonSchedules)
	if err != nil || salmonSchedules.Details == nil || salmonSchedules.Schedules == nil {
		return SalmonSchedules{}, errors.Wrap(err, "can't parse json to SalmonSchedules")
	}
	return salmonSchedules, nil
}

func (svc *impl) GetAllSalmonResults(iksm string, timezone timezone.Timezone, language language.Language) (SalmonSummary, error) {
	reqURL := "https://app.splatoon2.nintendo.net/api/coop_results"
	respJSON, err := svc.getSplatoon2RestfulJSON(reqURL, iksm, timezone.Minute(), language.IETF())
	if err != nil {
		return SalmonSummary{}, errors.Wrap(err, "can't get splatoon2 restful response")
	}
	if isCookiesExpired(respJSON) {
		return SalmonSummary{}, &ErrIKSMExpired{iksm}
	}
	log.Debug("get salmon summary", zap.ByteString("salmon_summary", respJSON))
	ret := SalmonSummary{}
	err = json.Unmarshal(respJSON, &ret)
	if err != nil {
		return SalmonSummary{}, errors.Wrap(err, "can't parse json to SalmonSummary")
	}
	return ret, nil
}

type rawLatestSalmonResult struct {
	RawResults []json.RawMessage `json:"results"`
}

func unmarshalRawLatestSalmonResult(lastID int32, data []byte) ([]SalmonResult, error) {
	raw := rawLatestSalmonResult{}
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return nil, err
	}
	ret := make([]SalmonResult, 0, len(raw.RawResults))
	for _, r := range raw.RawResults {
		res := SalmonResult{}
		err = json.Unmarshal(r, &res)
		if err != nil {
			return nil, err
		}
		if res.JobID == lastID {
			break
		}
		ret = append(ret, res)
	}
	return ret, nil
}

func (svc *impl) GetLatestSalmonResults(lastBattleNumber int32, iksm string, timezone timezone.Timezone, language language.Language) ([]SalmonResult, error) {
	reqURL := "https://app.splatoon2.nintendo.net/api/coop_results"
	respJSON, err := svc.getSplatoon2RestfulJSON(reqURL, iksm, timezone.Minute(), language.IETF())
	if err != nil {
		return nil, errors.Wrap(err, "can't get splatoon2 restful response")
	}
	if isCookiesExpired(respJSON) {
		return nil, &ErrIKSMExpired{iksm}
	}
	log.Debug("get latest salmon summary", zap.ByteString("latest_salmon_summary", respJSON))
	ret, err := unmarshalRawLatestSalmonResult(lastBattleNumber, respJSON)
	if err != nil {
		return nil, errors.Wrap(err, "can't parse json to slice of SalmonResult")
	}
	return ret, nil
}

func (svc *impl) GetDetailedSalmonResults(battleNumber int32, iksm string, timezone timezone.Timezone, language language.Language) (SalmonDetailedResult, error) {
	reqURL := "https://app.splatoon2.nintendo.net/api/coop_results/" + strconv.Itoa(int(battleNumber))
	respJSON, err := svc.getSplatoon2RestfulJSON(reqURL, iksm, timezone.Minute(), language.IETF())
	if err != nil {
		return SalmonDetailedResult{}, errors.Wrap(err, "can't get splatoon2 restful response")
	}
	if isCookiesExpired(respJSON) {
		return SalmonDetailedResult{}, &ErrIKSMExpired{iksm}
	}
	log.Debug("get detailed salmon results", zap.ByteString("detailed_salmon_results", respJSON))
	ret := SalmonDetailedResult{}
	err = json.Unmarshal(respJSON, &ret)
	if err != nil {
		return SalmonDetailedResult{}, errors.Wrap(err, "can't parse json to SalmonDetailedResult")
	}
	return ret, nil
}
