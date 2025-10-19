package home

import (
	"database/sql"
	"forum/backend/database"
	"html/template"
	"log"
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
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	isLoggedIn := false
	var username string

	// Check if session cookie exists
	cookie, err := r.Cookie("session_token")
	if err == nil {
		// Validate session
		err = database.Db.QueryRow(`
			SELECT username FROM sessions WHERE id = ? AND expires_at > datetime('now')
		`, cookie.Value).Scan(&username)

		if err == nil {
			isLoggedIn = true
		} else if err != sql.ErrNoRows {
			// Log unexpected DB error
			log.Printf("session lookup error: %v", err)
		}
	}

	pageData := database.AllPageData(r, "HomeData")
	pageData.IsLoggedIn = isLoggedIn
	pageData.Username = username

	if err := tmpl.Execute(w, pageData); err != nil {
		log.Printf("template execution error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func StaticHandler(w http.ResponseWriter, r *http.Request) {
	path := "./frontend/templates/" + r.URL.Path

	// Serve the file directly
	http.ServeFile(w, r, path)
}
