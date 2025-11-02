package likes

import (
	"forum/backend/auth"
	"forum/backend/database"
	"forum/backend/models"
	"net/http"
	"text/template"
)

func HandleLikedPosts(w http.ResponseWriter, r *http.Request) {
	username, _ := auth.GetUsernameFromCookie(r, "session_token")
	if username == "" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	likedPosts := database.GetAlllike(1, "post", username)
	if likedPosts == nil {
		http.Error(w, "No liked posts found", http.StatusNotFound)
		return
	}

	data := struct {
		Username string
		Posts    []models.LikedPosts
	}{
		Username: username,
		Posts:    likedPosts,
	}

	tmpl := template.Must(template.ParseFiles("./frontend/templates/liked-posts.html"))
	tmpl.Execute(w, data)
}
