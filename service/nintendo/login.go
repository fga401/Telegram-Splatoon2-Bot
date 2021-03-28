package nintendo

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	json "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"telegram-splatoon2-bot/common/log"
	"telegram-splatoon2-bot/common/util"
	"telegram-splatoon2-bot/service/language"
)

// IsRedirectLinkValid return true if the link doesn't start with the prefix.
func IsRedirectLinkValid(link string) bool {
	if !strings.HasPrefix(link, prefix) {
		return false
	}
	_, err := url.ParseQuery(link[prefixLen:])
	if err != nil {
		return false
	}
	return true
}

func (svc *impl) NewProofKey() ([]byte, error) {
	proofKey, err := randBytes(32)
	if err != nil {
		return nil, errors.Wrap(err, "can't generate proof key")
	}
	return base64UrlEncode(proofKey), nil
}

func (svc *impl) NewLoginLink(proofKey []byte) (string, error) {
	state, err := randBytes(36)
	if err != nil {
		return "", errors.Wrap(err, "can't generate state")
	}
	hashProofKey := sha256.Sum256(proofKey)
	challenge := string(base64UrlEncode(hashProofKey[:]))
	authState := string(base64UrlEncode(state))

	hardcodeURL := "https://accounts.nintendo.com/connect/1.0.0/authorize?" +
		"redirect_uri=npf71b963c1b7b6d119://auth" +
		"&client_id=71b963c1b7b6d119" +
		"&scope=openid%20user%20user.birthday%20user.mii%20user.screenName" +
		"&response_type=session_token_code" +
		"&session_token_code_challenge_method=S256" +
		"&theme=login_form" +
		"&state=" + authState +
		"&session_token_code_challenge=" + challenge

	return hardcodeURL, nil
}

func (svc *impl) GetSessionToken(link string, proofKey []byte, language language.Language) (string, error) {
	sessionTokenCode, err := svc.getSessionTokenCode(link)
	if err != nil {
		// todo: invalid operation count ++
		return "", errors.Wrap(err, "invalid redirect link")
	}

	var sessionToken string
	err = util.Retry(func() error {
		var err error
		sessionToken, err = svc.getSessionToken(proofKey, sessionTokenCode, language.IETF())
		return err
	}, svc.retryTimes)
	if err != nil {
		return "", errors.Wrap(err, "can't fetch sessionToken")
	}
	return sessionToken, nil
}

func (svc *impl) GetAccountMetadata(sessionToken string, language language.Language) (AccountMetadata, error) {
	var accessToken string
	err := util.Retry(func() error {
		var err error
		accessToken, err = svc.getAccessToken(sessionToken, language.IETF())
		return err
	}, svc.retryTimes)
	if err != nil {
		return AccountMetadata{}, errors.Wrap(err, "can't get access token")
	}

	var userInfo *userInfo
	err = util.Retry(func() error {
		var err error
		userInfo, err = svc.getUserInfo(accessToken, language.IETF())
		return err
	}, svc.retryTimes)
	if err != nil {
		return AccountMetadata{}, errors.Wrap(err, "can't get user info")
	}

	var splatoonAccessToken, nsName string
	err = util.Retry(func() error {
		var err error
		splatoonAccessToken, nsName, err = svc.getSplatoonAccessToken(accessToken, userInfo, language.IETF())
		return err
	}, svc.retryTimes)
	if err != nil {
		return AccountMetadata{}, errors.Wrap(err, "can't get splatoon access token")
	}

	var iksmSession string
	err = util.Retry(func() error {
		var err error
		iksmSession, err = svc.getIksmSession(splatoonAccessToken, language.IETF())
		return err
	}, svc.retryTimes)
	if err != nil {
		return AccountMetadata{}, errors.Wrap(err, "can't get iksm session")
	}

	return AccountMetadata{
		IKSM:        iksmSession,
		AccountName: userInfo.NickName,
		UserName:    nsName,
	}, nil
}

const (
	prefix    = "npf71b963c1b7b6d119://auth#"
	prefixLen = len(prefix)
)

