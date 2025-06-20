package router

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRouterRoutes(t *testing.T) {
	router := SetupRouter()

	tests := []struct {
		name       string
		method     string
		path       string
		body       string
		wantStatus int
	}{
		{"GET /books", http.MethodGet, "/books", "", http.StatusOK},
		{"POST /books", http.MethodPost, "/books", `{"title":"Go","author":"Riki","published_year":2024}`, http.StatusCreated},
		{"PUT /books/1", http.MethodPut, "/books/1", `{"title":"Updated","author":"Updated","published_year":2024}`, http.StatusOK},
		{"GET /books/1", http.MethodGet, "/books/1", "", http.StatusOK},
		{"DELETE /books/1", http.MethodDelete, "/books/1", "", http.StatusOK},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			res := httptest.NewRecorder()

			router.ServeHTTP(res, req)

			if res.Code != tc.wantStatus {
				t.Errorf("unexpected status: got %v, want %v", res.Code, tc.wantStatus)
			}
		})
	}
}
