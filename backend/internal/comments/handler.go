package comments

import (
	"forum/database"
	"net/http"
)

// GetUserIDFromCookies retrieves the user_id associated with the session_token cookie.
func GetUserIDFromCookies(r *http.Request) (int, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return 0, err // no cookie or other error
	}

	var userID int
	err = database.Db.QueryRow(`SELECT user_id FROM sessions WHERE id = ? AND expires_at > datetime('now')`, cookie.Value).Scan(&userID)
	if err != nil {
		return 0, err // invalid or expired session
	}

	return userID, nil
}

// func SaveCommentHandler(w http.ResponseWriter, r *http.Request) {
// 	var comment database.Comment

// 	// Get args from the request: postid, author, comment
// 	err := json.NewDecoder(r.Body).Decode(&comment)
// 	if err != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		json.NewEncoder(w).Encode(map[string]string{
// 			"success": "false",
// 			"error":   "Invalid JSON",
// 		})
// 		return
// 	}

// 	// add logic
// 	_, err = GetUserIDFromCookies(r)
// 	if err != nil {
// 		w.WriteHeader(http.StatusUnauthorized)
// 		json.NewEncoder(w).Encode(map[string]string{
// 			"success": "false",
// 			"error":   "Unauthenticated",
// 		})
// 		fmt.Println("Unauthenticated")
// 		return
// 	}

// 	// verify all feilds exist
// 	if comment.PostId == "" || comment.Content == "" {
// 		w.WriteHeader(http.StatusBadRequest)
// 		json.NewEncoder(w).Encode(map[string]string{
// 			"success": "false",
// 			"error":   "Missing required fields",
// 		})
// 		return
// 	}

// 	// Get comments from db
// 	err = database.SaveComment(comment.PostId, "UserId", comment.Content)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	// return success true and comments
// 	json.NewEncoder(w).Encode(map[string]any{
// 		"success": "true",
// 	})
// }

// func GetCommentsHandler(w http.ResponseWriter, r *http.Request) {
// 	postId := r.URL.Query().Get("postId")

// 	if postId == "" {
// 		w.WriteHeader(http.StatusBadRequest)
// 		json.NewEncoder(w).Encode(map[string]string{
// 			"success": "false",
// 			"error":   "Invalid JSON",
// 		})
// 		return
// 	}

// 	// Get comments from db
// 	comments, err := database.GetComments(postId)
// 	if err != nil {
// 		fmt.Println(err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		json.NewEncoder(w).Encode(map[string]string{
// 			"success": "false",
// 			"error":   "Database error",
// 		})
// 		return
// 	}

// 	// return success true and comments
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(map[string]interface{}{
// 		"success":  "true",
// 		"comments": comments,
// 	})
// }
