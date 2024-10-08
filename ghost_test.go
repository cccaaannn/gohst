package gohst

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

const (
	ServerHost          = "http://localhost"
	ServerHostTLS       = "https://localhost"
	ServerPort          = "8080"
	TestCertPath        = "test/cert/localhost.crt"
	TestKeyPath         = "test/cert/localhost.key"
	ApiResponse         = `{"message": "Hello, World!"}`
	ApiPageContent      = `{"message": "Hello, World!"}`
	AboutPageContent    = "<body><h1>About</h1><p>This is the about page</p></body>"
	NotFoundPageContent = "<body><h1>Not Found</h1><p>The page you are looking for does not exist</p></body>"
	UnauthorizedContent = "<body><h1>401 Unauthorized</h1></body>"
	TestHeaderContent1  = "banana"
	TestHeaderContent2  = "melon"
	TestHeaderContent3  = "apple"
	TestHeaderName      = "Test-Header"
)

type result struct {
	Query  map[string]string `json:"query"`
	Params map[string]string `json:"params"`
	Body   string            `json:"body"`
}

func createAPIServer() *Server {
	// Create a new server
	server := CreateServer()

	// Set default headers
	headers := map[string]string{"Content-Type": "application/json"}
	server.SetHeaders(headers)

	server.AddHandler("GET /tls", func(req *Request, res *Response) {
		res.Body = ApiResponse
	})

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
		res.Body = AboutPageContent
	})

	server.AddHandler("GET /api", func(req *Request, res *Response) {
		res.Headers["Content-Type"] = "application/json"
		res.Body = ApiPageContent
	})

	server.AddHandler("/*", func(req *Request, res *Response) {
		res.StatusCode = http.StatusNotFound
		res.Body = NotFoundPageContent
	})

	return server
}

func createBrokenPatternServer() {
	server := CreateServer()
	server.AddHandler("POST invalid pattern", func(req *Request, res *Response) {})
}

func createMiddlewareServer() *Server {
	// Create a new server
	server := CreateServer()

	// Middleware 1
	server.Use(func(next HandlerFunc) HandlerFunc {
		return func(req *Request, res *Response) {

			var testHeader string = req.Headers[TestHeaderName]

			if testHeader == "" {
				res.StatusCode = 401
				res.Body = UnauthorizedContent
				return
			}

			req.Context[TestHeaderName] = testHeader

			next(req, res)
		}
	})

	// Middleware 2
	server.Use(func(next HandlerFunc) HandlerFunc {
		return func(req *Request, res *Response) {
			req.Context[TestHeaderName] = fmt.Sprintf("%s%s", req.Context[TestHeaderName], TestHeaderContent2)
			next(req, res)
		}
	})

	server.AddHandler("GET /middleware", func(req *Request, res *Response) {
		body := fmt.Sprintf("%s%s", req.Context[TestHeaderName], TestHeaderContent3)
		res.Body = body
	})

	return server
}

func setup() {
	time.Sleep(500 * time.Millisecond)
}

func TestRequestParsing(t *testing.T) {
	// Given
	setup()
	server := createAPIServer()
	stop, err := server.ListenAndServe(fmt.Sprintf(":%s", ServerPort))
	if err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	defer close(stop)
	time.Sleep(500 * time.Millisecond) // Delay to allow the server to start

	// When
	resp, err := http.Post(
		fmt.Sprintf(
			"%s:%s/path1/test1/path2/5?query1=test2&query2=2",
			ServerHost,
			ServerPort,
		),
		"application/json",
		bytes.NewReader([]byte(ApiResponse)),
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
		Body:   ApiResponse,
	}

	expectedJson, _ := json.Marshal(expectedResult)
	body, _ := io.ReadAll(resp.Body)
	if string(expectedJson) != string(body) {
		t.Fatalf("Expected response body %v, got %v", string(expectedJson), string(body))
	}
}

func TestGenericNotFound(t *testing.T) {
	// Given
	setup()
	server := createAPIServer()
	stop, err := server.ListenAndServe(fmt.Sprintf(":%s", ServerPort))
	if err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	defer close(stop)
	time.Sleep(500 * time.Millisecond) // Delay to allow the server to start

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
	setup()
	server := createHTMLServer()
	stop, err := server.ListenAndServe(fmt.Sprintf(":%s", ServerPort))
	if err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	defer close(stop)
	time.Sleep(500 * time.Millisecond) // Delay to allow the server to start

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
	if AboutPageContent != string(body) {
		t.Fatalf("Expected response body %v, got %v", AboutPageContent, string(body))
	}
}

