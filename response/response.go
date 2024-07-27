package response

import "github.com/cccaaannn/gohst/constant"

type Response struct {
	Body       string
	Headers    map[string]string
	StatusCode constant.HTTPStatusCode
}
