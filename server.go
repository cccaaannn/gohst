package gohst

import (
	"fmt"
	"net"

	"github.com/cccaaannn/gohst/src/constant"
	"github.com/cccaaannn/gohst/src/request"
	"github.com/cccaaannn/gohst/src/response"
	"github.com/cccaaannn/gohst/src/url"
	"github.com/cccaaannn/gohst/src/util"
)

type handlerFunc func(*Request, *Response)

type handler struct {
	path        url.Path
	method      string
	handlerFunc handlerFunc
}

type Server struct {
	handlers []handler
	headers  map[string]string
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

func (server *Server) getMergedHeaders(response *Response) map[string]string {
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

func (server *Server) buildResponseString(response *Response) string {
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

	response := response.CreateOkResponse()

	if !matched {
		response.StatusCode = constant.NotFoundStatus
		responseStr := server.buildResponseString(response)
		conn.Write([]byte(responseStr))
		return
	}

	handler.handlerFunc(req, response)

	responseStr := server.buildResponseString(response)

	conn.Write([]byte(responseStr))
}
