package comments

import (
	"database/sql"
	"fmt"
	"net/http"

	Errorhandel "forum/backend/Errors"
	"forum/backend/auth"
	"forum/backend/posts"
)

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	postid := r.URL.Query().Get("postid")
	path := "/post-detail/?postid=" + postid

	username, ErroFromcookie := auth.GetUsernameFromCookie(r, "session_token")
	if ErroFromcookie.Error != nil&&ErroFromcookie.Error != sql.ErrNoRows && fmt.Sprintf("%v", ErroFromcookie.Error) != "http: named cookie not present" {
		Errorhandel.Errordirect(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err := r.ParseForm()
	if err != nil {
		posts.PostPageError(w, r, "Failed to parse form")
		return
	}

	postID := r.FormValue("postid")
	content := r.FormValue("content")

	if len(content) < 1 || len(content) > 300 {
		posts.PostPageError(w, r, "Comment must be between 1 and 300 characters")
		return
	}

	err = insertComment(postID, username, content)
	if err != nil {
		posts.PostPageError(w, r, "Internal server error, try later")
		fmt.Println(err)
		return
	}

	http.Redirect(w, r, path, http.StatusSeeOther)
}
