package home

import (
	"fmt"
	"html/template"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpt, err := template.ParseFiles("../frontend/templates/index.html")
	if err != nil {
		fmt.Println(err)
		return
	}

	tmpt.Execute(w, nil)
}

func StaticHandler(w http.ResponseWriter, r *http.Request) {
	path := "../frontend/templates/" + r.URL.Path

	// Serve the file directly
	http.ServeFile(w, r, path)
}
