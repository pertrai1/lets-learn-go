package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pertrai1/snippetbox/internal/assert"
)

func TestMiddleware(t *testing.T) {
	// initialize a new httptest.ResponseRecorder and dummy http.Request
	rr := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// create a mock http handler that we can pass to our secureHeaders
	// middleware, which writes a 200 status code and "OK" response body.
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// pass the mock http handler to our secureHeaders middleware. Because
	// secureHeaders returns a http.Handler we can call serveHTTP() method.
	// passing in the http.ResponseRecorder and dummy http.Request to execute it.
	secureHeaders(next).ServeHTTP(rr, r)

	// call the Result() method on http.ResponseRecorder to get the results
	// of the test
	rs := rr.Result()

	// check that the middleware has correctly set the Content-Security-Policy
	// header on the response
	expectedValue := "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com"
	assert.Equal(t, rs.Header.Get("Content-Security-Policy"), expectedValue)

	// check that the middleware has correctly set the Referrer-Policy
	// header on the response
	expectedValue = "origin-when-cross-origin"
	assert.Equal(t, rs.Header.Get("Referrer-Policy"), expectedValue)

	// check that the middleware has correctly set the X-Content-Type-Options
	// header on the response
	expectedValue = "nosniff"
	assert.Equal(t, rs.Header.Get("X-Content-Type-Options"), expectedValue)

	// check that the middleware correctly set the X-Frame-Options header
	// on the response
	expectedValue = "deny"
	assert.Equal(t, rs.Header.Get("X-Frame-Options"), expectedValue)

	// check that the middleware correctly sets the X-XSS-Protection header
	// on the response
	expectedValue = "0"
	assert.Equal(t, rs.Header.Get("X-XSS-Protection"), expectedValue)

	// check that the middleware has correctly set the next handler in line
	// and the response status code and body are as expected
	assert.Equal(t, rs.StatusCode, http.StatusOK)

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)

	assert.Equal(t, string(body), "OK")

}
