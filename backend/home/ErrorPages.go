package home

import (
	Errorhandel "forum/backend/Errors"
	"forum/backend/auth"
	"forum/backend/database"
	"html/template"
	"log"
	"net/http"
)

// func PageNotFound(w http.ResponseWriter) {
// 	tmpt, err := template.ParseFiles("./frontend/templates/not-found.html")
// 	if err != nil {
// 		http.Error(w, "Internal server error", http.StatusInternalServerError)
// 		return
// 	}
// 	tmpt.Execute(w, nil)
// }

func HomePageError(w http.ResponseWriter, r *http.Request, Error string ,statuscode int) {
	tmpl, err := template.ParseFiles("./frontend/templates/index.html")
	if err != nil {
		Errorhandel.Errordirect(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(statuscode)
	// Get the normal page data
	PageData := database.AllPageData(r, "HomeData")
	PageData.Username, _ = auth.GetUsernameFromCookie(r, "session_token")

	// Add error
	PageData.Error = Error
	// Execute template
	if err := tmpl.Execute(w, PageData); err != nil {
		log.Printf("template execution error: %v", err)
		Errorhandel.Errordirect(w, "Internal server error", http.StatusInternalServerError)
	}
}



