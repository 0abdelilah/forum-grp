package main

import (
	"fmt"
	"log"
	"net/http"

	"forum/backend/auth"
	"forum/backend/comments"
	databasecreate "forum/backend/database"
	"forum/backend/home"
	"forum/backend/likes"
	"forum/backend/posts"
)

func main() {
	databasecreate.Init()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /register", auth.RegisterHandlerGet)
	mux.HandleFunc("POST /api/register", auth.RegisterHandlerPost)
	mux.HandleFunc("GET  /login", auth.LoginHandlerGet)
	mux.HandleFunc("POST  /api/login", auth.LoginHandlerPost)
	mux.HandleFunc("GET  /logout", auth.LogoutHandler)

	// Posts
	mux.HandleFunc("POST /api/create_post", posts.CreatePostsHandler)
	mux.HandleFunc("GET /Profile/{username}", posts.Profile)

	mux.HandleFunc("POST /api/like", likes.LikeHandler)
	mux.HandleFunc("POST /api/dislike", likes.DislikeHandler)

	mux.HandleFunc("POST /api/comment", comments.CreateCommentHandler)

	mux.HandleFunc("POST /api/like_comment", likes.AddCommentLikeHandler)
	mux.HandleFunc("POST /api/dislike_comment", likes.AddCommentDislikeHandler)

	// Posts content
	mux.HandleFunc("GET /post-detail/", posts.SeePostdetail)

	// static files
	mux.HandleFunc("GET /static/", home.StaticHandler)

	fmt.Println("Listening on http://localhost:3001")
	if err := http.ListenAndServe(":3001", mux); err != nil {
		log.Println(err)
	}
}
