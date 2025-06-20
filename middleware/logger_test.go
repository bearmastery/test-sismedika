package middleware

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLoggerMiddleware(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	req := httptest.NewRequest(http.MethodGet, "/test-path", nil)
	rr := httptest.NewRecorder()

	loggedHandler := LoggerMiddleware(handler)

	loggedHandler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", rr.Code)
	}

	logOutput := buf.String()
	if !strings.Contains(logOutput, "GET /test-path") {
		t.Errorf("Expected log to contain 'GET /test-path', got %q", logOutput)
	}
}
