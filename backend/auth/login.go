package auth

import (
	"fmt"
	"forum/backend/database"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandlerGet(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../frontend/templates/login.html")
}

func LoginHandlerPost(w http.ResponseWriter, r *http.Request) {
	fmt.Println("dfs")
	var (
		userID     int
		storedHash string
	)

	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	err := database.Db.QueryRow(
		`SELECT id, password_hash FROM users WHERE username = ?`,
		username,
	).Scan(&userID, &storedHash)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		log.Println("Login failed (user not found):", err)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password)); err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		log.Println("Login failed (wrong password):", err)
		return
	}

	sessionID := uuid.New().String()
	createdAt := time.Now()
	expiresAt := createdAt.Add(24 * time.Hour)

	_, err = database.Db.Exec(`
		INSERT INTO sessions(id, user_id, created_at, expires_at)
		VALUES (?, ?, ?, ?)
	`, sessionID, userID, createdAt, expiresAt)
	if err != nil {
		http.Error(w, "Could not create session", http.StatusInternalServerError)
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
