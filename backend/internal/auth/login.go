package handlers

import (
	"log"
	"net/http"

	"forum/internal/database"

	"github.com/google/uuid"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	userid := 0
	switch r.Method {
	case "GET":

		http.ServeFile(w, r, "frontend/templates/login.html")
		return
	case "POST":
		r.ParseForm()
		username := r.FormValue("username")
		password := r.FormValue("password")
		if username == "" || password == "" {
			http.Error(w, "Username and password are required", http.StatusBadRequest)
			return
		}
		err := database.Db.QueryRow(`SELECT id FROM users WHERE username = ? AND password = ?`, username, password).Scan(&userid)
		if err != nil {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			log.Println("Login failed: wrong username/password")
			return
		}
		u2 := uuid.New().String()
		database.Db.Exec(`UPDATE users SET session_token = ? WHERE id = ?`, u2, userid)
		http.SetCookie(w, &http.Cookie{
			Name:  "session_token",
			Value: u2,
		})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

	}
}
