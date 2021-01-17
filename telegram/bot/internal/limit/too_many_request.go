package limit

import (
	"regexp"
	"strconv"
)

var tooManyRequestRegExp = regexp.MustCompile(`Too Many Requests: retry after (?P<sec>\d+)`)

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

func IsTooManyRequestString(s string) (bool, int) {
	results := tooManyRequestRegExp.FindStringSubmatch(s)
	if len(results) > 0 {
		sec, _ := strconv.Atoi(results[1])
		return true, sec
	}
	return false, 0
}

