package main

import (
	"book-api/router"
	"fmt"
	"net/http"
)

func main() {
	r := router.SetupRouter()

	port := ":8080"
	fmt.Printf("Server running at http://localhost%s\n", port)

	if err := http.ListenAndServe(port, r); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}
