package likes

import (
	"fmt"
	"forum/backend/auth"
	"forum/backend/home"
	"net/http"
)

func LikeHandler(w http.ResponseWriter, r *http.Request) {
	postid := r.URL.Query().Get("postid")
	page := r.URL.Query().Get("page")
	path := "/"
	if page == "postdetails" {
		path = "/post-detail/?postid=" + postid
	}

	username, err := auth.GetUsernameFromCookie(r, "session_token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	err = insertLike(postid, username)
	if err != nil {
		fmt.Println(err)
		home.PostPageError(w, r, "Internal server error, try later")
		return
	}

	http.Redirect(w, r, path, http.StatusSeeOther)
}

func DislikeHandler(w http.ResponseWriter, r *http.Request) {
	postid := r.URL.Query().Get("postid")
	page := r.URL.Query().Get("page")
	path := "/"
	if page == "postdetails" {
		path = "/post-detail/?postid=" + postid
	}

	username, err := auth.GetUsernameFromCookie(r, "session_token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	fmt.Println(postid)
	err = insertDislike(postid, username)
	if err != nil {
		home.PostPageError(w, r, "Internal server error, try later")
		fmt.Println(err)
		return
	}

	http.Redirect(w, r, path, http.StatusSeeOther)
}
