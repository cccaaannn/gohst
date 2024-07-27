package url

import "strings"

func ParseQuery(query string) map[string]string {
	params := make(map[string]string)
	parts := strings.Split(query, "&")
	for _, part := range parts {
		split := strings.Split(part, "=")
		if len(split) == 2 {
			params[split[0]] = split[1]
		} else {
			params[split[0]] = ""
		}
	}
	return params
}

func SplitQuery(path string) (string, string) {
	parts := strings.Split(path, "?")
	if len(parts) == 1 {
		return parts[0], ""
	}
	return parts[0], parts[1]
}
