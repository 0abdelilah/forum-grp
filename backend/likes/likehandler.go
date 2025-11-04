package likes

import (
	"net/http"
	"text/template"

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

	likedPosts := database.GetAlllike(1, "post", username)

	data := struct {
		Username string
		Posts    []models.Post
	}{
		Username: username,
		Posts:    likedPosts,
	}

	tmpl := template.Must(template.ParseFiles("./frontend/templates/liked-posts.html"))
	tmpl.Execute(w, data)
}
