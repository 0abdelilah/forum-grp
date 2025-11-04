package auth

import (
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	Errorhandel "forum/backend/Errors"
	"forum/backend/database"
	"forum/backend/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandlerGet(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("session_token")
	if err == nil {
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
	var Login models.LoginData
	if err != nil {
		Errorhandel.Errordirect(w, "Internal Server Error ", http.StatusInternalServerError)
		return
	}

	r.ParseForm()
	Login.LoginInput = r.FormValue("username")
	Login.Password = r.FormValue("password")
	Login.LoginInput = strings.TrimSpace(Login.LoginInput)

	if Login.LoginInput == "" || Login.Password == "" {
		tmpt.Execute(w, models.ErrorData{Error: "Username and password cannot be empty"})
		return
	}

	err = database.Db.QueryRow(
		`SELECT id, password_hash , username FROM users WHERE username = ? OR email = ?`,
		Login.LoginInput, Login.LoginInput,
	).Scan(&Login.PasswordHash, &Login.Username)
	if err != nil {
		tmpt.Execute(w, models.ErrorData{Error: "Invalid username or password"})
		log.Println("Login failed (user not found):", err)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(Login.PasswordHash), []byte(Login.Password)); err != nil {
		tmpt.Execute(w, models.ErrorData{Error: "Invalid username or password"})
		return
	}
	_, err = database.Db.Exec("DELETE FROM sessions A WHERE A.username=?", Login.LoginInput)
	if err != nil {
		log.Println("Error deleting old sessions:", err)
	}
	sessionID := uuid.New().String()
	createdAt := time.Now()

	expiresAt := createdAt.Add(24 * time.Hour)

	_, err = database.Db.Exec(`
		INSERT INTO sessions(id, username, created_at, expires_at)
		VALUES (?, ?, ?, ?)
	`, sessionID, Login.Username, createdAt, expiresAt)
	if err != nil {
		tmpt.Execute(w, models.ErrorData{Error: "Could not create session, try again later"})
		log.Println("Error inserting session:", err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionID,
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
