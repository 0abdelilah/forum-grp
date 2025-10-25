package likes

import (
	"fmt"
	"forum/backend/home"
	"net/http"
)

func LikeHandler(w http.ResponseWriter, r *http.Request) {
	postid := r.URL.Query().Get("postid")
	currentPath := "/post-detail/?postid=" + postid
	// TODO (abdelilah): use another parameter "?page=postdetails or ?page=homepage" for redirecting

	username, err := home.GetUsernameFromCookie(r, "session_token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	err = insertLike(postid, username)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		fmt.Println(err)
		return
	}

	http.Redirect(w, r, currentPath, http.StatusSeeOther)
}

func DislikeHandler(w http.ResponseWriter, r *http.Request) {
	postid := r.URL.Query().Get("postid")

	username, err := home.GetUsernameFromCookie(r, "session_token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	fmt.Println(postid)
	err = insertDislike(postid, username)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		fmt.Println(err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