func (svc *impl) getSessionTokenCode(link string) (string, error) {
	param, err := url.ParseQuery(link[prefixLen:])
	if err != nil {
		return "", errors.Wrap(err, "can't parse redirect link")
	}
	sessionTokenCode := param.Get("session_token_code")
	return sessionTokenCode, nil
}

func (svc *impl) getSessionToken(proofKey []byte, sessionTokenCode string, acceptLang string) (string, error) {
	reqURL := "https://accounts.nintendo.com/connect/1.0.0/api/session_token"
	bodyMap := map[string][]string{
		"client_id":                   {"71b963c1b7b6d119"},
		"session_token_code":          {sessionTokenCode},
		"session_token_code_verifier": {string(proofKey)},
	}
	bodyText := url.Values(bodyMap).Encode()
	reqBody := strings.NewReader(bodyText)
	req, err := http.NewRequest("POST", reqURL, reqBody)
	if err != nil {
		return "", errors.Wrap(err, "can't generate request")
	}
	req.Header = map[string][]string{
		"Accept":          {"application/json"},
		"Accept-Encoding": {"gzip"},
		"Accept-Language": {"en-US"},
		"Connection":      {"Keep-Alive"},
		"Content-Length":  {strconv.FormatInt(int64(len(bodyText)), 10)},
		"Content-Type":    {"application/x-www-form-urlencoded"},
		"Host":            {"accounts.nintendo.com"},
		"User-Agent":      {"OnlineLounge/1.10.1 NASDKAPI Android"},
	}
	resp, err := svc.client.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "can't get response")
	}
	if resp.StatusCode != http.StatusOK {
		n, _ := ioutil.ReadAll(resp.Body)
		log.Debug("get session token", zap.ByteString("json", n))
		return "", fmt.Errorf("status code not 200, got %d", resp.StatusCode)
	}
	defer closeBody(resp.Body)
	respBody := resp.Body
	if isGzip(resp.Header) {
		respBody, err = gzip.NewReader(respBody)
		if err != nil {
			return "", errors.Wrap(err, "can't unzip response body")
		}
	}
	respJSON, err := ioutil.ReadAll(respBody)
	if err != nil {
		return "", errors.Wrap(err, "can't read response body")
	}
	sessionToken := json.Get(respJSON, "session_token").ToString()
	log.Debug("get session token", zap.String("session token", sessionToken), zap.ByteString("json", respJSON))
	return sessionToken, nil
}

func (svc *impl) getAccessToken(sessionToken string, acceptLang string) (string, error) {
	reqURL := "https://accounts.nintendo.com/connect/1.0.0/api/token"
	bodyMap := map[string]string{
		"client_id":     "71b963c1b7b6d119",
		"session_token": sessionToken,
		"grant_type":    "urn:ietf:params:oauth:grant-type:jwt-bearer-session-token",
	}
	bodyText, err := json.Marshal(bodyMap)
	if err != nil {
		return "", errors.Wrap(err, "can't generate request body")
	}
	reqBody := bytes.NewReader(bodyText)
	req, err := http.NewRequest("POST", reqURL, reqBody)
	if err != nil {
		return "", errors.Wrap(err, "can't generate request")
	}
	req.Header = map[string][]string{
		"Accept":          {"application/json"},
		"Accept-Encoding": {"gzip"},
		"Accept-Language": {acceptLang},
		"Connection":      {"Keep-Alive"},
		"Content-Length":  {strconv.FormatInt(int64(len(bodyText)), 10)},
		"Content-Type":    {"application/json; charset=utf-8"},
		"Host":            {"accounts.nintendo.com"},
		"User-Agent":      {"OnlineLounge/1.10.1 NASDKAPI Android"},
	}
	resp, err := svc.client.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "can't get response")
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status code not 200, got %d", resp.StatusCode)
	}
	defer closeBody(resp.Body)
	respBody := resp.Body
	if isGzip(resp.Header) {
		respBody, err = gzip.NewReader(respBody)
		if err != nil {
			return "", errors.Wrap(err, "can't unzip response body")
		}
	}
	respJSON, err := ioutil.ReadAll(respBody)
	if err != nil {
		return "", errors.Wrap(err, "can't read response body")
	}
	accessToken := json.Get(respJSON, "access_token").ToString()
	log.Debug("get access token", zap.String("access token", accessToken), zap.ByteString("json", respJSON))
	return accessToken, nil
}

