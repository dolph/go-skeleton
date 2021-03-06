package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	// Google App Engine
	"google.golang.org/appengine/aetest"

	// Request routing
	"github.com/gorilla/mux"
)

func NewServer(t *testing.T) TestHandler {
	instance, err := aetest.NewInstance(nil)
	if err != nil {
		t.Fatalf("Failed to create instance: %v", err)
	}
	defer instance.Close()

	return TestHandler{t, Router(), instance}
}

type TestHandler struct {
	t *testing.T

	// The HTTP router to be tested.
	router *mux.Router

	// The HTTP server to route requests through.
	instance aetest.Instance
}

// Build an HTTP request, pass it to the HTTP handler, and return the response.
func (handler TestHandler) request(method, path string, headers map[string]string) TestResponse {
	request, err := handler.instance.NewRequest(method, path, nil)
	if err != nil {
		handler.t.Fatalf("Failed to create request: %v", err)
	}

	// Set request headers, if any.
	for header, value := range headers {
		request.Header.Set(header, value)
	}

	// Set an arbitrary remote address for logging purposes, etc.
	request.RemoteAddr = "1.2.3.4:80"

	response := httptest.NewRecorder()
	handler.router.ServeHTTP(response, request)
	return TestResponse{handler.t, response}
}

// Make a GET request to the HTTP handler, and return the response.
func (handler TestHandler) Get(path string, headers map[string]string) TestResponse {
	return handler.request(http.MethodGet, path, headers)
}

// Make a POST request to the HTTP handler, and return the response.
func (handler TestHandler) Post(path string, headers map[string]string) TestResponse {
	return handler.request(http.MethodPost, path, headers)
}

type TestResponse struct {
	t *testing.T

	// The HTTP response being asserted against.
	r *httptest.ResponseRecorder
}

// Ensure that the response contains the expected status code.
func (response TestResponse) AssertStatusEquals(expected int) {
	if response.r.Code != expected {
		response.t.Errorf(
			"Handler returned unexpected status code: got `%v` want `%v`",
			response.r.Code, expected)
	}
}

// Ensure that the response body is exactly as expected.
func (response TestResponse) AssertBodyEquals(expected string) {
	if actual := response.r.Body.String(); actual != expected {
		response.t.Errorf(
			"Handler returned unexpected body: got `%v` want `%v`",
			actual, expected)
	}
}

// Ensure that the response body contains a substring.
func (response TestResponse) AssertBodyContains(substr string) {
	if actual := response.r.Body.String(); !strings.Contains(actual, substr) {
		response.t.Errorf(
			"Handler returned unexpected body: did not find `%v` in `%v`",
			substr, actual)
	}
}

// Ensure that the response contains a specific header.
func (response TestResponse) AssertHeaderExists(header string) {
	if _, ok := response.r.Header()[header]; !ok {
		response.t.Errorf(
			"Handler did not set header `%v`",
			header)
	}
}

// Ensure that the response contains a specific header-value pair.
func (response TestResponse) AssertHeaderContains(header, expected string) {
	response.AssertHeaderExists(header)
	actuals, _ := response.r.Header()[header]
	for _, actual := range actuals {
		if actual == expected {
			return
		}
	}

	response.t.Errorf(
		"Handler returned unexpected %v: got `%v` want `%v`",
		header, actuals, expected)
}

func TestGetIndex(t *testing.T) {
	response := NewServer(t).Get("/", nil)
	response.AssertStatusEquals(http.StatusOK)
	response.AssertBodyEquals("1.2.3.4\n")
	response.AssertHeaderContains("Content-Type", "text/plain; charset=UTF-8")
}

func TestGetInvalidUrl(t *testing.T) {
	response := NewServer(t).Get("/non-existant", nil)
	response.AssertStatusEquals(http.StatusNotFound)
	response.AssertBodyEquals("404 Not Found\n")
	response.AssertHeaderContains("Content-Type", "text/plain; charset=UTF-8")
}
