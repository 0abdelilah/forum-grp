package auth

import (
	"net/http"
	"time"

	Errorhandel "forum/backend/Errors"
	"forum/backend/database"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		Errorhandel.Errordirect(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Remove session from DB using the cookie value
	_, err = database.Db.Exec("DELETE FROM sessions WHERE id = ?", cookie.Value)
	if err != nil {
		http.Error(w, "Failed to delete session", http.StatusInternalServerError)
		return
	}

	// Delete the cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