type userInfo struct {
	NickName string `json:"nickname"`
	Country  string `json:"country"`
	Birthday string `json:"birthday"`
	Language string `json:"language"`
}

func (svc *impl) getUserInfo(accessToken string, acceptLang string) (*userInfo, error) {
	reqURL := "https://api.accounts.nintendo.com/2.0.0/users/me"
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "can't generate request")
	}
	req.Header = map[string][]string{
		"Accept":          {"application/json"},
		"Accept-Language": {acceptLang},
		"Accept-Encoding": {"gzip"},
		"Authorization":   {"Bearer " + accessToken},
		"Host":            {"api.accounts.nintendo.com"},
		"Connection":      {"Keep-Alive"},
		"User-Agent":      {"OnlineLounge/1.10.1 NASDKAPI Android"},
	}
	resp, err := svc.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "can't get response")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code not 200, got %d", resp.StatusCode)
	}
	defer closeBody(resp.Body)
	respBody := resp.Body
	if isGzip(resp.Header) {
		respBody, err = gzip.NewReader(respBody)
		if err != nil {
			return nil, errors.Wrap(err, "can't unzip response body")
		}
	}
	respJSON, err := ioutil.ReadAll(respBody)
	if err != nil {
		return nil, errors.Wrap(err, "can't read response body")
	}
	userInfo := &userInfo{}
	err = json.Unmarshal(respJSON, userInfo)
	if err != nil {
		return nil, errors.Wrap(err, "can't unmarshal response body")
	}
	log.Debug("get user info",
		zap.String("nickname", userInfo.NickName),
		zap.String("country", userInfo.Country),
		zap.String("birthday", userInfo.Birthday),
		zap.String("language", userInfo.Language),
		zap.ByteString("json", respJSON))
	return userInfo, nil
}

type flapgResponse struct {
	Result flapgResponseResult `json:"result"`
}

type flapgResponseResult struct {
	F  string `json:"f"`
	P1 string `json:"p1"`
	P2 string `json:"p2"`
	P3 string `json:"p3"`
}

func (svc *impl) getFlapgResponse(guid, accessToken string, timestamp int64, iid string) (*flapgResponse, error) {
	hash, err := svc.getS2SResponse(accessToken, timestamp)
	if err != nil {
		return nil, errors.Wrap(err, "can't get hash")
	}
	reqURL := "https://flapg.com/ika2/api/login?public"
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "can't generate request")
	}
	req.Header = map[string][]string{
		"x-token": {accessToken},
		"x-time":  {strconv.FormatInt(timestamp, 10)},
		"x-guid":  {guid},
		"x-hash":  {hash},
		"x-ver":   {"3"},
		"x-iid":   {iid},
	}
	resp, err := svc.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "can't get response")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code not 200, got %d", resp.StatusCode)
	}
	defer closeBody(resp.Body)
	respJSON, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "can't read response body")
	}
	flapgResponse := &flapgResponse{}
	err = json.Unmarshal(respJSON, flapgResponse)
	if err != nil {
		return nil, errors.Wrap(err, "can't unmarshal response body")
	}
	log.Debug("get flapgResponse",
		zap.String("f", flapgResponse.Result.F),
		zap.String("p1", flapgResponse.Result.P1),
		zap.String("p2", flapgResponse.Result.P2),
		zap.String("p3", flapgResponse.Result.P3),
		zap.ByteString("json", respJSON))
	return flapgResponse, nil
}

func (svc *impl) getS2SResponse(accessToken string, timestamp int64) (string, error) {
	reqURL := "https://elifessler.com/s2s/api/gen2"
	bodyMap := map[string][]string{
		"naIdToken": {accessToken},
		"timestamp": {strconv.FormatInt(timestamp, 10)},
	}
	bodyText := url.Values(bodyMap).Encode()
	reqBody := strings.NewReader(bodyText)
	req, err := http.NewRequest("POST", reqURL, reqBody)
	if err != nil {
		return "", errors.Wrap(err, "can't generate request")
	}
	req.Header = map[string][]string{
		"User-Agent": {"splatnet2statink/1.5.8"}, // todo: use my own agent?
	}
	resp, err := svc.client.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "can't get response")
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status code not 200, got %d", resp.StatusCode)
	}
	defer closeBody(resp.Body)
	respJSON, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "can't read response body")
	}
	hash := json.Get(respJSON, "hash").ToString()
	log.Debug("get hash", zap.String("hash", hash), zap.ByteString("json", respJSON))
	return hash, nil
}

