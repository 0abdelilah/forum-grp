package likes

import (
	"net/http"
	"text/template"

	Errorhandel "forum/backend/Errors"
	"forum/backend/auth"
	"forum/backend/database"
	"forum/backend/models"
)

func HandleLikedPosts(w http.ResponseWriter, r *http.Request) {
	username, _ := auth.GetUsernameFromCookie(r, "session_token")
	if username == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	likedPosts, err := database.GetAlllike(1, username)
	if err != nil {
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
