package callbackquery

import "strings"

func SetPrefix(prefix, text string) string {
	return prefix + ":" + text
}

func GetPrefix(data string) string {
	index := strings.Index(data, ":")
	if index == -1 {
		return ""
	}
	return data[:index]
}

func GetText(data string) string {
	index := strings.Index(data, ":")
	if index == -1 {
		return data
	}
	return data[index+1:]
}
