package home

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"

	Errorhandel "forum/backend/Errors"
	"forum/backend/auth"
	"forum/backend/database"
)

// the Home function
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		Errorhandel.Errordirect(w, "Page not found", http.StatusNotFound)
		return
	}
	tmpl, err := template.ParseFiles("./frontend/templates/index.html")
	if err != nil {
		Errorhandel.Errordirect(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	PageData := database.AllPageData(r, "HomeData")
	username, err := auth.GetUsernameFromCookie(r, "session_token")
	PageData.Username = username
	if err != nil && err != sql.ErrNoRows && fmt.Sprintf("%v", err) != "http: named cookie not present" {
		Errorhandel.Errordirect(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, PageData); err != nil {
		Errorhandel.Errordirect(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// serve static files
func StaticHandler(w http.ResponseWriter, r *http.Request) {
	path := "./frontend/templates/" + strings.TrimSuffix(r.URL.Path, "/")
	if f, err := os.Stat(path); err != nil || f.IsDir() {
		Errorhandel.Errordirect(w, "Page Not Found", http.StatusNotFound)
		return
	}
	http.ServeFile(w, r, path)
}
