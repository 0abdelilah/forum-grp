package home

import (
	"html/template"
	"net/http"

	"forum/database"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpt, err := template.ParseFiles("../frontend/templates/index.html")
	if err != nil {
		http.Error(w, "Page not found", http.StatusInternalServerError)
		return
	}
	if r.URL.Path != "/" {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}
	PageData := database.ALlPageData(r, "HomeData")
	tmpt.Execute(w, PageData)
}

func StaticHandler(w http.ResponseWriter, r *http.Request) {
	path := "../frontend/templates/" + r.URL.Path

	// Serve the file directly
	http.ServeFile(w, r, path)
}
