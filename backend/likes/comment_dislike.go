package likes

import (
	"database/sql"
	"fmt"
	"net/http"

	"forum/backend/auth"
	"forum/backend/database"
	"forum/backend/home"
)

func AddCommentDislikeHandler(w http.ResponseWriter, r *http.Request) {
	username, err := auth.GetUsernameFromCookie(r, "session_token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	err = r.ParseForm()
	if err != nil {
		home.PostPageError(w, r, "Failed to parse form")
		return
	}

	postid := r.FormValue("postid")
	cmntid := r.FormValue("cmntid")

	path := "/post-detail/?postid=" + postid

	err = addCmntDislike(cmntid, username)
	if err != nil {
		home.PostPageError(w, r, "Internal server error, try later")
		fmt.Println(err)
		return
	}

	http.Redirect(w, r, path, http.StatusSeeOther)
}

func addCmntDislike(cmntId, username string) error {
	var exists int
	fmt.Println(cmntId, username)
	err := database.Db.QueryRow(
		`SELECT 1 FROM comment_dislikes WHERE comment_id = ? AND username = ?`,
		cmntId, username,
	).Scan(&exists)

	// If user already liked → remove it
	if err == nil {
		_, err = database.Db.Exec(
			`DELETE FROM comment_dislikes WHERE comment_id = ? AND username = ?`,
			cmntId, username,
		)
		if err != nil {
			return err
		}
		_, err = database.Db.Exec(
			`UPDATE comments SET dislikes_count = MAX(dislikes_count - 1, 0) WHERE id = ?`,
			cmntId,
		)
		return err
	}

	// If other error (not just "no rows")
	if err != sql.ErrNoRows {
		return err
	}

	// Add new like
	_, err = database.Db.Exec(
		`INSERT INTO comment_dislikes (username, comment_id) VALUES (?, ?)`,
		username, cmntId,
	)
	if err != nil {
		return err
	}

	_, err = database.Db.Exec(
		`UPDATE comments SET dislikes_count = dislikes_count + 1 WHERE id = ?`,
		cmntId,
	)
	return err
}
