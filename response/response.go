package response

import "github.com/cccaaannn/gohst/constant"

type Response struct {
	Body       string
	Headers    map[string]string
	StatusCode constant.HTTPStatusCode
}

func CreateOkResponse() *Response {
	return &Response{
		Headers:    make(map[string]string),
		Body:       "",
		StatusCode: constant.OK,
	}
}
