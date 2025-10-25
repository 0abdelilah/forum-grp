package comments

import (
	"fmt"
	"forum/backend/home"
	"net/http"
)

func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	postid := r.URL.Query().Get("postid")
	path := "/post-detail/?postid=" + postid

	username, err := home.GetUsernameFromCookie(r, "session_token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	err = r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Failed to parse form")
		return
	}

	postID := r.FormValue("postid")
	content := r.FormValue("content")

	fmt.Println(postID, content, username)

	err = insertComment(postID, content, username)
	if err != nil {
		http.Redirect(w, r, path, http.StatusSeeOther)
		fmt.Println(err)
		return
	}

	http.Redirect(w, r, path, http.StatusSeeOther)
}
