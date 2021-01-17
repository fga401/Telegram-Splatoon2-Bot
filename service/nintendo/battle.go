package nintendo

import (
	json "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	log "telegram-splatoon2-bot/common/log"
)

func unmarshalBattleResult(raw []byte) (ret BattleResult,err error) {
	t := json.Get(raw, "type").ToString()
	switch BattleResultType(t) {
	case BattleResultTypeEnum.Regular:
		{
			temp := RegularBattleResult{}
			err = json.Unmarshal(raw, &temp)
			if err != nil {
				return nil, err
			}
			ret = temp
		}
	case BattleResultTypeEnum.Gachi:
		{
			temp := GachiBattleResult{}
			err = json.Unmarshal(raw, &temp)
			if err != nil {
				return nil, err
			}
			ret = temp
		}
	case BattleResultTypeEnum.League:
		{
			temp := LeagueBattleResult{}
			err = json.Unmarshal(raw, &temp)
			if err != nil {
				return nil, err
			}
			ret = temp
		}
	case BattleResultTypeEnum.Festival:
		{
			temp := FesBattleResult{}
			err = json.Unmarshal(raw, &temp)
			if err != nil {
				return nil, err
			}
			ret = temp
		}
	default:
		return nil, errors.Errorf("unknown type")
	}
	return
}

type rawBattleResults struct {
	ID      string            `json:"unique_id"`
	Summary *BattleSummary    `json:"summary"`
	Results []json.RawMessage `json:"results"`
}

func (b *rawBattleResults) UnmarshalJSON(data []byte) error {
	rawBattleResults := &rawBattleResults{}
	err := json.Unmarshal(data, rawBattleResults)
	if err != nil {
		return err
	}
	ret := make([]BattleResult, len(rawBattleResults.Results))
	for i, raw := range rawBattleResults.Results {
		ret[i], err = unmarshalBattleResult(raw)
		if err != nil {
			return err
		}
	}
	b.ID = rawBattleResults.ID
	b.Summary = rawBattleResults.Summary
	//b.Results = ret
	return nil
}

func GetAllBattleResults(iksm string, timezone int, acceptLang string) (*BattleResults, error) {
	reqUrl := "https://app.splatoon2.nintendo.net/api/results"
	respJson, err := getSplatoon2RestfulJson(reqUrl, iksm, timezone, acceptLang)
	if err != nil {
		return nil, errors.Wrap(err, "can't get splatoon2 restful response")
	}
	if isCookiesExpired(respJson) {
		return nil, &ExpirationError{iksm}
	}
	log.Debug("get stage schedules", zap.ByteString("stage_schedules", respJson))
	battleResults := &BattleResults{}
	err = json.Unmarshal(respJson, battleResults)
	if err != nil {
		return nil, errors.Wrap(err, "can't parse json to StageSchedules")
	}
	return battleResults, nil
}

type rawLatestBattleResult struct {
	lastBattleNumber string
	RawResults []json.RawMessage `json:"results"`
	Results []BattleResult
}

func (b *rawLatestBattleResult) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, b)
	if err != nil {
		return err
	}
	ret := make([]BattleResult, 0, len(b.RawResults))
	for _, raw := range b.RawResults {
		res, err := unmarshalBattleResult(raw)
		if err != nil {
			return err
		}
		if res.Metadata().BattleNumber == b.lastBattleNumber {
			break
		}
		ret = append(ret, res)
	}
	b.Results = ret
	return nil
}

func GetLatestBattleResults(lastBattleNumber string, iksm string, timezone int, acceptLang string) ([]BattleResult, error) {
	reqUrl := "https://app.splatoon2.nintendo.net/api/results"
	respJson, err := getSplatoon2RestfulJson(reqUrl, iksm, timezone, acceptLang)
	if err != nil {
		return nil, errors.Wrap(err, "can't get splatoon2 restful response")
	}
	if isCookiesExpired(respJson) {
		return nil, &ExpirationError{iksm}
	}
	log.Debug("get stage schedules", zap.ByteString("stage_schedules", respJson))
	latestBattleResults := &rawLatestBattleResult{lastBattleNumber: lastBattleNumber}
	err = json.Unmarshal(respJson, latestBattleResults)
	if err != nil {
		return nil, errors.Wrap(err, "can't parse json to StageSchedules")
	}
	return latestBattleResults.Results, nil
}
