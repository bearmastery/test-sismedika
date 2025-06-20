package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"book-api/handler"
	"book-api/model"

	"github.com/go-chi/chi/v5"
)

func requestWithBody(method, url string, body interface{}) *httptest.ResponseRecorder {
	var buf bytes.Buffer
	if body != nil {
		json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(method, url, &buf)
	rr := httptest.NewRecorder()

	r := chi.NewRouter()
	switch method {
	case http.MethodPost:
		r.Post("/books", handler.CreateBookHandler)
	case http.MethodPut:
		r.Put("/books/{id}", handler.UpdateBookHandler)
	case http.MethodDelete:
		r.Delete("/books/{id}", handler.DeleteBookHandler)
	}
	r.ServeHTTP(rr, req)
	return rr
}

func requestWithPath(method, path string, h http.HandlerFunc) *httptest.ResponseRecorder {
	r := chi.NewRouter()
	r.MethodFunc(method, "/books/{id}", h)
	req := httptest.NewRequest(method, path, nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	return rr
}

func TestGetBooksHandler(t *testing.T) {
	model.ResetBookStore()
	model.GetBookStore().AddBook(model.Book{Title: "Go Book", Author: "Riki", PublishedYear: 2024})

	req := httptest.NewRequest(http.MethodGet, "/books", nil)
	rr := httptest.NewRecorder()
	handler.GetBooksHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestGetBookHandlerInvalidID(t *testing.T) {
	rr := requestWithPath(http.MethodGet, "/books/abc", handler.GetBookHandler)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestGetBookHandlerNotFound(t *testing.T) {
	rr := requestWithPath(http.MethodGet, "/books/99", handler.GetBookHandler)
	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}
}

func TestGetBookHandler(t *testing.T) {
	model.ResetBookStore()
	book := model.GetBookStore().AddBook(model.Book{Title: "Go Book", Author: "Riki", PublishedYear: 2024})
	rr := requestWithPath(http.MethodGet, "/books/"+strconv.Itoa(book.ID), handler.GetBookHandler)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestCreateBookHandler(t *testing.T) {
	model.ResetBookStore()
	payload := model.Book{Title: "New Book", Author: "Author", PublishedYear: 2023}
	rr := requestWithBody(http.MethodPost, "/books", payload)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", rr.Code)
	}
}

func TestCreateBookHandlerInvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewBufferString(`{invalid`))
	rr := httptest.NewRecorder()
	handler.CreateBookHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestCreateBookHandlerMissingFields(t *testing.T) {
	payload := model.Book{}
	rr := requestWithBody(http.MethodPost, "/books", payload)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestUpdateBookHandler(t *testing.T) {
	model.ResetBookStore()
	book := model.GetBookStore().AddBook(model.Book{Title: "Old", Author: "Old", PublishedYear: 2000})

	updated := model.Book{Title: "Updated", Author: "New", PublishedYear: 2025}
	rr := requestWithBody(http.MethodPut, "/books/"+strconv.Itoa(book.ID), updated)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestUpdateBookHandlerInvalidBookID(t *testing.T) {
	body := model.Book{Title: "Dummy", Author: "X", PublishedYear: 2020}
	rr := requestWithBody(http.MethodPut, "/books/abc", body)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestUpdateBookHandlerInvalidBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPut, "/books/1", strings.NewReader("{invalid json"))
	rr := httptest.NewRecorder()
	r := chi.NewRouter()
	r.Put("/books/{id}", handler.UpdateBookHandler)
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestUpdateBookHandlerRequiredBody(t *testing.T) {
	body := model.Book{}
	rr := requestWithBody(http.MethodPut, "/books/1", body)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestUpdateBookHandlerNotFound(t *testing.T) {
	body := model.Book{Title: "Nope", Author: "Ghost", PublishedYear: 2020}
	rr := requestWithBody(http.MethodPut, "/books/999", body)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}
}

func TestDeleteBookHandler(t *testing.T) {
	model.ResetBookStore()
	book := model.GetBookStore().AddBook(model.Book{Title: "ToDelete", Author: "Someone", PublishedYear: 2022})
	rr := requestWithPath(http.MethodDelete, "/books/"+strconv.Itoa(book.ID), handler.DeleteBookHandler)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestDeleteBookHandlerInvalidID(t *testing.T) {
	rr := requestWithPath(http.MethodDelete, "/books/abc", handler.DeleteBookHandler)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestDeleteBookHandlerNotFound(t *testing.T) {
	rr := requestWithPath(http.MethodDelete, "/books/999", handler.DeleteBookHandler)
	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}
}
