package nintendo

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"github.com/google/uuid"
	json "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	log "telegram-splatoon2-bot/logger"
	"time"
)

func NewProofKey() ([]byte, error) {
	proofKey, err := randBytes(32)
	if err != nil {
		return nil, errors.Wrap(err, "can't generate proof key")
	}
	return base64UrlEncode(proofKey), nil
}

func NewLoginLink(proofKey []byte) (string, error) {
	state, err := randBytes(36)
	if err != nil {
		return "", errors.Wrap(err, "can't generate status")
	}
	hashProofKey := sha256.Sum256(proofKey)
	challenge := string(base64UrlEncode(hashProofKey[:]))
	authState := string(base64UrlEncode(state))

	hardcodeUrl := "https://accounts.nintendo.com/connect/1.0.0/authorize?" +
		"redirect_uri=npf71b963c1b7b6d119://auth" +
		"&client_id=71b963c1b7b6d119" +
		"&scope=openid%20user%20user.birthday%20user.mii%20user.screenName" +
		"&response_type=session_token_code" +
		"&session_token_code_challenge_method=S256" +
		"&theme=login_form" +
		"&state=" + authState +
		"&session_token_code_challenge=" + challenge

	return hardcodeUrl, nil
}

var prefix = "npf71b963c1b7b6d119://auth#"
var prefixLen = len(prefix)

func GetCookies(link string, proofKey []byte) (string, error) {
	if !strings.HasPrefix(link, prefix) {
		return "", fmt.Errorf("unknown URI")
	}
	param, err := url.ParseQuery(link[prefixLen:])
	if err != nil {
		return "", errors.Wrap(err, "can't parse redirect link")
	}
	//state := param.Get("state")
	//sessionState := param.Get("session_state")
	sessionTokenCode := param.Get("session_token_code")

	language := "en-US"
	sessionToken, err := getSessionToken(proofKey, sessionTokenCode, language)
	if err != nil {
		return "", errors.Wrap(err, "can't get session token")
	}
	accessToken, err := getAccessToken(sessionToken, language)
	if err != nil {
		return "", errors.Wrap(err, "can't get access token")
	}
	userInfo, err := getUserInfo(accessToken, language)
	if err != nil {
		return "", errors.Wrap(err, "can't get user info")
	}
	splatoonAccessToken, err := getSplatoonAccessToken(accessToken, userInfo, language)
	if err != nil {
		return "", errors.Wrap(err, "can't get splatoon access token")
	}
	iksmSession, err := getIksmSession(splatoonAccessToken, language)
	if err != nil {
		return "", errors.Wrap(err, "can't get iksm session")
	}
	return iksmSession, nil
}

func getSessionToken(proofKey []byte, sessionTokenCode string, acceptLang string) (string, error) {
	reqUrl := "https://accounts.nintendo.com/connect/1.0.0/api/session_token"
	bodyMap := map[string][]string{
		"client_id":                   {"71b963c1b7b6d119"},
		"session_token_code":          {sessionTokenCode},
		"session_token_code_verifier": {string(proofKey)},
	}
	bodyText := url.Values(bodyMap).Encode()
	reqBody := strings.NewReader(bodyText)
	req, err := http.NewRequest("POST", reqUrl, reqBody)
	if err != nil {
		return "", errors.Wrap(err, "can't generate request")
	}
	req.Header = map[string][]string{
		"Accept":          {"application/json"},
		"Accept-Encoding": {"gzip"},
		"Accept-Language": {acceptLang},
		"Connection":      {"Keep-Alive"},
		"Content-Length":  {strconv.FormatInt(int64(len(bodyText)), 10)},
		"Content-Type":    {"application/x-www-form-urlencoded"},
		"Host":            {"accounts.nintendo.com"},
		"User-Agent":      {"OnlineLounge/1.9.0 NASDKAPI Android"},
	}
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
	sessionToken := json.Get(respJson, "session_token").ToString()
	log.Debug("Get session token", zap.String("session token", sessionToken), zap.ByteString("json", respJson))
	return sessionToken, nil
}

