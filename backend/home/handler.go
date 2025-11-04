package home

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	Errorhandel "forum/backend/Errors"
	"forum/backend/auth"
	"forum/backend/database"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		Errorhandel.Errordirect(w, "Page not Found", http.StatusNotFound)
		fmt.Println("fhsghfwe")
		return
	}
	tmpl, err := template.ParseFiles("./frontend/templates/index.html")
	if err != nil {
		// http.Error(w, "internal server error", http.StatusInternalServerError)
		Errorhandel.Errordirect(w, "Internal server error", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	PageData := database.AllPageData(r, "HomeData")
	PageData.Username, err = auth.GetUsernameFromCookie(r, "session_token")
	if err != nil {
		fmt.Println(err)
	}

	if err := tmpl.Execute(w, PageData); err != nil {
		log.Printf("template execution error: %v", err)
		Errorhandel.Errordirect(w, "Internal server error", http.StatusInternalServerError)
	}
}

func StaticHandler(w http.ResponseWriter, r *http.Request) {
	path := "./frontend/templates/" + r.URL.Path
	// Serve the file directly
	http.ServeFile(w, r, path)
}
