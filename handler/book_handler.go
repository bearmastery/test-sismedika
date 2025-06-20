package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"book-api/model"
	"book-api/utils"

	"github.com/go-chi/chi/v5"
)

func GetBooksHandler(w http.ResponseWriter, r *http.Request) {
	store := model.GetBookStore()
	books := store.GetAllBooks()
	utils.WriteJSON(w, http.StatusOK, books)
}

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
