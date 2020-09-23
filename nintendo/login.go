package nintendo

import (
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

func GetLoginLink() (string, error) {
	state, err := randBytes(36)
	if err != nil {
		return "", errors.Wrap(err, "can't generate status")
	}
	proofKey, err := randBytes(32)
	if err != nil {
		return "", errors.Wrap(err, "can't generate proof key")
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

func getSessionToken() error {
	state, err := randBytes(36)
	if err != nil {
		return errors.Wrap(err, "can't generate status")
	}
	proofKey, err := randBytes(32)
	if err != nil {
		return errors.Wrap(err, "can't generate proof key")
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
	req, err := http.NewRequest("GET", hardcodeUrl, nil)
	if err != nil {
		return errors.Wrap(err, "can't generate request")
	}
	req.Header = map[string][]string{
		"Host":                      {"accounts.nintendo.com"},
		"Connection":                {"keep-alive"},
		"Cache-Control":             {"max-age=0"},
		"Upgrade-Insecure-Requests": {"1"},
		"User-Agent":                {"Mozilla/5.0 (Linux; Android 7.1.2; Pixel Build/NJH47D; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/59.0.3071.125 Mobile Safari/537.36"},
		"Accept":                    {"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8n"},
		"DNT":                       {"1"},
		"Accept-Encoding":           {"gzip"},
	}

	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "can't get response")
	}
	defer closeBody(resp.Body)
	body, err := gzip.NewReader(resp.Body)
	if err != nil {
		return errors.Wrap(err, "can't unzip response body")
	}
	ret, _ := ioutil.ReadAll(body)
	fmt.Println(string(ret))

	return nil
}
