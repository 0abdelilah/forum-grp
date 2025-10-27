package auth

import (
	"database/sql"
	"forum/backend/database"
	"net/http"
)

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
