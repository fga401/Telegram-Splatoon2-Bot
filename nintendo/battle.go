package nintendo

import (
	"compress/gzip"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

func GetResults(iksm string, timezone int, acceptLang string) (interface{}, error) {
	reqUrl := "https://app.splatoon2.nintendo.net/api/results"
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return nil, errors.Wrap(err, "can't generate request")
	}
	req.Header = getAppHeader(iksm, timezone, acceptLang, true)
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "can't get response")
	}
	defer closeBody(resp.Body)
	respBody, err := gzip.NewReader(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "can't unzip response body")
	}
	respJson, err := ioutil.ReadAll(respBody)
	if err != nil {
		return "", errors.Wrap(err, "can't read response body")
	}
	//sessionToken := json.Get(respJson, "session_token").ToString()
	//log.Debug("get session token", zap.String("session token", sessionToken), zap.ByteString("json", respJson))
	return respJson, nil
}
