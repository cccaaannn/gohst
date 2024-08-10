package server

import (
	"github.com/cccaaannn/gohst/src/request"
	"github.com/cccaaannn/gohst/src/response"
	"github.com/cccaaannn/gohst/src/url"
)

type HandlerFunc func(*request.Request, *response.Response)

type handler struct {
	path        url.Path
	method      string
	handlerFunc HandlerFunc
}
