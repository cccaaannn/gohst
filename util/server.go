package util

import (
	"strings"
	"time"
)

func GetHttpTime() string {
	now := time.Now()
	dateString := now.UTC().Format(time.RFC1123)
	dateString = strings.ReplaceAll(dateString, "UTC", "GMT")
	return dateString
}

func ParseRequestPattern(requestPattern string) (string, string, bool) {
	splitRequestPattern := strings.Split(requestPattern, " ")
	pathText := ""
	var method = ""
	if len(splitRequestPattern) == 2 {
		method = splitRequestPattern[0]
		pathText = splitRequestPattern[1]
	} else if len(splitRequestPattern) == 1 {
		pathText = splitRequestPattern[0]
	} else {
		return "", "", false
	}

	return pathText, method, true
}
