package limit

import (
	"regexp"
	"strconv"
)

var tooManyRequestRegExp = regexp.MustCompile(`Too Many Requests: retry after (?P<sec>\d+)`)

// IsTooManyRequestError returns true and the recommended waiting time in second if an error is caused by sending too many request to telegram.
func IsTooManyRequestError(e error) (bool, int) {
	if e == nil {
		return false, 0
	}
	results := tooManyRequestRegExp.FindStringSubmatch(e.Error())
	if len(results) > 0 {
		sec, _ := strconv.Atoi(results[1])
		return true, sec
	}
	return false, 0
}

// IsTooManyRequestString returns true and the recommended waiting time in second if the string is the error string of TooManyRequestError.
func IsTooManyRequestString(s string) (bool, int) {
	results := tooManyRequestRegExp.FindStringSubmatch(s)
	if len(results) > 0 {
		sec, _ := strconv.Atoi(results[1])
		return true, sec
	}
	return false, 0
}
