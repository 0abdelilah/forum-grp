package main

import (
	"fmt"
	"net/http"

	"forum/internal/comments"
	"forum/internal/database"
	"forum/internal/home"
)

func main() {
	db, err := database.Init()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer db.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", home.HomeHandler)

	mux.HandleFunc("GET /static/", home.StaticHandler)

	mux.HandleFunc("POST /api/comment", comments.SaveCommentHandler)

	mux.HandleFunc("GET /api/comments", comments.GetCommentsHandler)

	fmt.Println("Listening on http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}
