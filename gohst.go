package gohst

import (
	"fmt"
	"net"

	"github.com/cccaaannn/gohst/request"
	"github.com/cccaaannn/gohst/response"
	"github.com/cccaaannn/gohst/url"
	"github.com/cccaaannn/gohst/util"
)

type Request = request.Request
type Response = response.Response

func CreateServer() *Server {
	return &Server{
		handlers: make([]handler, 0),
	}
}

func (server *Server) AddHandler(requestPattern string, handlerFunc HandlerFunc) {
	pathText, method, ok := util.ParseRequestPattern(requestPattern)
	if !ok {
		fmt.Printf("Cannot add handler with request pattern of %s\n", requestPattern)
		return
	}

	path := url.CreatePath(pathText)
	handler := handler{
		path:        path,
		method:      method,
		handlerFunc: handlerFunc,
	}

	server.handlers = append(server.handlers, handler)
}

func (server *Server) SetHeaders(headers map[string]string) {
	server.headers = headers
}

func (server *Server) ListenAndServe(address string) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("Error listening ", err.Error())
		return
	}
	defer listener.Close()
	fmt.Printf("Listening on %s\n", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting ", err.Error())
			return
		}
		go server.handleConnection(conn)
	}
}
