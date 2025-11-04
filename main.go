package main

import (
	"fmt"
	"log"
	"net/http"

	"forum/backend/auth"
	"forum/backend/comments"
	"forum/backend/database"
	"forum/backend/home"
	"forum/backend/likes"
	"forum/backend/middleware"
	"forum/backend/posts"
)

func main() {
	database.Init()

	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("GET /", home.HomeHandler)
	mux.HandleFunc("GET /register", auth.RegisterHandlerGet)
	mux.HandleFunc("POST /api/register", auth.RegisterHandlerPost)
	mux.HandleFunc("GET /login", auth.LoginHandlerGet)
	mux.HandleFunc("POST /api/login", auth.LoginHandlerPost)
	mux.HandleFunc("GET /logout", auth.LogoutHandler)
	mux.HandleFunc("GET /Profile/{username}", posts.Profile)
	mux.HandleFunc("GET /post-detail/", posts.SeePostdetail)
	mux.HandleFunc("GET /static/", home.StaticHandler)
mux.HandleFunc("GET /liked-posts", likes.HandleLikedPosts)
	// Protected routes (require authentication)
	mux.Handle("POST /api/comment",
		middleware.AuthMiddleware(http.HandlerFunc(comments.CreateCommentHandler)),
	)

	mux.Handle("POST /api/create_post",
		middleware.AuthMiddleware(http.HandlerFunc(posts.CreatePostsHandler)),
	)

	mux.Handle("POST /api/react",
		middleware.AuthMiddleware(http.HandlerFunc(likes.HandleLikeOrDislike)),
	)
	mux.Handle("POST /api/Delete/{Postid}",
		middleware.AuthMiddleware(http.HandlerFunc(posts.PostDelete)),
    )


	fmt.Println("Listening on http://localhost:3001")
	if err := http.ListenAndServe(":3001", mux); err != nil {
		log.Println(err)
	}
}
