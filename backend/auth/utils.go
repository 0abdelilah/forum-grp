package auth

import (
	"database/sql"
	"net/http"

	"forum/backend/database"
)
	type ErroFromcookie struct{
		Error error
		ErrorType string
	}
func GetUsernameFromCookie(r *http.Request, cookie_name string) (string, ErroFromcookie) {
	var ErroFromcookie ErroFromcookie

	c, err := r.Cookie(cookie_name)
	if err != nil {
		ErroFromcookie.Error=err
		ErroFromcookie.ErrorType="cookie"
		return "", ErroFromcookie
	}
	var username string
	err = database.Db.QueryRow("SELECT username FROM sessions WHERE id = ? AND expires_at > datetime('now')", c.Value).Scan(&username)
	if err == sql.ErrNoRows {
		ErroFromcookie.Error=err
		ErroFromcookie.ErrorType="Sql"
		return "", ErroFromcookie //
	}
	if err != nil {
		ErroFromcookie.Error=err
		ErroFromcookie.ErrorType="Sqlinternal"
		return "", ErroFromcookie
	}
	return username, ErroFromcookie
}
