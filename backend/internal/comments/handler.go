package comments

import (
	"encoding/json"
	"fmt"
	"net/http"

	"forum/internal/database"
)

func SaveCommentHandler(w http.ResponseWriter, r *http.Request) {
	var comment database.Comment

	// Get args from the request: postid, author, comment
	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"success": "false",
			"error":   "Invalid JSON",
		})
		return
	}

	cookie, err := r.Cookie("userid")
	userid := cookie.Value
	if err != nil || userid == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"success": "false",
			"error":   "Unauthenticated",
		})
		fmt.Println("Unauthenticated")
		return
	}

	// verify all feilds exist
	if comment.PostId == "" || comment.Content == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"success": "false",
			"error":   "Missing required fields",
		})
		return
	}

	fmt.Println(comment.PostId, userid, comment.Content)

	// Get comments from db
	err = database.SaveComment(comment.PostId, userid, comment.Content)
	if err != nil {
		fmt.Println(err)
		return
	}

	// return success true and comments
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": "true",
	})
}

func GetCommentsHandler(w http.ResponseWriter, r *http.Request) {
	postId := r.URL.Query().Get("postId")

	if postId == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"success": "false",
			"error":   "Invalid JSON",
		})
		return
	}

	// Get comments from db
	comments, err := database.GetComments(postId)
	if err != nil {
		fmt.Println(err)
		return
	}

	// return success true and comments
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":  "true",
		"comments": comments,
	})
}