func (svc *impl) getSplatoonAccessTokenFirstStep(flapgNsoResponse *flapgResponse, userInfo *userInfo, acceptLang string) (string, string, error) {
	reqURL := "https://api-lp1.znc.srv.nintendo.net/v1/Account/Login"
	bodyMap := map[string]map[string]string{
		"parameter": {
			"f":          flapgNsoResponse.Result.F,
			"naIdToken":  flapgNsoResponse.Result.P1,
			"timestamp":  flapgNsoResponse.Result.P2,
			"requestId":  flapgNsoResponse.Result.P3,
			"naCountry":  userInfo.Country,
			"naBirthday": userInfo.Birthday,
			"language":   userInfo.Language,
		},
	}
	bodyText, err := json.Marshal(bodyMap)
	if err != nil {
		return "", "", errors.Wrap(err, "can't generate request body")
	}
	reqBody := bytes.NewReader(bodyText)
	req, err := http.NewRequest("POST", reqURL, reqBody)
	if err != nil {
		return "", "", errors.Wrap(err, "can't generate request")
	}
	req.Header = map[string][]string{
		"Accept":           {"application/json"},
		"Accept-Encoding":  {"gzip"},
		"Accept-Language":  {acceptLang},
		"Authorization":    {"Bearer"},
		"Connection":       {"Keep-Alive"},
		"Content-Length":   {strconv.FormatInt(int64(len(bodyText)), 10)},
		"Content-Type":     {"application/json; charset=utf-8"},
		"Host":             {"api-lp1.znc.srv.nintendo.net"},
		"User-Agent":       {"com.nintendo.znca/1.10.1 (Android/7.1.2)"},
		"X-Platform":       {"Android"},
		"X-ProductVersion": {"1.10.1"},
	}
	resp, err := svc.client.Do(req)
	if err != nil {
		return "", "", errors.Wrap(err, "can't get response")
	}
	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("status code not 200, got %d", resp.StatusCode)
	}
	defer closeBody(resp.Body)
	respBody := resp.Body
	if isGzip(resp.Header) {
		respBody, err = gzip.NewReader(respBody)
		if err != nil {
			return "", "", errors.Wrap(err, "can't unzip response body")
		}
	}
	respJSON, err := ioutil.ReadAll(respBody)
	if err != nil {
		return "", "", errors.Wrap(err, "can't read response body")
	}
	splatoonAccessToken := json.Get(respJSON, "result", "webApiServerCredential", "accessToken").ToString()
	nsName := json.Get(respJSON, "result", "user", "name").ToString()
	log.Debug("get splatoon access token first step",
		zap.String("splatoon webApiServerCredential access token", splatoonAccessToken),
		zap.ByteString("json", respJSON))
	return splatoonAccessToken, nsName, nil
}

