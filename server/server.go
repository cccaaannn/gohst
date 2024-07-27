package server

import (
	"fmt"
	"net"

	"github.com/cccaaannn/gohst/constant"
	"github.com/cccaaannn/gohst/request"
	"github.com/cccaaannn/gohst/response"
	"github.com/cccaaannn/gohst/url"
)

type HandlerFunc func(*request.Request, *response.Response)

type handler struct {
	path        url.Path
	method      constant.HttpMethod
	handlerFunc HandlerFunc
}

type server struct {
	handlers []handler
	headers  map[string]string
}

func CreateServer() *server {
	return &server{
		handlers: make([]handler, 0),
	}
}

func (server *server) AddHandler(requestPattern string, handlerFunc HandlerFunc) {
	pathText, method, ok := parseRequestPattern(requestPattern)
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

func (server *server) SetHeaders(headers map[string]string) {
	server.headers = headers
}

func (server *server) getResponseContentTypeHeader(responseHeaders map[string]string) string {
	if _, ok := responseHeaders[constant.ContentTypeHeader.String()]; ok {
		return ""
	}
	if val, ok := server.headers[constant.ContentTypeHeader.String()]; ok {
		return fmt.Sprintf("%s: %s", constant.ContentTypeHeader.String(), val)
	}
	return fmt.Sprintf("%s: %s", constant.ContentTypeHeader.String(), constant.TextHtml.String())
}

func (server *server) getContentLengthHeader(body string) string {
	contentLength := []byte(body)
	if len(contentLength) == 0 {
		return ""
	}
	return fmt.Sprintf("%s: %d", constant.ContentLength, len(contentLength))
}

func (server *server) matchHandler(path string, method string) (handler, map[string]string, bool) {
	for _, handler := range server.handlers {
		params, ok := handler.path.Match(path)
		if ok && (handler.method == "" || handler.method == constant.HttpMethod(method)) {
			return handler, params, true
		}
	}
	return handler{}, nil, false
}

func (server *server) handleGenericNotFound(responseStr string) string {
	dateString := getHttpTime()
	return fmt.Sprintf(
		responseStr,
		constant.HTTPVersion, constant.NotFound.String(), constant.NotFound.Verb(),
		dateString,
		constant.ServerName,
		server.getResponseContentTypeHeader(nil),
		"",
		"",
		"",
	)
}

func (server *server) handleConnection(conn net.Conn) {
	defer conn.Close()

	req, err := request.ParseRequest(conn)
	if err != nil {
		fmt.Println("Error parsing request:", err)
		return
	}

	dateString := getHttpTime()

	responseStr := "" +
		"%s %s %s\r\n" +
		"Date: %s\r\n" +
		"Server: %s\r\n" +
		"Connection: close\r\n" +
		"%s\r\n" +
		"%s\r\n" +
		"%s" +
		"\r\n" +
		"%s\r\n"

	// Query parsing
	path, query := url.SplitQuery(req.Path)
	req.Query = url.ParseQuery(query)

	// Path parsing
	handler, params, ok := server.matchHandler(path, req.Method)
	req.Params = params

	if !ok {
		server.handleGenericNotFound(responseStr)
		conn.Write([]byte(responseStr))
		return
	}

	response := response.Response{
		Body:    "",
		Headers: make(map[string]string),
	}

	handler.handlerFunc(req, &response)

	userHeaders := ""
	for key, val := range response.Headers {
		userHeaders += fmt.Sprintf("%s: %s\r\n", key, val)
	}

	responseStr = fmt.Sprintf(
		responseStr,
		constant.HTTPVersion, constant.OK.String(), constant.OK.Verb(),
		dateString,
		constant.ServerName,
		server.getResponseContentTypeHeader(response.Headers),
		server.getContentLengthHeader(response.Body),
		userHeaders,
		response.Body,
	)

	conn.Write([]byte(responseStr))
}
