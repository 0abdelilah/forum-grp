package main

import (
	"fmt"
	"net/http"

	"forum/internal/home"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", home.HomeHandler)

	fmt.Println("Listening on http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
