package home

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"forum/backend/database"
	"forum/backend/filters"
)

func PageNotFound(w http.ResponseWriter) {
	tmpt, err := template.ParseFiles("./frontend/templates/not-found.html")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	tmpt.Execute(w, nil)
}

func GetUsernameFromCookie(r *http.Request, cookie_name string) (string, error) {
	c, err := r.Cookie(cookie_name)
	if err != nil {
		return "", err
	}

	var username string
	err = database.Db.QueryRow("SELECT username FROM sessions WHERE id = ? AND expires_at > datetime('now')", c.Value).Scan(&username)
	if err == sql.ErrNoRows {
		return "", err //
	}
	if err != nil {
		return "", err
	}
	return username, nil
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		PageNotFound(w)
		return
	}

	tmpl, err := template.ParseFiles("./frontend/templates/index.html")
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	PageData := database.AllPageData(r, "HomeData")
	PageData.Username, err = GetUsernameFromCookie(r, "session_token")
	if err != nil {
		fmt.Println(err)
	}

	if err := tmpl.Execute(w, PageData); err != nil {
		log.Printf("template execution error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func PostHomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./frontend/templates/index.html")
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	PageData := database.AllPageData(r, "HomeData")

	//hna bax nfiltery bmethod post ghida nkamal
	r.ParseForm()
	if r.Form["Category"] != nil {
		PageData.AllPosts = filters.FelterbyCategory(PageData, r.Form["Category"][0])
		for i := 0; i < len(PageData.CategoryChoice); i++ {
			if PageData.CategoryChoice[i].Category == r.Form["Category"][0] {
				PageData.CategoryChoice[i].Selected = "true"
			}
		}
		fmt.Println(PageData)
		tmpl.Execute(w, PageData)
		return
	}

	if err := tmpl.Execute(w, PageData); err != nil {
		log.Printf("template execution error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func StaticHandler(w http.ResponseWriter, r *http.Request) {
	path := "./frontend/templates/" + r.URL.Path

	// Serve the file directly
	http.ServeFile(w, r, path)
}
