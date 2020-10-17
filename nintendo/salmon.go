package nintendo

import (
	"compress/gzip"
	json "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	log "telegram-splatoon2-bot/logger"
)

// GetSalmonSchedules returns SalmonSchedules, whether cookie is expired, and error
func GetSalmonSchedules(iksm string, timezone int, acceptLang string) (*SalmonSchedules, bool, error) {
	reqUrl := "https://app.splatoon2.nintendo.net/api/coop_schedules"
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return nil, false, errors.Wrap(err, "can't generate request")
	}
	req.Header = getAppHeader(iksm, timezone, acceptLang, true)
	resp, err := client.Do(req)
	if err != nil {
		return nil, false, errors.Wrap(err, "can't get response")
	}
	defer closeBody(resp.Body)
	respBody, err := gzip.NewReader(resp.Body)
	if err != nil {
		return nil, false, errors.Wrap(err, "can't unzip response body")
	}
	respJson, err := ioutil.ReadAll(respBody)
	if err != nil {
		return nil, false, errors.Wrap(err, "can't read response body")
	}
	log.Debug("get salmon schedules", zap.ByteString("salmon_schedules", respJson))
	salmonSchedules := &SalmonSchedules{}
	err = json.Unmarshal(respJson, salmonSchedules)
	if err == nil && salmonSchedules.Details != nil && salmonSchedules.Schedules != nil{
		return salmonSchedules, false, nil
	}
	expired := json.Get(respJson, "code").ToString() == "AUTHENTICATION_ERROR"
	if expired {
		return nil, true, nil
	}
	return nil, false, errors.Wrap(err, "can't parse json to SalmonSchedules")
}
