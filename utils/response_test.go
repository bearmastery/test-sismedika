package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteJSON(t *testing.T) {
	rr := httptest.NewRecorder()

	data := map[string]string{
		"message": "hello world",
	}

	WriteJSON(rr, http.StatusOK, data)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rr.Code)
	}

	if contentType := rr.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}

	var response map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response["message"] != "hello world" {
		t.Errorf("Expected message 'hello world', got '%s'", response["message"])
	}
}

func TestWriteError(t *testing.T) {
	rr := httptest.NewRecorder()

	WriteError(rr, http.StatusBadRequest, "invalid request")

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}

	if contentType := rr.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}

	var response map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response["error"] != "invalid request" {
		t.Errorf("Expected error message 'invalid request', got '%s'", response["error"])
	}
}

func TestWriteJSONEncodeError(t *testing.T) {
	rr := httptest.NewRecorder()

	data := func() {}

	WriteJSON(rr, http.StatusOK, data)

	expected := "failed to encode JSON\n"
	if rr.Body.String() != expected {
		t.Errorf("Expected body %q, got %q", expected, rr.Body.String())
	}
}
