package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"book-api/model"
	"book-api/utils"

	"github.com/go-chi/chi/v5"
)

// GetBooksHandler menangani permintaan GET /books.
//
// Params:
//   - w: http.ResponseWriter untuk menulis response ke client.
//   - r: *http.Request yang berisi informasi request dari client.
func GetBooksHandler(w http.ResponseWriter, r *http.Request) {
	store := model.GetBookStore()
	books := store.GetAllBooks()
	utils.WriteJSON(w, http.StatusOK, books)
}

// GetBookHandler menangani permintaan GET /books/{id}.
//
// Params:
//   - w: http.ResponseWriter untuk menulis response ke client.
//   - r: *http.Request yang mengandung parameter URL "id"
//
// Response:
//   - 200 OK jika buku ditemukan
//   - 400 Bad Request jika ID tidak valid
//   - 404 Not Found jika buku tidak ditemukan
func GetBookHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid book ID")
		return
	}

	store := model.GetBookStore()
	book, err := store.GetBookByID(id)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, book)
}

// CreateBookHandler menangani permintaan POST /books untuk menambahkan buku baru.
//
// Params:
//   - w: http.ResponseWriter untuk menulis response ke client.
//   - r: *http.Request yang membawa data JSON dari body.
//
// Response:
//   - 201 Created jika sukses
//   - 400 Bad Request jika body tidak valid atau field kosong
func CreateBookHandler(w http.ResponseWriter, r *http.Request) {
	var book model.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if book.Title == "" || book.Author == "" || book.PublishedYear == 0 {
		utils.WriteError(w, http.StatusBadRequest, "all fields are required")
		return
	}

	store := model.GetBookStore()
	created := store.AddBook(book)
	utils.WriteJSON(w, http.StatusCreated, created)
}

// UpdateBookHandler menangani permintaan PUT /books/{id} untuk memperbarui data buku.
//
// Params:
//   - w: http.ResponseWriter untuk menulis response ke client.
//   - r: *http.Request yang mengandung parameter URL "id" dan data JSON baru dari body.
//
// Response:
//   - 200 OK jika update berhasil
//   - 400 Bad Request jika ID/body tidak valid atau field kosong
//   - 404 Not Found jika ID buku tidak ditemukan
func UpdateBookHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid book ID")
		return
	}

	var book model.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if book.Title == "" || book.Author == "" || book.PublishedYear == 0 {
		utils.WriteError(w, http.StatusBadRequest, "all fields are required")
		return
	}

	store := model.GetBookStore()
	updated, err := store.UpdateBook(id, book)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, updated)
}

// DeleteBookHandler menangani permintaan DELETE /books/{id} untuk menghapus buku.
//
// Params:
//   - w: http.ResponseWriter untuk menulis response ke client.
//   - r: *http.Request yang mengandung parameter URL "id".
//
// Response:
//   - 200 OK jika buku berhasil dihapus
//   - 400 Bad Request jika ID tidak valid
//   - 404 Not Found jika ID buku tidak ditemukan
func DeleteBookHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid book ID")
		return
	}

	store := model.GetBookStore()
	err = store.DeleteBook(id)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "book deleted"})
}
