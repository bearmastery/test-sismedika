package utils_test

import (
	"book-api/utils"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWriteJSON_Success(t *testing.T) {
	rr := httptest.NewRecorder()
	data := map[string]string{"message": "success"}

	utils.WriteJSON(rr, http.StatusOK, data)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	if contentType := rr.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", contentType)
	}

	var resp utils.APIResponse
	err := json.NewDecoder(rr.Body).Decode(&resp)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	dataMap, ok := resp.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("expected map in response data, got: %T", resp.Data)
	}
	if dataMap["message"] != "success" {
		t.Errorf("expected message 'success', got '%v'", dataMap["message"])
	}

	if resp.Error != "" {
		t.Errorf("expected no error, got: %s", resp.Error)
	}
}

func TestWriteError_Success(t *testing.T) {
	rr := httptest.NewRecorder()
	utils.WriteError(rr, http.StatusBadRequest, "invalid input")

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}

	if contentType := rr.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", contentType)
	}

	var resp utils.APIResponse
	err := json.NewDecoder(rr.Body).Decode(&resp)
	if err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}

	if resp.Error != "invalid input" {
		t.Errorf("expected error message 'invalid input', got '%s'", resp.Error)
	}

	if resp.Data != nil {
		t.Errorf("expected data to be nil, got: %v", resp.Data)
	}
}

func TestWriteJSON_EncodingError(t *testing.T) {
	rr := httptest.NewRecorder()

	utils.WriteJSON(rr, http.StatusOK, make(chan int))

	res := rr.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status 500 on encoding error, got %d", res.StatusCode)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "failed to encode response") {
		t.Errorf("expected encoding error message, got: %s", body)
	}
}
