package auth

import (
	"database/sql"
	"errors"
	"html/template"
	"log"
	"net/http"
	"strings"

	Errorhandel "forum/backend/Errors"
	"forum/backend/database"
	"forum/backend/models"

	"golang.org/x/crypto/bcrypt"
)

func LoginHandlerGet(w http.ResponseWriter, r *http.Request) {
	username, err := GetUsernameFromCookie(r, "session_token")
	if err != nil {
		if err != sql.ErrNoRows && err != errors.New("http: named cookie not present") {
			Errorhandel.Errordirect(w, "InternalServerError", http.StatusInternalServerError)
			return
		}
	} else if username != "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	tmpt, err := template.ParseFiles("./frontend/templates/login.html")
	if err != nil {
		Errorhandel.Errordirect(w, "Internal Server Error ", http.StatusInternalServerError)
		return
	}
	tmpt.Execute(w, nil)
}

func LoginHandlerPost(w http.ResponseWriter, r *http.Request) {
	tmpt, err := template.ParseFiles("./frontend/templates/login.html")
	if err != nil {
		Errorhandel.Errordirect(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		tmpt.Execute(w, models.ErrorData{Error: "Invalid form data"})
		return
	}

	var login models.LoginData
	login.LoginInput = strings.TrimSpace(r.FormValue("username"))
	login.Password = r.FormValue("password")

	if login.LoginInput == "" || login.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		tmpt.Execute(w, models.ErrorData{Error: "Username and password cannot be empty"})
		return
	}

	err = database.Db.QueryRow(
		`SELECT username, password_hash FROM users WHERE username = ? OR email = ?`,
		login.LoginInput, login.LoginInput,
	).Scan(&login.Username, &login.PasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusUnauthorized)
			tmpt.Execute(w, models.ErrorData{Error: "Invalid username or password"})
			return
		} else {
			log.Println("DB error:", err)
			w.WriteHeader(http.StatusInternalServerError)
			tmpt.Execute(w, models.ErrorData{Error: "Internal server error"})
			return
		}
	}

	if err = bcrypt.CompareHashAndPassword([]byte(login.PasswordHash), []byte(login.Password)); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		tmpt.Execute(w, models.ErrorData{Error: "Invalid username or password"})
		return
	}

	err = CreateSession(login.Username, w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		tmpt.Execute(w, models.ErrorData{Error: "Could not create session, try again later"})
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
