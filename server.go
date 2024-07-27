package gohst

import (
	"fmt"
	"net"

	"github.com/cccaaannn/gohst/constant"
	"github.com/cccaaannn/gohst/request"
	"github.com/cccaaannn/gohst/response"
	"github.com/cccaaannn/gohst/url"
	"github.com/cccaaannn/gohst/util"
)

type HandlerFunc func(*request.Request, *response.Response)

type handler struct {
	path        url.Path
	method      constant.HttpMethod
	handlerFunc HandlerFunc
}

type Server struct {
	handlers []handler
	headers  map[string]string
}

func (server *Server) getResponseContentTypeHeader(responseHeaders map[string]string) string {
	if _, ok := responseHeaders[constant.ContentTypeHeader.String()]; ok {
		return ""
	}
	if val, ok := server.headers[constant.ContentTypeHeader.String()]; ok {
		return fmt.Sprintf("%s: %s\r\n", constant.ContentTypeHeader.String(), val)
	}
	return fmt.Sprintf("%s: %s\r\n", constant.ContentTypeHeader.String(), constant.TextHtml.String())
}

func (server *Server) getContentLengthHeader(body string) string {
	contentLength := []byte(body)
	if len(contentLength) == 0 {
		return ""
	}
	return fmt.Sprintf("%s: %d\r\n", constant.ContentLength, len(contentLength))
}

func (server *Server) getBodyString(body string) string {
	if body != "" {
		return fmt.Sprintf("\r\n%s\r\n", body)
	}
	return ""
}

func (server *Server) matchHandler(path string, method string) (handler, map[string]string, bool) {
	for _, handler := range server.handlers {
		params, ok := handler.path.Match(path)
		if ok && (handler.method == "" || handler.method == constant.HttpMethod(method)) {
			return handler, params, true
		}
	}
	return handler{}, nil, false
}

func (server *Server) handleGenericNotFound() string {
	responseStr := "" +
		"%s %s %s\r\n" +
		"Date: %s\r\n" +
		"Server: %s\r\n" +
		"Connection: close\r\n" +
		"%s"

	return fmt.Sprintf(
		responseStr,
		constant.HTTPVersion, constant.NotFound.String(), constant.NotFound.Verb(),
		util.GetHttpTime(),
		constant.ServerName,
		server.getResponseContentTypeHeader(nil),
	)
}

func (server *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	req, err := request.ParseRequest(conn)
	if err != nil {
		fmt.Println("Error parsing request:", err)
		return
	}

	responseStr := "" +
		"%s %s %s\r\n" +
		"Date: %s\r\n" +
		"Server: %s\r\n" +
		"Connection: close\r\n" +
		"%s" +
		"%s" +
		"%s" +
		"%s"

	// Query parsing
	path, query := url.SplitQuery(req.Path)
	req.Query = url.ParseQuery(query)

	// Path parsing
	handler, params, ok := server.matchHandler(path, req.Method)
	req.Params = params

	if !ok {
		responseStr = server.handleGenericNotFound()
		conn.Write([]byte(responseStr))
		return
	}

	response := response.Response{
		Body:       "",
		Headers:    make(map[string]string),
		StatusCode: constant.OK,
	}

	handler.handlerFunc(req, &response)

	userHeaders := ""
	for key, val := range response.Headers {
		userHeaders += fmt.Sprintf("%s: %s\r\n", key, val)
	}

	responseStr = fmt.Sprintf(
		responseStr,
		constant.HTTPVersion, response.StatusCode.String(), response.StatusCode.Verb(),
		util.GetHttpTime(),
		constant.ServerName,
		server.getResponseContentTypeHeader(response.Headers),
		server.getContentLengthHeader(response.Body),
		userHeaders,
		server.getBodyString(response.Body),
	)

	conn.Write([]byte(responseStr))
}
