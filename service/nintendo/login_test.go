package nintendo

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	json "github.com/json-iterator/go"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	log "telegram-splatoon2-bot/common/log"
)

func prepareTest() {
	viper.SetConfigName("dev")
	viper.SetConfigType("json")
	viper.AddConfigPath("../config/")
	viper.AddConfigPath("./config/")
	viper.ReadInConfig()
	log.InitLogger()
	InitClient()
}

var (
	proofKey            = []byte{0}
	acceptLang          = "en-US"
	sessionTokenCode    = ""
	sessionToken        = ""
	accessToken         = ""
	guid                = ""
	timestamp           = int64(0)
	iid                 = ""
	splatoonAccessToken = ""
	flapgNsoResponse    = &flapgResponse{Result: flapgResponseResult{
		F:  "e294a519a2b08a9de23ccd769e1e731994928b7eb94f9b6e3463ae09f8ddf9005b5bcbfb77d446cd2b0d757c65",
		P1: "eyJhbGciOiJSUzI1NiIsImprdSI6Imh0dHBzOi8vYWNjb3VudHMubmludGVuZG8uY29tLzEuMC4wL2NlcnRpZmljYXRlcyIsImtpZCI6IjZmMGY3ZWM4LWI2NGQtNGFmMC1iZTk0LThhZWVjMWVmOGIyZSJ9.eyJleHAiOjE2MDE0Mjg0ODAsImFjOnNjcCI6WzAsOCw5LDE3LDIzXSwiYWM6Z3J0Ijo2NCwidHlwIjoidG9rZW4iLCJqdGkiOiJmYzQxYjFlYS01ZDM1LTQwYzctYjgyMC01YzZiNWNiZTg0N2EiLCJhdWQiOiI3MWI5NjNjMWI3YjZkMTE5IiwiaXNzIjoiaHR0cHM6Ly9hY2NvdW50cy5uaW50ZW5kby5jb20iLCJzdWIiOiJkYTQzMjBlMzUyNWI1OWExIiwiaWF0IjoxNjAxNDI3NTgwfQ.MyVIClEGK-7SB-N4O_5Rr4CuVpNVcYbMyB5zSnjVHivrSEUlFWoIwkRk8qE8FLQjHuAZ8TIStO8iAspY5dfdcPbpDC356Q4-JBS0MqeXzE6y3LZgrC1B26Oodq3CNfiMF0K2XIfmthns8IzsMCUjDrU8wmISGTQj6rUAuDpUdDwUIU5UgSQTLV7-Vv5EorvVlGiOU7Xqc4qoR29pS9CAKonC_tWyVvrcEQIk2nXaKexCZVdZe2Qqa3sFO70EBf0rlbeajfWQFxV6VCOoi_MM3OH4HWN1zF6ZYwAekVVH231VXY8a_VeRSLG_Bbp19F7z4Y5YduS3iQdZ0hMQhH5cyQ",
		P2: "1601427580",
		P3: "c9515e7d-6f09-4b12-b605-2b859ba0267b",
	}}
	flapgAppResponse = &flapgResponse{Result: flapgResponseResult{
		F:  "",
		P1: "",
		P2: "",
		P3: "",
	}}
	userInfo = &userInfo{
		NickName: "",
		Country:  "",
		Birthday: "",
		Language: "",
	}
)

func TestGetSessionToken(t *testing.T) {
	prepareTest()
	_, err := getSessionToken(proofKey, sessionTokenCode, acceptLang)
	assert.Nil(t, err)
}

func TestGetAccessToken(t *testing.T) {
	prepareTest()
	_, err := getAccessToken(sessionToken, acceptLang)
	assert.Nil(t, err)
}

func TestGetUserInfo(t *testing.T) {
	prepareTest()
	_, err := getUserInfo(accessToken, acceptLang)
	assert.Nil(t, err)
}

func TestGetFlapgResponse(t *testing.T) {
	prepareTest()
	_, err := getFlapgResponse(guid, accessToken, timestamp, iid)
	assert.Nil(t, err)
}

func TestGetS2SResponse(t *testing.T) {
	prepareTest()
	_, err := getS2SResponse(accessToken, timestamp)
	assert.Nil(t, err)
}

func TestGetSplatoonAccessTokenFirstStep(t *testing.T) {
	prepareTest()
	_, _, err := getSplatoonAccessTokenFirstStep(flapgNsoResponse, userInfo, acceptLang)
	assert.Nil(t, err)
}

func TestGetSplatoonAccessTokenSecondStep(t *testing.T) {
	prepareTest()
	_, err := getSplatoonAccessTokenSecondStep(splatoonAccessToken, flapgAppResponse, acceptLang)
	assert.Nil(t, err)
}

func BenchmarkStringConcat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bodyText := `
{
	"parameter": {
		"f": "` + flapgNsoResponse.Result.F + `",
		"naIdToken": "` + flapgNsoResponse.Result.P1 + `",
		"timestamp": "` + flapgNsoResponse.Result.P2 + `",
		"requestId": "` + flapgNsoResponse.Result.P3 + `",
		"naCountry": "` + userInfo.Country + `",
		"naBirthday": "` + userInfo.Birthday + `",
		"language": "` + userInfo.Language + `"
	}
}`
		reqBody := strings.NewReader(bodyText)
		_, _ = ioutil.ReadAll(reqBody)
	}
}

func BenchmarkJsonEncode(b *testing.B) {
	for i := 0; i < b.N; i++ {
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
		bodyText, _ := json.Marshal(bodyMap)
		reqBody := bytes.NewReader(bodyText)
		_, _ = ioutil.ReadAll(reqBody)
	}
}