func TestCustomRequestHeader(t *testing.T) {
	// Given
	setup()
	server := createHTMLServer()
	stop, err := server.ListenAndServe(fmt.Sprintf(":%s", ServerPort))
	if err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	defer close(stop)
	time.Sleep(500 * time.Millisecond) // Delay to allow the server to start

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
	if ApiPageContent != string(body) {
		t.Fatalf("Expected response body %v, got %v", ApiPageContent, string(body))
	}
}

func TestWildcardPath(t *testing.T) {
	// Given
	setup()
	server := createHTMLServer()
	stop, err := server.ListenAndServe(fmt.Sprintf(":%s", ServerPort))
	if err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	defer close(stop)
	time.Sleep(500 * time.Millisecond) // Delay to allow the server to start

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
	if string(NotFoundPageContent) != string(body) {
		t.Fatalf("Expected response body %v, got %v", string(NotFoundPageContent), string(body))
	}
}

func TestBrokenHandlerPattern(t *testing.T) {
	setup()

	// Then
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("Expected panic, but code did not panic")
		}
	}()

	// Given
	createBrokenPatternServer()
}

func TestTLS(t *testing.T) {
	// Given
	setup()
	server := createAPIServer()
	stop, err := server.ListenAndServeTLS(fmt.Sprintf(":%s", ServerPort), TestCertPath, TestKeyPath)
	if err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	defer close(stop)
	time.Sleep(500 * time.Millisecond) // Delay to allow the server to start

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // Test only
			},
		},
	}

	// When
	resp, err := client.Get(
		fmt.Sprintf(
			"%s:%s/tls",
			ServerHostTLS,
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

	body, _ := io.ReadAll(resp.Body)
	expectedBody := ApiResponse
	bodyStr := string(body)
	if bodyStr != expectedBody {
		t.Fatalf("Expected response body %v, got %v", string(expectedBody), bodyStr)
	}
}

func TestMiddlewareInterception(t *testing.T) {
	// Given
	setup()
	server := createMiddlewareServer()
	stop, err := server.ListenAndServe(fmt.Sprintf(":%s", ServerPort))
	if err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	defer close(stop)
	time.Sleep(500 * time.Millisecond) // Delay to allow the server to start

	// When
	resp, err := http.Get(
		fmt.Sprintf(
			"%s:%s/middleware",
			ServerHost,
			ServerPort,
		),
	)
	if err != nil {
		t.Fatalf("Failed to send GET request: %v", err)
	}
	defer resp.Body.Close()

	// Then
	expectedStatusCode := http.StatusUnauthorized
	if resp.StatusCode != expectedStatusCode {
		t.Fatalf("Expected status code %v, got %v", expectedStatusCode, resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if UnauthorizedContent != string(body) {
		t.Fatalf("Expected response body %v, got %v", UnauthorizedContent, string(body))
	}
}

func TestMiddlewareRunOrder(t *testing.T) {
	// Server middleware should run in the order they are added
	// Expected response body: bananamelonapple

	// Given
	setup()
	server := createMiddlewareServer()
	stop, err := server.ListenAndServe(fmt.Sprintf(":%s", ServerPort))
	if err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	defer close(stop)
	time.Sleep(500 * time.Millisecond) // Delay to allow the server to start

	// When
	req, err := http.NewRequest("GET",
		fmt.Sprintf(
			"%s:%s/middleware",
			ServerHost,
			ServerPort,
		),
		nil,
	)
	if err != nil {
		t.Fatalf("Failed to create GET request: %v", err)
	}

	req.Header.Add(TestHeaderName, TestHeaderContent1)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to send GET request: %v", err)
	}
	defer resp.Body.Close()

	// Then
	expectedStatusCode := http.StatusOK
	if resp.StatusCode != expectedStatusCode {
		t.Fatalf("Expected status code %v, got %v", expectedStatusCode, resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	expectedBody := fmt.Sprintf("%s%s%s", TestHeaderContent1, TestHeaderContent2, TestHeaderContent3)
	if expectedBody != string(body) {
		t.Fatalf("Expected response body %v, got %v", expectedBody, string(body))
	}
}
