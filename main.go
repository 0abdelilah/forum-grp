package main

import (
	"fmt"
	"log"
	"net/http"

	"forum/backend/auth"
	"forum/backend/database"
)

func main() {
	database.Init()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /register", auth.RegisterHandlerGet)
	mux.HandleFunc("POST /api/register", auth.RegisterHandlerPost)
	mux.HandleFunc("GET  /login", auth.LoginHandlerGet)
	mux.HandleFunc("POST  /api/login", auth.LoginHandlerPost)
	mux.HandleFunc("GET  /logout", auth.LogoutHandler)

	fmt.Println("Listening on http://localhost:3001")
	if err := http.ListenAndServe(":3001", mux); err != nil {
		log.Println(err)
	}
}
