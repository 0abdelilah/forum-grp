
package auth

import (
	"log"
	"net/http"
	"time"

	"forum/backend/database"

	"github.com/google/uuid"
)



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

func DeleteSessionByID(id string) error {
	_, err := database.Db.Exec("DELETE FROM sessions WHERE username = ?", id)
	if err != nil {
		return err
	}
	return nil
}
