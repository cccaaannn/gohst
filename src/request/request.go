package request

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

type Request struct {
	Method   string
	Path     string
	Protocol string
	Body     string
	Query    map[string]string
	Params   map[string]string
	Headers  map[string]string
	Context  map[string]any
}

func readUntilBody(reader *bufio.Reader) (string, error) {
	var headers strings.Builder

	for {
		header, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}

			fmt.Println("Error reading:", err.Error())
			return "", err
		}

		headers.WriteString(header)

		if header == "\r\n" {
			break
		}
	}

	return headers.String(), nil
}

func readBody(reader *bufio.Reader, headers string) (string, error) {
	var bodyBuffer strings.Builder
	contentLengthHeaderStr := "Content-Length: "

	if strings.Contains(headers, contentLengthHeaderStr) {
		// Extract the content length
		start := strings.Index(headers, contentLengthHeaderStr) + len(contentLengthHeaderStr)
		end := strings.Index(headers[start:], "\r\n")
		contentLength, err := strconv.Atoi(headers[start : start+end])
		if err != nil {
			fmt.Println("Error parsing Content-Length:", err.Error())
			return "", err
		}

		// Read the body based on the content length
		body := make([]byte, contentLength)
		_, err = io.ReadFull(reader, body)
		if err != nil {
			fmt.Println("Error reading body:", err.Error())
			return "", err
		}
		bodyBuffer.Write(body)
	}

	return bodyBuffer.String(), nil
}

func parseRequestLine(requestLine string) (string, string, string, error) {
	requestParts := strings.Fields(requestLine)

	if len(requestParts) < 3 {
		return "", "", "", fmt.Errorf("invalid request line")
	}
	return requestParts[0], requestParts[1], requestParts[2], nil
}

func parseHeaders(headers string) map[string]string {
	headerMap := make(map[string]string)

	headerLines := strings.Split(headers, "\r\n")
	for _, headerLine := range headerLines {
		if headerLine == "" {
			continue
		}

		headerParts := strings.Split(headerLine, ": ")
		if len(headerParts) != 2 {
			continue
		}
		headerMap[headerParts[0]] = headerParts[1]
	}

	return headerMap
}

func ParseRequest(conn net.Conn) (*Request, error) {
	reader := bufio.NewReader(conn)

	headers, err := readUntilBody(reader)
	if err != nil {
		return nil, err
	}

	method, path, protocol, err := parseRequestLine(headers)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%s %s %s\n", method, path, protocol)

	headerMap := parseHeaders(headers)
	body, err := readBody(reader, headers)
	if err != nil {
		return nil, err
	}

	req := &Request{
		Method:   method,
		Path:     path,
		Protocol: protocol,
		Body:     body,
		Headers:  headerMap,
		Context:  make(map[string]any),
	}

	return req, nil
}
