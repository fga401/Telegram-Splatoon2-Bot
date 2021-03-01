package callbackquery

import "strings"

// SetPrefix returns a string combine the prefix and text.
func SetPrefix(prefix, text string) string {
	return prefix + ":" + text
}

// GetPrefix extracts the prefix from the combined string.
func GetPrefix(data string) string {
	index := strings.Index(data, ":")
	if index == -1 {
		return ""
	}
	return data[:index]
}

// GetText extracts text from the combined string.
func GetText(data string) string {
	index := strings.Index(data, ":")
	if index == -1 {
		return data
	}
	return data[index+1:]
}
