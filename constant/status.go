package constant

import "strconv"

type HTTPStatusCode int

const (
	OK       HTTPStatusCode = 200
	NotFound HTTPStatusCode = 404
)

func (c HTTPStatusCode) String() string {
	return strconv.Itoa(int(c))
}

func (c HTTPStatusCode) Verb() string {
	switch c {
	case OK:
		return "OK"
	case NotFound:
		return "Not Found"
	default:
		return "Unknown"
	}
}
