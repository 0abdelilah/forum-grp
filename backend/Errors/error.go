package Errorhandel

import (
	"fmt"
	"html/template"
	"net/http"

	"forum/backend/models"
)

func Errordirect(w http.ResponseWriter, messageErr string, statuscode int) {
	template, err := template.ParseGlob("./frontend/templates/ErrorPage.html")
	if err != nil {
		fmt.Println("the err",err)
		http.Error(w, "internel server Error", http.StatusInternalServerError)
		return
	}
	Data := models.Error{
		MassageError: messageErr,
		Errorstatus:  statuscode,
	}
	w.WriteHeader(statuscode)
	if err := template.Execute(w, Data); err != nil {
		http.Error(w, "internel server Error", http.StatusInternalServerError)
		return
	}
}
