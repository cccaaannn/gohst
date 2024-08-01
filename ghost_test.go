package gohst

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

const (
	ServerHost = "http://localhost"
	ServerPort = "8080"
)

type result struct {
	Query  map[string]string `json:"query"`
	Params map[string]string `json:"params"`
	Body   string            `json:"body"`
}

const apiResponse = `{"message": "Hello, World!"}`
const apiPageContent = `{"message": "Hello, World!"}`
const aboutPageContent = "<body><h1>About</h1><p>This is the about page</p></body>"
const notFoundPageContent = "<body>	<h1>Not Found</h1><p>The page you are looking for does not exist</p></body>"

func createAPIServer() *Server {
	// Create a new server
	server := CreateServer()

	// Set default headers
	headers := map[string]string{"Content-Type": "application/json"}
	server.SetHeaders(headers)

	server.AddHandler("POST /path1/:param1/path2/:param2", func(req *Request, res *Response) {
		result := result{
			Params: req.Params,
			Query:  req.Query,
			Body:   req.Body,
		}
		jsonResponse, _ := json.Marshal(result)
		res.Body = string(jsonResponse)
	})

	return server
}

func createHTMLServer() *Server {
	// Create a new server
	server := CreateServer()

	// Set default headers
	headers := map[string]string{"Content-Type": "text/html"}
	server.SetHeaders(headers)

	server.AddHandler("GET /about", func(req *Request, res *Response) {
		res.Body = aboutPageContent
	})

	server.AddHandler("GET /api", func(req *Request, res *Response) {
		res.Headers["Content-Type"] = "application/json"
		res.Body = apiPageContent
	})

	server.AddHandler("/*", func(req *Request, res *Response) {
		res.StatusCode = http.StatusNotFound
		res.Body = notFoundPageContent
	})

	return server
}

func createWithBrokenPattern() {
	server := CreateServer()
	server.AddHandler("POST invalid pattern", func(req *Request, res *Response) {})
}

func TestRequestParsing(t *testing.T) {
	// Given
	server := createAPIServer()
	shutdown, err := server.ListenAndServe(fmt.Sprintf(":%s", ServerPort))
	if err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	defer shutdown()
	time.Sleep(1 * time.Second) // Delay to allow the server to start

	// When
	resp, err := http.Post(
		fmt.Sprintf(
			"%s:%s/path1/test1/path2/5?query1=test2&query2=2",
			ServerHost,
			ServerPort,
		),
		"application/json",
		bytes.NewReader([]byte(apiResponse)),
	)
	if err != nil {
		t.Fatalf("Failed to send GET request: %v", err)
	}
	defer resp.Body.Close()

	// Then
	expectedStatusCode := http.StatusOK
	if resp.StatusCode != expectedStatusCode {
		t.Fatalf("Expected status code %v, got %v", expectedStatusCode, resp.StatusCode)
	}

	expectedContentType := "application/json"
	if resp.Header.Get("Content-Type") != expectedContentType {
		t.Fatalf("Expected content type %v, got %v", expectedContentType, resp.Header.Get("Content-Type"))
	}

	expectedResult := result{
		Params: map[string]string{"param1": "test1", "param2": "5"},
		Query:  map[string]string{"query1": "test2", "query2": "2"},
		Body:   apiResponse,
	}

	expectedJson, _ := json.Marshal(expectedResult)
	body, _ := io.ReadAll(resp.Body)
	if string(expectedJson) != string(body) {
		t.Fatalf("Expected response body %v, got %v", string(expectedJson), string(body))
	}
}

func TestGenericNotFound(t *testing.T) {
	// Given
	server := createAPIServer()
	shutdown, err := server.ListenAndServe(fmt.Sprintf(":%s", ServerPort))
	if err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	defer shutdown()
	time.Sleep(1 * time.Second) // Delay to allow the server to start

	// When
	resp, err := http.Get(
		fmt.Sprintf(
			"%s:%s/melon/5",
			ServerHost,
			ServerPort,
		),
	)
	if err != nil {
		t.Fatalf("Failed to send GET request: %v", err)
	}
	defer resp.Body.Close()

	// Then
	expectedStatusCode := http.StatusNotFound
	if resp.StatusCode != expectedStatusCode {
		t.Fatalf("Expected status code %v, got %v", expectedStatusCode, resp.StatusCode)
	}

	expectedContentType := "application/json"
	if resp.Header.Get("Content-Type") != expectedContentType {
		t.Fatalf("Expected content type %v, got %v", expectedContentType, resp.Header.Get("Content-Type"))
	}
}

