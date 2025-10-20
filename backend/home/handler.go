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

	IsLoggedIn := false
	var username string

	// Check if session cookie exists
	cookie, err := r.Cookie("session_token")
	if err == nil {
		// Validate session
		err = database.Db.QueryRow(`
			SELECT username FROM sessions WHERE id = ? AND expires_at > datetime('now')
		`, cookie.Value).Scan(&username)

		if err == nil {
			IsLoggedIn = true
		} else if err != sql.ErrNoRows {
			// Log unexpected DB error
			log.Printf("session lookup error: %v", err)
		}
	}
	PageData := database.AllPageData(r, "HomeData")
	if r.Method != http.MethodPost {
		tmpl.Execute(w, PageData)
		return
	}
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

	PageData.IsLoggedIn = IsLoggedIn
	PageData.Username = username

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
