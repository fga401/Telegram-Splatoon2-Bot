package nintendo
import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"github.com/pkg/errors"
	"io"
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

func randBytes(n int) ([]byte,error){
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


