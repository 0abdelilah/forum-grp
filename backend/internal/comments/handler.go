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
	// TODO: verify all feilds exist

	// Get comments from db
	err = database.SaveComment(comment.PostId, comment.UserId, comment.Content)
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

	// Get comments from db
	comments, err := database.GetComments(comment.PostId)
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
