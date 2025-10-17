package main

import (
	"fmt"
	"log"
	"net/http"

	handlers "forum/internal/auth"
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

	mux.HandleFunc("GET /register/", handlers.RegisterHandlerGet)

	mux.HandleFunc("POST /register/", handlers.RegisterHandlerPost)

	mux.HandleFunc("GET  /login/", handlers.LoginHandlerGet)

	mux.HandleFunc("POST  /login/", handlers.LoginHandlerPost)

	mux.HandleFunc("GET /static/", home.StaticHandler)

	mux.HandleFunc("POST /api/comment", comments.SaveCommentHandler)

	mux.HandleFunc("GET /api/comments", comments.GetCommentsHandler)

	fmt.Println("Listening on http://localhost:8081")
	if err := http.ListenAndServe(":8081", mux); err != nil {
		log.Println(err)
	}
}