func (svc *impl) getSplatoonAccessTokenSecondStep(accessToken string, flapgAppResponse *flapgResponse, acceptLang string) (string, error) {
	reqURL := "https://api-lp1.znc.srv.nintendo.net/v2/Game/GetWebServiceToken"
	bodyMap := map[string]map[string]interface{}{
		"parameter": {
			"id":                int64(5741031244955648),
			"f":                 flapgAppResponse.Result.F,
			"registrationToken": flapgAppResponse.Result.P1,
			"timestamp":         flapgAppResponse.Result.P2,
			"requestId":         flapgAppResponse.Result.P3,
		},
	}
	bodyText, err := json.Marshal(bodyMap)
	if err != nil {
		return "", errors.Wrap(err, "can't generate request body")
	}
	reqBody := bytes.NewReader(bodyText)
	req, err := http.NewRequest("POST", reqURL, reqBody)
	if err != nil {
		return "", errors.Wrap(err, "can't generate request")
	}
	req.Header = map[string][]string{
		"Accept":           {"application/json"},
		"Accept-Encoding":  {"gzip"},
		"Accept-Language":  {acceptLang},
		"Authorization":    {"Bearer " + accessToken},
		"Connection":       {"Keep-Alive"},
		"Content-Length":   {strconv.FormatInt(int64(len(bodyText)), 10)},
		"Content-Type":     {"application/json; charset=utf-8"},
		"Host":             {"api-lp1.znc.srv.nintendo.net"},
		"User-Agent":       {"com.nintendo.znca/1.10.1 (Android/7.1.2)"},
		"X-Platform":       {"Android"},
		"X-ProductVersion": {"1.10.1"},
	}
	resp, err := svc.client.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "can't get response")
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status code not 200, got %d", resp.StatusCode)
	}
	defer closeBody(resp.Body)
	respBody := resp.Body
	if isGzip(resp.Header) {
		respBody, err = gzip.NewReader(respBody)
		if err != nil {
			return "", errors.Wrap(err, "can't unzip response body")
		}
	}
	respJSON, err := ioutil.ReadAll(respBody)
	if err != nil {
		return "", errors.Wrap(err, "can't read response body")
	}
	splatoonAccessToken := json.Get(respJSON, "result", "accessToken").ToString()
	log.Debug("get splatoon access token",
		zap.String("splatoon access token", splatoonAccessToken),
		zap.ByteString("json", respJSON))
	return splatoonAccessToken, nil
}

func (svc *impl) getSplatoonAccessToken(accessToken string, userInfo *userInfo, acceptLang string) (string, string, error) {
	uuid4, err := uuid.NewRandom()
	if err != nil {
		return "", "", errors.Wrap(err, "can't generate uuid4")
	}
	guid := uuid4.String()
	timestamp := time.Now().Unix()

	flapgResponse, err := svc.getFlapgResponse(guid, accessToken, timestamp, "nso")
	if err != nil {
		return "", "", errors.Wrap(err, "can't get flapg response")
	}
	firstSplatoonAccessToken, name, err := svc.getSplatoonAccessTokenFirstStep(flapgResponse, userInfo, acceptLang)
	if err != nil {
		return "", "", errors.Wrap(err, "can't get first splatoon access token")
	}

	flapgResponse, err = svc.getFlapgResponse(guid, firstSplatoonAccessToken, timestamp, "app")
	if err != nil {
		return "", "", errors.Wrap(err, "can't get flapg response")
	}

	SecondSplatoonAccessToken, err := svc.getSplatoonAccessTokenSecondStep(firstSplatoonAccessToken, flapgResponse, acceptLang)
	if err != nil {
		return "", "", errors.Wrap(err, "can't get second splatoon access token")
	}

	return SecondSplatoonAccessToken, name, nil
}

func (svc *impl) getIksmSession(splatoonAccessToken string, acceptLang string) (string, error) {
	reqURL := "https://app.splatoon2.nintendo.net/?lang=" + acceptLang
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return "", errors.Wrap(err, "can't generate request")
	}
	req.Header = map[string][]string{
		"Accept":                  {"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"},
		"Accept-Encoding":         {"gzip"},
		"Accept-Language":         {acceptLang},
		"Connection":              {"Keep-Alive"},
		"DNT":                     {"0"},
		"Host":                    {"app.splatoon2.nintendo.net"},
		"User-Agent":              {"Mozilla/5.0 (Linux; Android 7.1.2; Pixel Build/NJH47D; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/59.0.3071.125 Mobile Safari/537.36"},
		"X-GameWebToken":          {splatoonAccessToken},
		"X-IsAppAnalyticsOptedIn": {"false"},
		"X-IsAnalyticsOptedIn":    {"false"},
		"X-Requested-With":        {"com.nintendo.znca"},
	}
	resp, err := svc.client.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "can't get response")
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status code not 200, got %d", resp.StatusCode)
	}
	defer closeBody(resp.Body)
	cookies := resp.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "iksm_session" {
			iksmSession := cookie.Value
			log.Debug("get iksm session", zap.String("iksm session", iksmSession))
			return iksmSession, nil
		}
	}
	return "", errors.Errorf("iksm_session not in response's cookies")
}
