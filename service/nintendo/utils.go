package nintendo

import (
	"bytes"
	"compress/gzip"
	"crypto/rand"
	"encoding/base64"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	json "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"telegram-splatoon2-bot/common/util"
)

// base64UrlEncode encodes a []byte to a base64 coding url
// which replace '+' to '-', '/' to '_', and omit the padding '='
func base64UrlEncode(src []byte) []byte {
	base64Url := make([]byte, base64.StdEncoding.EncodedLen(len(src)))
	base64.StdEncoding.Encode(base64Url, src)
	index := bytes.LastIndex(base64Url, []byte{'='})
	if index > 0 {
		base64Url = base64Url[:index]
	}
	base64Url = bytes.ReplaceAll(base64Url, []byte{'+'}, []byte{'-'})
	base64Url = bytes.ReplaceAll(base64Url, []byte{'/'}, []byte{'_'})

	return base64Url
}

// base64UrlDecode decodes a base64 coding url to []byte,
// which replace '-' to '+', '_' to '/', and pad '='
func base64UrlDecode(base64Url []byte) ([]byte, error) {
	padding := 4 - (len(base64Url) % 4)
	if padding > 2 {
		padding = 0
	}
	temp := make([]byte, len(base64Url)+padding)
	copy(temp, base64Url)
	for i := 0; i < padding; i += 1 {
		temp[len(base64Url)+i] = '='
	}
	temp = bytes.ReplaceAll(temp, []byte{'_'}, []byte{'/'})
	temp = bytes.ReplaceAll(temp, []byte{'-'}, []byte{'+'})
	ret := make([]byte, base64.StdEncoding.DecodedLen(len(temp)))
	_, err := base64.StdEncoding.Decode(ret, temp)
	if err != nil {
		return nil, errors.Wrap(err, "can't decode base64 url")
	}
	return ret[:len(ret)-padding], nil
}

func randBytes(n int) ([]byte, error) {
	buf := make([]byte, n)
	_, err := rand.Read(buf)
	if err != nil {
		return nil, errors.Wrap(err, "can't generate random bytes slice")
	}
	return buf, nil
}

// closeBody closes http response body to eliminate the annoying warning of IDE
func closeBody(c io.Closer) {
	_ = c.Close()
}

func getAppHeader(iksm string, timezone int, acceptLang string, gzip bool) map[string][]string {
	acceptEncoding := ""
	if gzip {
		acceptEncoding = "gzip"
	}
	return map[string][]string{
		"Accept":            {"*/*"},
		"Accept-Encoding":   {acceptEncoding},
		"Accept-Language":   {acceptLang},
		"Connection":        {"Keep-Alive"},
		"Host":              {"app.splatoon2.nintendo.net"},
		"Referer":           {"https://app.splatoon2.nintendo.net/home"},
		"User-Agent":        {"Mozilla/5.0 (Linux; Android 7.1.2; Pixel Build/NJH47D; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/59.0.3071.125 Mobile Safari/537.36"},
		"x-requested-with":  {"XMLHttpRequest"},
		"x-timezone-offset": {strconv.Itoa(timezone)}, // todo: seems useless
		"x-unique-id":       {"32449507786579989234"},
		"Cookie":            {"iksm_session=" + iksm},
	}
}

func isCookiesExpired(respJson []byte) bool {
	return json.Get(respJson, "code").ToString() == "AUTHENTICATION_ERROR"
}

func (svc *impl) getSplatoon2RestfulJson(url string, iksm string, timezone int, acceptLang string) ([]byte, error) {
	var respJson []byte
	err := util.Retry(func() error {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return errors.Wrap(err, "can't generate request")
		}
		req.Header = getAppHeader(iksm, timezone, acceptLang, true)
		resp, err := svc.client.Do(req)
		if err != nil {
			return errors.Wrap(err, "can't get response")
		}
		defer closeBody(resp.Body)
		respBody := resp.Body
		if isGzip(resp.Header) {
			respBody, err = gzip.NewReader(respBody)
			if err != nil {
				return errors.Wrap(err, "can't unzip response body")
			}
		}
		respJson, err = ioutil.ReadAll(respBody)
		if err != nil {
			return errors.Wrap(err, "can't read response body")
		}
		return nil
	}, svc.retryTimes)
	return respJson, err
}

func isGzip(header http.Header) bool {
	return strings.Contains(header.Get("content-Encoding"), "gzip")
}
