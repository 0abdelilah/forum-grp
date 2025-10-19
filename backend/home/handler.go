package home

import (
	"forum/backend/database"
	"html/template"
	"net/http"
)

func PageNotFound(w http.ResponseWriter) {
	tmpt, err := template.ParseFiles("./frontend/templates/not-found.html")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	tmpt.Execute(w, nil)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		PageNotFound(w)
		return
	}

	tmpl, err := template.ParseFiles("./frontend/templates/index.html")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	IsLoggedIn := false

	// Check if session cookie exists
	cookie, err := r.Cookie("session_token")
	if err == nil {
		// Verify that session is valid and not expired
		var userID int
		err = database.Db.QueryRow(
			`SELECT user_id FROM sessions WHERE id = ? AND expires_at > datetime('now')`,
			cookie.Value,
		).Scan(&userID)

		if err == nil {
			IsLoggedIn = true
		}
	}

	PageData := database.AllPageData(r, "HomeData")
	PageData.IsLoggedIn = IsLoggedIn
	PageData.Username = "Username"

	tmpl.Execute(w, PageData)
}

func StaticHandler(w http.ResponseWriter, r *http.Request) {
	path := "./frontend/templates/" + r.URL.Path

	// Serve the file directly
	http.ServeFile(w, r, path)
}
