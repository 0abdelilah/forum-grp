package auth

import (
	"net/http"
	"time"

	"forum/backend/database"
)

func GetUsernameFromCookie(r *http.Request, cookieName string) (string, error) {
	c, err := r.Cookie(cookieName)
	if err != nil {
		return "", err
	}

	var username string
	var expiresAt time.Time

	err = database.Db.QueryRow(
		"SELECT username, expires_at FROM sessions WHERE id = ?",
		c.Value,
	).Scan(&username, &expiresAt)
	if err != nil {
		return "", err
	}

	if expiresAt.Before(time.Now()) {
		err := DeleteSessionByID(c.Value)
		if err != nil {
			return "", err
		}
		return "", nil
	}

	return username, nil
}
