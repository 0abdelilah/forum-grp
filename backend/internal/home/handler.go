package home

import (
	"fmt"
	"html/template"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpt, err := template.ParseFiles("../frontend/templates/home.html")
	if err != nil {
		fmt.Println(err)
		return
	}

	tmpt.Execute(w, nil)
}
