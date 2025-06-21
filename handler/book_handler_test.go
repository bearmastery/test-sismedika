package handler_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"book-api/handler"
	"book-api/model"

	"github.com/go-chi/chi/v5"
)

var service = model.NewBookStore()
var bookHandler = handler.NewBookHandler(service)
var mockBooks = []model.Book{
	{ID: 1, Title: "Book 1", Author: "Author A", PublishedYear: 2020},
	{ID: 2, Title: "Book 2", Author: "Author B", PublishedYear: 2021},
	{ID: 3, Title: "Book 3", Author: "Author C", PublishedYear: 2022},
}

func setupRequestWithID(method, path string, body *bytes.Buffer, handlerFunc http.HandlerFunc) *httptest.ResponseRecorder {
	r := chi.NewRouter()
	r.Route("/books", func(r chi.Router) {
		r.MethodFunc(method, "/{id}", handlerFunc)
	})

	var reader io.Reader
	if body != nil {
		reader = body
	}

	req := httptest.NewRequest(method, path, reader)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	return rr
}

func setupHandlerWithData() {

	for _, b := range mockBooks {
		service.AddBook(b)
	}
}

func TestGetBooksHandler_Success(t *testing.T) {
	setupHandlerWithData()
	req := httptest.NewRequest("GET", "/books", nil)
	rr := httptest.NewRecorder()

	bookHandler.GetBooksHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("GetBooks: expected 200, got %d", rr.Code)
	}

	var response struct {
		Data  []model.Book `json:"data"`
		Error string       `json:"error,omitempty"`
	}

	json.NewDecoder(rr.Body).Decode(&response)

	if len(response.Data) != len(mockBooks) {
		t.Errorf("GetBooks: unexpected response length: got %+v, want %+v", len(response.Data), len(mockBooks))
	}
}

func TestCreateBookHandler_Success(t *testing.T) {
	book := model.Book{Title: "Test Book", Author: "Tester", PublishedYear: 2023}
	body, _ := json.Marshal(book)

	req := httptest.NewRequest("POST", "/books", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	bookHandler.CreateBookHandler(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("CreateBook: expected 201, got %d", rr.Code)
	}

	var response struct {
		Data  model.Book `json:"data"`
		Error string     `json:"error,omitempty"`
	}

	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("CreateBook: failed to decode response: %v", err)
	}

	if response.Data.ID == 0 || response.Data.Title != "Test Book" {
		t.Errorf("CreateBook: unexpected book data: %+v", response.Data)
	}
}

func TestCreateBookHandler_InvalidBody(t *testing.T) {
	req := httptest.NewRequest("POST", "/books", nil)
	rr := httptest.NewRecorder()

	bookHandler.CreateBookHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("CreateBook: expected 400, got %d", rr.Code)
	}
}

func TestCreateBookHandler_BadRequest(t *testing.T) {
	book := model.Book{}
	body, _ := json.Marshal(book)

	req := httptest.NewRequest("POST", "/books", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	bookHandler.CreateBookHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("CreateBook: expected 400, got %d", rr.Code)
	}
}

func TestGetBookHandler_Success(t *testing.T) {

	rr := setupRequestWithID("GET", "/books/"+strconv.Itoa(mockBooks[0].ID), nil, bookHandler.GetBookHandler)

	if rr.Code != http.StatusOK {
		t.Errorf("GetBook: expected 200, got %d", rr.Code)
	}
}

func TestGetBookHandler_InvalidID(t *testing.T) {
	rr := setupRequestWithID("GET", "/books/abc", nil, bookHandler.GetBookHandler)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("GetBook: expected 400, got %d", rr.Code)
	}
}

func TestGetBookHandler_NotFound(t *testing.T) {
	rr := setupRequestWithID("GET", "/books/999", nil, bookHandler.GetBookHandler)

	if rr.Code != http.StatusNotFound {
		t.Errorf("GetBook: expected 404, got %d", rr.Code)
	}
}

func TestUpdateBookHandler_Success(t *testing.T) {

	updated := model.Book{Title: "Updated", Author: "Someone", PublishedYear: 2000}
	body, _ := json.Marshal(updated)

	rr := setupRequestWithID("PUT", "/books/"+strconv.Itoa(mockBooks[0].ID), bytes.NewBuffer(body), bookHandler.UpdateBookHandler)

	if rr.Code != http.StatusOK {
		t.Errorf("UpdateBook: expected 200, got %d", rr.Code)
	}

	var response struct {
		Data  model.Book `json:"data"`
		Error string     `json:"error,omitempty"`
	}

	json.NewDecoder(rr.Body).Decode(&response)

	if response.Data.Title != "Updated" {
		t.Errorf("UpdateBook Title: not updated correctly: %+v", response.Data.Title)
	}

	if response.Data.Author != "Someone" {
		t.Errorf("UpdateBook Author: not updated correctly: %+v", response.Data.Author)
	}

	if response.Data.PublishedYear != 2000 {
		t.Errorf("UpdateBook PublishedYear: not updated correctly: %+v", response.Data.PublishedYear)
	}
}

func TestUpdateBookHandler_InvalidID(t *testing.T) {
	updated := model.Book{Title: "Updated", Author: "Someone", PublishedYear: 2000}
	body, _ := json.Marshal(updated)

	rr := setupRequestWithID("PUT", "/books/abc", bytes.NewBuffer(body), bookHandler.UpdateBookHandler)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("UpdateBook: expected 400, got %d", rr.Code)
	}

}

func TestUpdateBookHandler_InvalidBody(t *testing.T) {
	rr := setupRequestWithID("PUT", "/books/1", nil, bookHandler.UpdateBookHandler)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("UpdateBook: expected 400, got %d", rr.Code)
	}

}

func TestUpdateBookHandler_BadRequest(t *testing.T) {
	updated := model.Book{}

	body, _ := json.Marshal(updated)
	rr := setupRequestWithID("PUT", "/books/1", bytes.NewBuffer(body), bookHandler.UpdateBookHandler)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("UpdateBook: expected 400, got %d", rr.Code)
	}

}

func TestUpdateBookHandler_NotFound(t *testing.T) {
	updated := model.Book{Title: "Updated", Author: "Someone", PublishedYear: 2000}
	body, _ := json.Marshal(updated)

	rr := setupRequestWithID("PUT", "/books/999", bytes.NewBuffer(body), bookHandler.UpdateBookHandler)

	if rr.Code != http.StatusNotFound {
		t.Errorf("UpdateBook: expected 404, got %d", rr.Code)
	}
}

func TestDeleteBookHandler_Success(t *testing.T) {
	rr := setupRequestWithID("DELETE", "/books/"+strconv.Itoa(mockBooks[0].ID), nil, bookHandler.DeleteBookHandler)

	if rr.Code != http.StatusOK {
		t.Errorf("DeletedBook: expected 200, got %d", rr.Code)
	}

	_, err := service.GetBookByID(mockBooks[0].ID)

	if err == nil {
		t.Errorf("DeletedBook: expected book to be deleted, but it still exists")
	}
}

func TestDeleteBookHandler_InvalidID(t *testing.T) {
	rr := setupRequestWithID("DELETE", "/books/abc", nil, bookHandler.DeleteBookHandler)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("DeletedBook: expected 400, got %d", rr.Code)
	}
}

func TestDeleteBookHandler_NotFound(t *testing.T) {
	rr := setupRequestWithID("DELETE", "/books/999", nil, bookHandler.DeleteBookHandler)

	if rr.Code != http.StatusNotFound {
		t.Errorf("DeletedBook: expected 404, got %d", rr.Code)
	}
}
