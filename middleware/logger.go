package middleware

import (
	"log"
	"net/http"
	"time"
)

// LoggerMiddleware mencatat metode HTTP, path, dan durasi waktu eksekusi setiap request.
func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)

		log.Printf("%s %s - %v\n", r.Method, r.URL.Path, duration)
	})
}
