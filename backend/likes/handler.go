package likes

import (
	"encoding/json"
	"fmt"
	"forum/backend/home"
	"net/http"
)

func LikeHandler(w http.ResponseWriter, r *http.Request) {

	postid := r.URL.Query().Get("postid")

	username, err := home.GetUsernameFromCookie(r, "session_token")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"success": "false",
			"error":   "Unauthenticated",
		})
		fmt.Println(err)
		return
	}

	insertLike(postid, username)

}

func DislikeHandler(w http.ResponseWriter, r *http.Request) {

	//id := r.URL.Query().Get("postid")

}
