package likes

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	Errorhandel "forum/backend/Errors"
	"forum/backend/auth"
	"forum/backend/database"
)

// func (r *http.Request, cookieName string) (int, error) {
// 	cookie, err := r.Cookie(cookieName)
// 	if err != nil {
// 		return 0, errors.New("session cookie not found")
// 	}

// 	sessionToken := cookie.Value
// 	if sessionToken == "" {
// 		return 0, errors.New("empty session token")
// 	}

// 	var userID int
// 	err = database.Db.QueryRow(`
//         SELECT user_id FROM sessions WHERE token = ?
//     `, sessionToken).Scan(&userID)

// 	if err == sql.ErrNoRows {
// 		return 0, errors.New("invalid session token")
// 	}
// 	if err != nil {
// 		return 0, err
// 	}

// 	return userID, nil
// }

func HandleLikeOrDislike(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	username, err := auth.GetUsernameFromCookie(r, "session_token")
	if err != nil && err != sql.ErrNoRows && fmt.Sprintf("%v", err) != "http: named cookie not present" {
		Errorhandel.Errordirect(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	fmt.Println("1")
	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	fmt.Println("2")

	targetType := r.FormValue("target_type")
	targetID := r.FormValue("target_id")
	fmt.Println("The post is", targetID)
	value := r.FormValue("value")

	var intValue int
	if value == "1" {
		intValue = 1
	} else {
		intValue = -1
	}
	fmt.Println("3")

	tID, err := strconv.Atoi(targetID)
	if err != nil {
		fmt.Println("err atoi")
		return
	}
	err = toggleLikeDislike(username, targetType, tID, intValue)
	if err != nil {
		fmt.Println("Error handling like/dislike:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
	fmt.Println("4")
}

func toggleLikeDislike(username string, targetType string, targetID int, value int) error {
	var existingValue int
	var query string
	var args []interface{}
	var tableName string
	var likeColumn string

	switch targetType {
	case "post":
		query = "SELECT value FROM likes WHERE username = ? AND post_id = ?"
		args = []interface{}{username, targetID}
		tableName = "posts"
		likeColumn = "post_id"
	case "comment":
		query = "SELECT value FROM likes WHERE username = ? AND comment_id = ?"
		args = []interface{}{username, targetID}
		tableName = "comments"
		likeColumn = "comment_id"
	default:
		return fmt.Errorf("unknown target type: %s", targetType)
	}

	err := database.Db.QueryRow(query, args...).Scan(&existingValue)

	if err == sql.ErrNoRows {
		insertQuery := fmt.Sprintf("INSERT INTO likes (username, %s, value) VALUES (?, ?, ?)", likeColumn)
		_, err = database.Db.Exec(insertQuery, username, targetID, value)
		if err != nil {
			return err
		}

		if value == 1 {
			_, err = database.Db.Exec(fmt.Sprintf("UPDATE %s SET likes_count = likes_count + 1 WHERE id = ?", tableName), targetID)
		} else {
			_, err = database.Db.Exec(fmt.Sprintf("UPDATE %s SET dislikes_count = dislikes_count + 1 WHERE id = ?", tableName), targetID)
		}
		return err
	}

	if err != nil {
		return err
	}

	if existingValue == value {
		deleteQuery := fmt.Sprintf("DELETE FROM likes WHERE username = ? AND %s = ?", likeColumn)
		_, err = database.Db.Exec(deleteQuery, username, targetID)
		if err != nil {
			return err
		}

		if value == 1 {
			_, err = database.Db.Exec(fmt.Sprintf("UPDATE %s SET likes_count = likes_count - 1 WHERE id = ?", tableName), targetID)
		} else {
			_, err = database.Db.Exec(fmt.Sprintf("UPDATE %s SET dislikes_count = dislikes_count - 1 WHERE id = ?", tableName), targetID)
		}
		return err
	}

	updateQuery := fmt.Sprintf("UPDATE likes SET value = ? WHERE username = ? AND %s = ?", likeColumn)
	_, err = database.Db.Exec(updateQuery, value, username, targetID)
	if err != nil {
		return err
	}

	if value == 1 {
		_, err = database.Db.Exec(fmt.Sprintf("UPDATE %s SET likes_count = likes_count + 1, dislikes_count = dislikes_count - 1 WHERE id = ?", tableName), targetID)
	} else {
		_, err = database.Db.Exec(fmt.Sprintf("UPDATE %s SET dislikes_count = dislikes_count + 1, likes_count = likes_count - 1 WHERE id = ?", tableName), targetID)
	}

	return err
}