func getAccessToken(sessionToken string, acceptLang string) (string, error) {
	reqUrl := "https://accounts.nintendo.com/connect/1.0.0/api/token"
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
	req, err := http.NewRequest("POST", reqUrl, reqBody)
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
		"User-Agent":      {"OnlineLounge/1.9.0 NASDKAPI Android"},
	}
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
	accessToken := json.Get(respJson, "access_token").ToString()
	log.Debug("Get access token", zap.String("access token", accessToken), zap.ByteString("json", respJson))
	return accessToken, nil
}

type UserInfo struct {
	NickName string `json:"nickname"`
	Country  string `json:"country"`
	Birthday string `json:"birthday"`
	Language string `json:"language"`
}

func getUserInfo(accessToken string, acceptLang string) (*UserInfo, error) {
	reqUrl := "https://api.accounts.nintendo.com/2.0.0/users/me"
	req, err := http.NewRequest("GET", reqUrl, nil)
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
		"User-Agent":      {"OnlineLounge/1.9.0 NASDKAPI Android"},
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "can't get response")
	}
	defer closeBody(resp.Body)
	respBody, err := gzip.NewReader(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "can't unzip response body")
	}
	respJson, err := ioutil.ReadAll(respBody)
	if err != nil {
		return nil, errors.Wrap(err, "can't read response body")
	}
	userInfo := &UserInfo{}
	err = json.Unmarshal(respJson, userInfo)
	if err != nil {
		return nil, errors.Wrap(err, "can't unmarshal response body")
	}
	log.Debug("Get user info",
		zap.String("nickname", userInfo.NickName),
		zap.String("country", userInfo.Country),
		zap.String("birthday", userInfo.Birthday),
		zap.String("language", userInfo.Language),
		zap.ByteString("json", respJson))
	return userInfo, nil
}

type FlapgResponse struct {
	Result FlapgResponseResult `json:"result"`
}

type FlapgResponseResult struct {
	F  string `json:"f"`
	P1 string `json:"p1"`
	P2 string `json:"p2"`
	P3 string `json:"p3"`
}

func getFlapgResponse(guid, accessToken string, timestamp int64, iid string) (*FlapgResponse, error) {
	hash, err := getS2SResponse(accessToken, timestamp)
	if err != nil {
		return nil, errors.Wrap(err, "can't get hash")
	}
	reqUrl := "https://flapg.com/ika2/api/login?public"
	req, err := http.NewRequest("GET", reqUrl, nil)
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
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "can't get response")
	}
	defer closeBody(resp.Body)
	respJson, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "can't read response body")
	}
	flapgResponse := &FlapgResponse{}
	err = json.Unmarshal(respJson, flapgResponse)
	if err != nil {
		return nil, errors.Wrap(err, "can't unmarshal response body")
	}
	log.Debug("Get flapgResponse",
		zap.String("f", flapgResponse.Result.F),
		zap.String("p1", flapgResponse.Result.P1),
		zap.String("p2", flapgResponse.Result.P2),
		zap.String("p3", flapgResponse.Result.P3),
		zap.ByteString("json", respJson))
	return flapgResponse, nil
}

func getS2SResponse(accessToken string, timestamp int64) (string, error) {
	reqUrl := "https://elifessler.com/s2s/api/gen2"
	bodyMap := map[string][]string{
		"naIdToken": {accessToken},
		"timestamp": {strconv.FormatInt(timestamp, 10)},
	}
	bodyText := url.Values(bodyMap).Encode()
	reqBody := strings.NewReader(bodyText)
	req, err := http.NewRequest("POST", reqUrl, reqBody)
	if err != nil {
		return "", errors.Wrap(err, "can't generate request")
	}
	req.Header = map[string][]string{
		"User-Agent": {"splatnet2statink/1.5.6"}, // todo: use my own agent?
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "can't get response")
	}
	defer closeBody(resp.Body)
	respJson, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "can't read response body")
	}
	hash := json.Get(respJson, "hash").ToString()
	log.Debug("Get hash", zap.String("hash", hash), zap.ByteString("json", respJson))
	return hash, nil
}

