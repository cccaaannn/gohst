package server

import (
	"crypto/tls"
	"fmt"
	"net"
	"sync"

	"github.com/cccaaannn/gohst/src/constant"
	"github.com/cccaaannn/gohst/src/request"
	"github.com/cccaaannn/gohst/src/response"
	"github.com/cccaaannn/gohst/src/url"
	"github.com/cccaaannn/gohst/src/util"
)

type HandlerFunc func(*request.Request, *response.Response)

type handler struct {
	path        url.Path
	method      string
	handlerFunc HandlerFunc
}

type Server struct {
	handlers []handler
	headers  map[string]string
}

func CreateServer() *Server {
	return &Server{
		handlers: make([]handler, 0),
		headers:  getDefaultHeaders(),
	}
}

func (sv *Server) AddHandler(requestPattern string, handlerFunc HandlerFunc) {
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

	sv.handlers = append(sv.handlers, handler)
}

func (sv *Server) SetHeaders(headers map[string]string) {
	sv.headers = headers
}

func (sv *Server) ListenAndServeTLS(address string, certFile string, keyFile string) (chan struct{}, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("error loading certificate: %v", err)
	}

	config := &tls.Config{Certificates: []tls.Certificate{cert}}
	listener, err := tls.Listen("tcp", address, config)
	if err != nil {
		return nil, fmt.Errorf("error listening: %v", err)
	}
	fmt.Printf("Listening on %s\n", address)

	return sv.listenAndServe(listener)
}

func (sv *Server) ListenAndServe(address string) (chan struct{}, error) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("error listening: %v", err)
	}
	fmt.Printf("Listening on %s\n", address)

	return sv.listenAndServe(listener)
}

func (sv *Server) listenAndServe(listener net.Listener) (chan struct{}, error) {
	var once sync.Once
	var wg sync.WaitGroup
	stop := make(chan struct{})

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				// Non blocking select to check errors and stop signal
				select {
				// When listener is closed it errors, so we check if channel is closed
				case <-stop:
					return
				// If it's not closed we continue with logging the error
				default:
					fmt.Println("Error accepting: ", err.Error())
					continue
				}
			}

			// Handle new connection in a separate goroutine
			wg.Add(1)
			go func() {
				defer wg.Done()
				sv.handleConnection(conn)
			}()
		}
	}()

	// Goroutine to handle gracefully shuting down the server
	go func() {
		<-stop
		once.Do(func() {
			listener.Close()
			wg.Wait()
			fmt.Println("Server stopped")
		})
	}()

	return stop, nil
}

func getDefaultHeaders() map[string]string {
	return map[string]string{
		constant.ServerHeader.String():      constant.ServerName,
		constant.ConnectionHeader.String():  "close",
		constant.ContentTypeHeader.String(): constant.TextHtml.String(),
	}
}

func (server *Server) matchHandler(path string, method string) (handler, map[string]string, bool) {
	for _, handler := range server.handlers {
		params, ok := handler.path.Match(path)
		if ok && (handler.method == "" || handler.method == method) {
			return handler, params, true
		}
	}
	return handler{}, nil, false
}

func (server *Server) getMergedHeaders(response *response.Response) map[string]string {
	contentLength := len([]byte(response.Body))

	requestHeaders := map[string]string{
		constant.DateHeader.String():          util.GetHttpTime(),
		constant.ContentLengthHeader.String(): fmt.Sprintf("%d", contentLength),
	}

	mergedHeaders := make(map[string]string)

	// Add default headers
	for key, val := range server.headers {
		mergedHeaders[key] = val
	}

	// Request headers override default headers
	for key, val := range requestHeaders {
		mergedHeaders[key] = val
	}

	// User headers override request headers
	for key, val := range response.Headers {
		mergedHeaders[key] = val
	}

	return mergedHeaders
}

func (server *Server) buildResponseString(response *response.Response) string {
	responseStr := "" +
		"%s\r\n" +
		"%s" +
		"\r\n%s"

	requestLineStr := fmt.Sprintf(
		"%s %s %s",
		constant.HTTPVersion, response.StatusCode.String(), response.StatusCode.Verb(),
	)

	mergedHeaders := server.getMergedHeaders(response)
	mergedHeadersStr := ""
	for key, val := range mergedHeaders {
		mergedHeadersStr += fmt.Sprintf("%s: %s\r\n", key, val)
	}

	responseStr = fmt.Sprintf(
		responseStr,
		requestLineStr,
		mergedHeadersStr,
		response.Body,
	)

	return responseStr
}

func (server *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	req, err := request.ParseRequest(conn)
	if err != nil {
		fmt.Println("Error parsing request:", err)
		return
	}

	// Query parsing
	path, query := url.SplitQuery(req.Path)
	req.Query = url.ParseQuery(query)

	// Path parsing
	handler, params, matched := server.matchHandler(path, req.Method)
	req.Params = params

	res := response.CreateOkResponse()

	if !matched {
		res.StatusCode = constant.NotFoundStatus
		responseStr := server.buildResponseString(res)
		conn.Write([]byte(responseStr))
		return
	}

	handler.handlerFunc(req, res)

	responseStr := server.buildResponseString(res)

	conn.Write([]byte(responseStr))
}
