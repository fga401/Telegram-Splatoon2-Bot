package nintendo

import (
	json "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	log "telegram-splatoon2-bot/common/log"
	"telegram-splatoon2-bot/service/language"
	"telegram-splatoon2-bot/service/timezone"
)

func unmarshalBattleResult(raw []byte) (ret BattleResult, err error) {
	t := json.Get(raw, "type").ToString()
	switch BattleResultType(t) {
	case BattleResultTypeEnum.Regular:
		{
			temp := &RegularBattleResult{}
			err = json.Unmarshal(raw, temp)
			if err != nil {
				return nil, err
			}
			ret = temp
		}
	case BattleResultTypeEnum.Gachi:
		{
			temp := &GachiBattleResult{}
			err = json.Unmarshal(raw, temp)
			if err != nil {
				return nil, err
			}
			ret = temp
		}
	case BattleResultTypeEnum.League:
		{
			temp := &LeagueBattleResult{}
			err = json.Unmarshal(raw, temp)
			if err != nil {
				return nil, err
			}
			ret = temp
		}
	case BattleResultTypeEnum.Festival:
		{
			temp := &FesBattleResult{}
			err = json.Unmarshal(raw, temp)
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
	Summary BattleSummary     `json:"summary"`
	Results []json.RawMessage `json:"results"`
}

// UnmarshalJSON implements Unmarshaler interface.
func (b *BattleResults) UnmarshalJSON(data []byte) error {
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
	b.Results = ret
	return nil
}

func (svc *impl) GetAllBattleResults(iksm string, timezone timezone.Timezone, language language.Language) (BattleResults, error) {
	reqURL := "https://app.splatoon2.nintendo.net/api/results"
	respJSON, err := svc.getSplatoon2RestfulJSON(reqURL, iksm, timezone.Minute(), language.IETF())
	if err != nil {
		return BattleResults{}, errors.Wrap(err, "can't get splatoon2 restful response")
	}
	if isCookiesExpired(respJSON) {
		return BattleResults{}, &ErrIKSMExpired{iksm}
	}
	log.Debug("get all battle results", zap.ByteString("all_battle_results", respJSON))
	battleResults := BattleResults{}
	err = json.Unmarshal(respJSON, &battleResults)
	if err != nil {
		return BattleResults{}, errors.Wrap(err, "can't parse json to BattleResults")
	}
	return battleResults, nil
}

type rawLatestBattleResult struct {
	RawResults []json.RawMessage `json:"results"`
}

func unmarshalRawLatestBattleResult(lastID string, min int, data []byte) ([]BattleResult, error) {
	raw := rawLatestBattleResult{}
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return nil, err
	}
	ret := make([]BattleResult, 0, len(raw.RawResults))
	found := false
	for _, r := range raw.RawResults {
		res, err := unmarshalBattleResult(r)
		if err != nil {
			return nil, err
		}
		if res.Metadata().BattleNumber == lastID {
			found = true
		}
		if found && len(ret) >= min {
			break
		}
		ret = append(ret, res)
	}
	return ret, nil
}

func (svc *impl) GetLatestBattleResults(lastID string, min int, iksm string, timezone timezone.Timezone, language language.Language) ([]BattleResult, error) {
	reqURL := "https://app.splatoon2.nintendo.net/api/results"
	respJSON, err := svc.getSplatoon2RestfulJSON(reqURL, iksm, timezone.Minute(), language.IETF())
	if err != nil {
		return nil, errors.Wrap(err, "can't get splatoon2 restful response")
	}
	if isCookiesExpired(respJSON) {
		return nil, &ErrIKSMExpired{iksm}
	}
	log.Debug("get last battle results", zap.ByteString("last_battle_results", respJSON))
	ret, err := unmarshalRawLatestBattleResult(lastID, min, respJSON)
	if err != nil {
		return nil, errors.Wrap(err, "can't parse json to slice of BattleResult")
	}
	return ret, nil
}

func unmarshalDetailedBattleResult(raw []byte) (ret DetailedBattleResult, err error) {
	t := json.Get(raw, "type").ToString()
	switch BattleResultType(t) {
	case BattleResultTypeEnum.Regular:
		{
			temp := &DetailedRegularBattleResult{}
			err = json.Unmarshal(raw, temp)
			if err != nil {
				return nil, err
			}
			ret = temp
		}
	case BattleResultTypeEnum.Gachi:
		{
			temp := &DetailedGachiBattleResult{}
			err = json.Unmarshal(raw, temp)
			if err != nil {
				return nil, err
			}
			ret = temp
		}
	case BattleResultTypeEnum.League:
		{
			temp := &DetailedLeagueBattleResult{}
			err = json.Unmarshal(raw, temp)
			if err != nil {
				return nil, err
			}
			ret = temp
		}
	case BattleResultTypeEnum.Festival:
		{
			temp := &DetailedFesBattleResult{}
			err = json.Unmarshal(raw, temp)
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

func (svc *impl) GetDetailedBattleResults(battleNumber string, iksm string, timezone timezone.Timezone, language language.Language) (DetailedBattleResult, error) {
	reqURL := "https://app.splatoon2.nintendo.net/api/results/" + battleNumber
	respJSON, err := svc.getSplatoon2RestfulJSON(reqURL, iksm, timezone.Minute(), language.IETF())
	if err != nil {
		return nil, errors.Wrap(err, "can't get splatoon2 restful response")
	}
	if isCookiesExpired(respJSON) {
		return nil, &ErrIKSMExpired{iksm}
	}
	log.Debug("get detailed battle results", zap.ByteString("detailed_battle_results", respJSON))
	ret, err := unmarshalDetailedBattleResult(respJSON)
	if err != nil {
		return nil, errors.Wrap(err, "can't parse json to DetailedBattleResult")
	}
	return ret, nil
}

type rawBattleSummary struct {
	Summary BattleSummary `json:"summary"`
}

func (svc *impl) GetBattleSummary(iksm string, timezone timezone.Timezone, language language.Language) (BattleSummary, error) {
	reqURL := "https://app.splatoon2.nintendo.net/api/results"
	respJSON, err := svc.getSplatoon2RestfulJSON(reqURL, iksm, timezone.Minute(), language.IETF())
	if err != nil {
		return BattleSummary{}, errors.Wrap(err, "can't get splatoon2 restful response")
	}
	if isCookiesExpired(respJSON) {
		return BattleSummary{}, &ErrIKSMExpired{iksm}
	}
	log.Debug("get battle summary", zap.ByteString("battle_summary", respJSON))
	ret := rawBattleSummary{}
	err = json.Unmarshal(respJSON, &ret)
	if err != nil {
		return BattleSummary{}, errors.Wrap(err, "can't parse json to BattleResults")
	}
	return ret.Summary, nil
}
