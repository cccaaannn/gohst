package constant

import "strconv"

type HTTPStatusCode int

const (
	OK                  HTTPStatusCode = 200
	Created             HTTPStatusCode = 201
	BadRequest          HTTPStatusCode = 400
	Unauthorized        HTTPStatusCode = 401
	NotFound            HTTPStatusCode = 404
	InternalServerError HTTPStatusCode = 500
)

func (c HTTPStatusCode) String() string {
	return strconv.Itoa(int(c))
}

func (c HTTPStatusCode) Verb() string {
	switch c {
	case OK:
		return "OK"
	case Created:
		return "Created"
	case BadRequest:
		return "Bad Request"
	case Unauthorized:
		return "Unauthorized"
	case NotFound:
		return "Not Found"
	case InternalServerError:
		return "Internal Server Error"
	default:
		return "Unknown Status Code"
	}
}
