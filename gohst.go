package gohst

import (
	"github.com/cccaaannn/gohst/src/request"
	"github.com/cccaaannn/gohst/src/response"
	"github.com/cccaaannn/gohst/src/server"
)

type Request = request.Request
type Response = response.Response
type HandlerFunc = server.HandlerFunc
type Server = server.Server

func CreateServer() *Server {
	return server.CreateServer()
}
