package auth

import (
	"log"
	"net/http"
	"time"

	"forum/backend/database"

	"github.com/google/uuid"
)

func NotFakeSession(str string) bool {
	var CHECK string
	err := database.Db.QueryRow(`SELECT username FROM sessions WHERE id = ?`, str).Scan(&CHECK)
	if err != nil {
		return false
	}
	return CHECK != ""
}

func CreateSession(username string, w http.ResponseWriter) error {
	// Delete old sessions
	_, err := database.Db.Exec("DELETE FROM sessions WHERE username = ?", username)
	if err != nil {
		log.Println("Error deleting old sessions:", err)
	}

	// Create new session
	sessionID := uuid.New().String()
	createdAt := time.Now()
	expiresAt := createdAt.Add(24 * time.Hour)
	_, err = database.Db.Exec(
		`INSERT INTO sessions(id, username, created_at, expires_at) VALUES (?, ?, ?, ?)`,
		sessionID, username, createdAt, expiresAt,
	)
	if err != nil {
		return err
	}

	// Set cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionID,
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
	})

	return nil
}
