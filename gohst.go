package gohst

import (
	"fmt"
	"net"
	"sync"

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
		headers:  getDefaultHeaders(),
	}
}

func (server *Server) AddHandler(requestPattern string, handlerFunc handlerFunc) {
	pathText, method, ok := util.ParseRequestPattern(requestPattern)
	if !ok {
		panic(fmt.Sprintf("Cannot add handler with request pattern of %s\n", requestPattern))
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

// func (server *Server) ListenAndServe(address string) {
// 	listener, err := net.Listen("tcp", address)
// 	if err != nil {
// 		fmt.Println("Error listening ", err.Error())
// 		return
// 	}
// 	defer listener.Close()
// 	fmt.Printf("Listening on %s\n", address)

// 	for {
// 		conn, err := listener.Accept()
// 		if err != nil {
// 			fmt.Println("Error accepting ", err.Error())
// 			return
// 		}
// 		go server.handleConnection(conn)
// 	}
// }

func (server *Server) ListenAndServe(address string) (func(), error) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("error listening: %v", err)
	}
	fmt.Printf("Listening on %s\n", address)

	var wg sync.WaitGroup
	stop := make(chan struct{})

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				select {
				case <-stop:
					return
				default:
					fmt.Println("Error accepting: ", err.Error())
					continue
				}
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
				server.handleConnection(conn)
			}()
		}
	}()

	shutdown := func() {
		close(stop)
		listener.Close()
		wg.Wait()
		fmt.Println("Server stopped")
	}

	return shutdown, nil
}
