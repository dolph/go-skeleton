package main

import (
    "testing"
    "net/http"
    "net/http/httptest"
    "strings"
)

type TestHandler struct {
    t *testing.T

    // The HTTP handler function to be tested.
    f func (http.ResponseWriter, *http.Request)
}

// Build an HTTP request, pass it to the HTTP handler, and return the response.
func (handler TestHandler) request(method string, path string, headers map[string]string) TestResponse {
    request, err := http.NewRequest(method, path, nil)
    if err != nil {
        handler.t.Fatal(err)
    }

    // Set request headers, if any.
    for header, value := range headers {
        request.Header.Set(header, value)
    }

    // Set an arbitrary remote address for logging purposes, etc.
    request.RemoteAddr = "1.2.3.4"

    response := httptest.NewRecorder()
    http.HandlerFunc(handler.f).ServeHTTP(response, request)
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

func TestGetIndex(t *testing.T) {
    response := TestHandler{t, IndexHandler}.Get("/", nil)
    response.AssertStatusEquals(http.StatusOK)
    response.AssertBodyEquals("1.2.3.4")
}

func TestGetInvalidUrl(t *testing.T) {
    response := TestHandler{t, IndexHandler}.Get("/non-existant", nil)
    response.AssertStatusEquals(http.StatusNotFound)
    response.AssertBodyEquals("404 Not Found")
}
