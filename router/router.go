package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"book-api/handler"
	middleware2 "book-api/middleware"
)

func SetupRouter() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware2.LoggerMiddleware)

	r.Route("/books", func(r chi.Router) {
		r.Get("/", handler.GetBooksHandler)
		r.Get("/{id}", handler.GetBookHandler)
		r.Post("/", handler.CreateBookHandler)
		r.Put("/{id}", handler.UpdateBookHandler)
		r.Delete("/{id}", handler.DeleteBookHandler)
	})

	return r
}
