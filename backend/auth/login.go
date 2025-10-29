package auth

import (
	"html/template"
	"log"
	"net/http"
	"time"

	databasecreate "forum/backend/database"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandlerGet(w http.ResponseWriter, r *http.Request) {
	tmpt, err := template.ParseFiles("./frontend/templates/login.html")
	if err != nil {
		log.Fatal(err)
		return
	}
	tmpt.Execute(w, nil)
}

func LoginHandlerPost(w http.ResponseWriter, r *http.Request) {
	Db := databasecreate.Open()
	tmpt, err := template.ParseFiles("./frontend/templates/login.html")
	if err != nil {
		log.Fatal(err)
		return
	}

	var (
		userID     int
		storedHash string
	)

	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		tmpt.Execute(w, struct{ Error string }{Error: "Username and password are required"})
		return
	}

	err = Db.QueryRow(
		`SELECT id, password_hash FROM users WHERE username = ?`,
		username,
	).Scan(&userID, &storedHash)
	if err != nil {
		tmpt.Execute(w, struct{ Error string }{Error: "Invalid username or password"})
		log.Println("Login failed (user not found):", err)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password)); err != nil {
		tmpt.Execute(w, struct{ Error string }{Error: "Invalid username or password"})
		log.Println("Login failed (wrong password):", err)
		return
	}

	sessionID := uuid.New().String()
	createdAt := time.Now()
	expiresAt := createdAt.Add(24 * time.Hour)

	_, err = Db.Exec(`
		INSERT INTO sessions(id, username, created_at, expires_at)
		VALUES (?, ?, ?, ?)
	`, sessionID, username, createdAt, expiresAt)
	if err != nil {
		tmpt.Execute(w, struct{ Error string }{Error: "Could not create session, try again later"})
		log.Println("Error inserting session:", err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionID,
		Path:    "/",
		Expires: expiresAt,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
