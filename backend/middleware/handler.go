package middleware

import (
	"database/sql"
	"errors"
	"net/http"

	Errorhandel "forum/backend/Errors"
	"forum/backend/auth"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, err := auth.GetUsernameFromCookie(r, "session_token")
		if err != nil {
			if err == errors.New("http: named cookie not present") || err == sql.ErrNoRows {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			} else if username != "" {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			Errorhandel.Errordirect(w, "InternalServerError", http.StatusInternalServerError)
		}

		next.ServeHTTP(w, r)
	})
}