func TestHtmlContent(t *testing.T) {
	// Given
	server := createHTMLServer()
	shutdown, err := server.ListenAndServe(fmt.Sprintf(":%s", ServerPort))
	if err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	defer shutdown()
	time.Sleep(1 * time.Second) // Delay to allow the server to start

	// When
	resp, err := http.Get(
		fmt.Sprintf(
			"%s:%s/about",
			ServerHost,
			ServerPort,
		),
	)
	if err != nil {
		t.Fatalf("Failed to send GET request: %v", err)
	}
	defer resp.Body.Close()

	// Then
	expectedStatusCode := http.StatusOK
	if resp.StatusCode != expectedStatusCode {
		t.Fatalf("Expected status code %v, got %v", expectedStatusCode, resp.StatusCode)
	}

	expectedContentType := "text/html"
	if resp.Header.Get("Content-Type") != expectedContentType {
		t.Fatalf("Expected content type %v, got %v", expectedContentType, resp.Header.Get("Content-Type"))
	}

	body, _ := io.ReadAll(resp.Body)
	if string(aboutPageContent) != string(body) {
		t.Fatalf("Expected response body %v, got %v", string(aboutPageContent), string(body))
	}
}

func TestCustomRequestHeader(t *testing.T) {
	// Given
	server := createHTMLServer()
	shutdown, err := server.ListenAndServe(fmt.Sprintf(":%s", ServerPort))
	if err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	defer shutdown()
	time.Sleep(1 * time.Second) // Delay to allow the server to start

	// When
	resp, err := http.Get(
		fmt.Sprintf(
			"%s:%s/api",
			ServerHost,
			ServerPort,
		),
	)
	if err != nil {
		t.Fatalf("Failed to send GET request: %v", err)
	}
	defer resp.Body.Close()

	// Then
	expectedStatusCode := http.StatusOK
	if resp.StatusCode != expectedStatusCode {
		t.Fatalf("Expected status code %v, got %v", expectedStatusCode, resp.StatusCode)
	}

	expectedContentType := "application/json"
	if resp.Header.Get("Content-Type") != expectedContentType {
		t.Fatalf("Expected content type %v, got %v", expectedContentType, resp.Header.Get("Content-Type"))
	}

	body, _ := io.ReadAll(resp.Body)
	if string(apiPageContent) != string(body) {
		t.Fatalf("Expected response body %v, got %v", string(aboutPageContent), string(body))
	}
}

func TestWildcardPath(t *testing.T) {
	// Given
	server := createHTMLServer()
	shutdown, err := server.ListenAndServe(fmt.Sprintf(":%s", ServerPort))
	if err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	defer shutdown()
	time.Sleep(1 * time.Second) // Delay to allow the server to start

	// When
	resp, err := http.Get(
		fmt.Sprintf(
			"%s:%s/bananaaaaa",
			ServerHost,
			ServerPort,
		),
	)
	if err != nil {
		t.Fatalf("Failed to send GET request: %v", err)
	}
	defer resp.Body.Close()

	// Then
	expectedStatusCode := http.StatusNotFound
	if resp.StatusCode != expectedStatusCode {
		t.Fatalf("Expected status code %v, got %v", expectedStatusCode, resp.StatusCode)
	}

	expectedContentType := "text/html"
	if resp.Header.Get("Content-Type") != expectedContentType {
		t.Fatalf("Expected content type %v, got %v", expectedContentType, resp.Header.Get("Content-Type"))
	}

	body, _ := io.ReadAll(resp.Body)
	if string(notFoundPageContent) != string(body) {
		t.Fatalf("Expected response body %v, got %v", string(notFoundPageContent), string(body))
	}
}

func TestBrokenHandlerPattern(t *testing.T) {
	// Then
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("Expected panic, but code did not panic")
		}
	}()

	// Given
	createWithBrokenPattern()
}
