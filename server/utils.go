package server

import (
	"strings"
	"time"

	"github.com/cccaaannn/gohst/constant"
)

func getHttpTime() string {
	now := time.Now()
	dateString := now.UTC().Format(time.RFC1123)
	dateString = strings.ReplaceAll(dateString, "UTC", "GMT")
	return dateString
}

func parseRequestPattern(requestPattern string) (string, constant.HttpMethod, bool) {
	splitRequestPattern := strings.Split(requestPattern, " ")
	pathText := ""
	var method constant.HttpMethod = ""
	if len(splitRequestPattern) == 2 {
		method = constant.HttpMethod(splitRequestPattern[0])
		pathText = splitRequestPattern[1]
	} else if len(splitRequestPattern) == 1 {
		pathText = splitRequestPattern[0]
	} else {
		return "", "", false
	}

	return pathText, method, true
}
