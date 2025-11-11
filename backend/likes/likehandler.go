package likes

import (
	"database/sql"
	"fmt"
	"net/http"
	"text/template"

	Errorhandel "forum/backend/Errors"
	"forum/backend/auth"
	"forum/backend/database"
	"forum/backend/models"
)

func HandleLikedPosts(w http.ResponseWriter, r *http.Request) {
	username, err := auth.GetUsernameFromCookie(r, "session_token")
	if err != nil && err != sql.ErrNoRows && fmt.Sprintf("%v", err) != "http: named cookie not present" {
		Errorhandel.Errordirect(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	likedPosts, err := database.GetAlllike(1, username)
	if err != nil {
		fmt.Println("err", err)
		Errorhandel.Errordirect(w, "InternalServerError", http.StatusInternalServerError)
	}

	data := struct {
		Username string
		Posts    []models.Post
	}{
		Username: username,
		Posts:    likedPosts,
	}

	tmpl, err := template.ParseFiles("./frontend/templates/liked-posts.html")
	if err != nil {
		Errorhandel.Errordirect(w, "InternalServerError", http.StatusInternalServerError)
	}
	tmpl.Execute(w, data)
}