func getSplatoonAccessTokenFirstStep(flapgNsoResponse *FlapgResponse, userInfo *UserInfo, acceptLang string) (string, error) {
	reqUrl := "https://api-lp1.znc.srv.nintendo.net/v1/Account/Login"
	bodyMap:= map[string]map[string]string{
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
		return "", errors.Wrap(err, "can't generate request body")
	}
	reqBody := bytes.NewReader(bodyText)
	req, err := http.NewRequest("POST", reqUrl, reqBody)
	if err != nil {
		return "", errors.Wrap(err, "can't generate request")
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
		"User-Agent":       {"com.nintendo.znca/1.9.0 (Android/7.1.2)"},
		"X-Platform":       {"Android"},
		"X-ProductVersion": {"1.9.0"},
	}
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
	splatoonAccessToken := json.Get(respJson, "result", "webApiServerCredential", "accessToken").ToString()
	log.Debug("Get splatoon access token first step",
		zap.String("splatoon webApiServerCredential access token", splatoonAccessToken),
		zap.ByteString("json", respJson))
	return splatoonAccessToken, nil
}

func getSplatoonAccessTokenSecondStep(accessToken string, flapgAppResponse *FlapgResponse, acceptLang string) (string, error) {
	reqUrl := "https://api-lp1.znc.srv.nintendo.net/v2/Game/GetWebServiceToken"
	bodyMap := map[string]map[string]interface{}{
		"parameter": {
			"id":                5741031244955648,
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
	req, err := http.NewRequest("POST", reqUrl, reqBody)
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
		"User-Agent":       {"com.nintendo.znca/1.9.0 (Android/7.1.2)"},
		"X-Platform":       {"Android"},
		"X-ProductVersion": {"1.9.0"},
	}
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
	splatoonAccessToken := json.Get(respJson, "result", "accessToken").ToString()
	log.Debug("Get splatoon access token",
		zap.String("splatoon access token", splatoonAccessToken),
		zap.ByteString("json", respJson))
	return splatoonAccessToken, nil
}

func getSplatoonAccessToken(accessToken string, userInfo *UserInfo, acceptLang string) (string, error) {
	uuid4, err := uuid.NewRandom()
	if err != nil {
		return "", errors.Wrap(err, "can't generate uuid4")
	}
	guid := uuid4.URN()[9:]
	timestamp := time.Now().Unix()

	flapgResponse, err := getFlapgResponse(guid, accessToken, timestamp, "nso")
	if err != nil {
		return "", errors.Wrap(err, "can't get flapg response")
	}
	firstSplatoonAccessToken, err := getSplatoonAccessTokenFirstStep(flapgResponse, userInfo, acceptLang)
	if err != nil {
		return "", errors.Wrap(err, "can't get first splatoon access token")
	}

	flapgResponse, err = getFlapgResponse(guid, firstSplatoonAccessToken, timestamp, "app")
	if err != nil {
		return "", errors.Wrap(err, "can't get flapg response")
	}

	SecondSplatoonAccessToken, err := getSplatoonAccessTokenSecondStep(firstSplatoonAccessToken, flapgResponse, acceptLang)
	if err != nil {
		return "", errors.Wrap(err, "can't get second splatoon access token")
	}

	return SecondSplatoonAccessToken, nil
}

func getIksmSession(splatoonAccessToken string, acceptLang string) (string, error) {
	reqUrl := "https://app.splatoon2.nintendo.net/?lang=" + acceptLang
	req, err := http.NewRequest("GET", reqUrl, nil)
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
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "can't get response")
	}
	defer closeBody(resp.Body)
	cookies := resp.Cookies()
	for _, cookie := range cookies{
		if cookie.Name == "iksm_session" {
			iksmSession := cookie.Value
			log.Debug("Get iksm session", zap.String("iksm session", iksmSession))
			return iksmSession, nil
		}
	}
	return "", fmt.Errorf("iksm_session not in response's cookies")
}
