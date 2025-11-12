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

	// Public routes
	http.HandleFunc("/", home.HomeHandler)
	http.HandleFunc("/register", auth.RegisterHandlerGet)
	http.HandleFunc("/api/register", auth.RegisterHandlerPost)
	http.HandleFunc("/login", auth.LoginHandlerGet)
	http.HandleFunc("/api/login", auth.LoginHandlerPost)
	http.HandleFunc("/logout", auth.LogoutHandler)
	http.HandleFunc("/Profile/{username}", posts.Profile)
	http.HandleFunc("/post-detail/", posts.SeePostdetail)
	http.HandleFunc("/static/", home.StaticHandler)
	http.HandleFunc("/liked-posts", likes.HandleLikedPosts)

	// Protected routes (require authentication)
	http.Handle("/api/comment",
		middleware.AuthMiddleware(http.HandlerFunc(comments.CreateCommentHandler)),
	)

	http.Handle("/api/create_post",
		middleware.AuthMiddleware(http.HandlerFunc(posts.CreatePostsHandler)),
	)

	http.Handle("/api/react",
		middleware.AuthMiddleware(http.HandlerFunc(likes.HandleLikeOrDislike)),
	)
	http.Handle("/api/Delete/{Postid}",
		middleware.AuthMiddleware(http.HandlerFunc(posts.PostDelete)),
	)

	fmt.Println("Listening on http://localhost:3001")
	if err := http.ListenAndServe(":3002", nil); err != nil {
		log.Println(err)
	}
}
