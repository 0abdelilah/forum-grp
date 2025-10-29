package home

import (
	"html/template"
	"log"
	"net/http"

	"forum/backend/auth"
	databasecreate "forum/backend/database"
)

func HomePageError(w http.ResponseWriter, r *http.Request, Error string) {
	tmpl, err := template.ParseFiles("./frontend/templates/index.html")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Get the normal page data
	PageData := databasecreate.AllPageData(r, "HomeData")
	PageData.Username, _ = auth.GetUsernameFromCookie(r, "session_token")

	// Add error
	PageData.Error = Error

	// Execute template
	if err := tmpl.Execute(w, PageData); err != nil {
		log.Printf("template execution error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func PostPageError(w http.ResponseWriter, r *http.Request, Error string) {
	tmpl, err := template.ParseFiles("./frontend/templates/post-detail.html")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Get the normal page data
	PageData := databasecreate.AllPageData(r, "postContent")

	PageData.Username, _ = auth.GetUsernameFromCookie(r, "session_token")
	// Add error
	PageData.Error = Error

	// Execute template
	if err := tmpl.Execute(w, PageData); err != nil {
		log.Printf("template execution error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
