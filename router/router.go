package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"book-api/handler"
	middleware2 "book-api/middleware"
	"book-api/model"
)

// SetupRouter mengatur dan mengembalikan konfigurasi HTTP router utama.
// Fungsi ini menggunakan chi router dan menambahkan middleware serta route untuk resource /books.
func SetupRouter() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware2.LoggerMiddleware)

	bookService := model.NewBookStore()
	bookHandler := handler.NewBookHandler(bookService)

	r.Route("/books", func(r chi.Router) {
		r.Get("/", bookHandler.GetBooksHandler)
		r.Get("/{id}", bookHandler.GetBookHandler)
		r.Post("/", bookHandler.CreateBookHandler)
		r.Put("/{id}", bookHandler.UpdateBookHandler)
		r.Delete("/{id}", bookHandler.DeleteBookHandler)
	})

	return r
}
