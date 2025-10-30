package likes

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/backend/database"
	"net/http"
)

func GetUserIDFromCookie(r *http.Request, cookieName string) (int, error) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return 0, errors.New("session cookie not found")
	}

	sessionToken := cookie.Value
	if sessionToken == "" {
		return 0, errors.New("empty session token")
	}

	var userID int
	err = database.Db.QueryRow(`
        SELECT user_id FROM sessions WHERE token = ?
    `, sessionToken).Scan(&userID)

	if err == sql.ErrNoRows {
		return 0, errors.New("invalid session token")
	}
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func HandleLikeOrDislike(w http.ResponseWriter, r *http.Request) {
	  if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
	userID, err := GetUserIDFromCookie(r, "session_token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	targetType := r.FormValue("target_type") // "post" أو "comment"
	targetID := r.FormValue("target_id")
	value := r.FormValue("value")

	var intValue int
	if value == "1" {
		intValue = 1
	} else {
		intValue = -1
	}

	err = toggleLikeDislike(userID, targetType, targetID, intValue)
	if err != nil {
		fmt.Println("Error handling like/dislike:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
}
func toggleLikeDislike(userID int, targetType, targetID string, value int) error {
	var existingValue int
	err := database.Db.QueryRow(`
        SELECT value FROM likes
        WHERE user_id = ? AND target_type = ? AND target_id = ?`,
		userID, targetType, targetID,
	).Scan(&existingValue)

	tableName := ""
	switch targetType {
	case "post":
		tableName = "posts"
	case "comment":
		tableName = "comments"
	default:
		return fmt.Errorf("unknown target type: %s", targetType)
	}

	if err == sql.ErrNoRows {
		// المستخدم ما دارش reaction من قبل → نضيفه
		_, err = database.Db.Exec(`
            INSERT INTO likes (user_id, target_type, target_id, value)
            VALUES (?, ?, ?, ?)`,
			userID, targetType, targetID, value)
		if err != nil {
			return err
		}

		// تحديث العدادات مباشرة
		if value == 1 {
			_, err = database.Db.Exec(fmt.Sprintf(
				"UPDATE %s SET likes_count = likes_count + 1 WHERE id = ?", tableName), targetID)
		} else {
			_, err = database.Db.Exec(fmt.Sprintf(
				"UPDATE %s SET dislikes_count = dislikes_count + 1 WHERE id = ?", tableName), targetID)
		}
		return err
	}

	if err != nil {
		return err
	}

	if existingValue == value {
		// نفس القيمة → نحيد reaction
		_, err = database.Db.Exec(`
            DELETE FROM likes
            WHERE user_id = ? AND target_type = ? AND target_id = ?`,
			userID, targetType, targetID)
		if err != nil {
			return err
		}

		// تحديث العدادات مباشرة
		if value == 1 {
			_, err = database.Db.Exec(fmt.Sprintf(
				"UPDATE %s SET likes_count = likes_count - 1 WHERE id = ?", tableName), targetID)
		} else {
			_, err = database.Db.Exec(fmt.Sprintf(
				"UPDATE %s SET dislikes_count = dislikes_count - 1 WHERE id = ?", tableName), targetID)
		}
		return err
	}

	// كانت مختلفة → نبدلها (من like إلى dislike أو العكس)
	_, err = database.Db.Exec(`
        UPDATE likes
        SET value = ?
        WHERE user_id = ? AND target_type = ? AND target_id = ?`,
		value, userID, targetType, targetID)
	if err != nil {
		return err
	}

	// تحديث العدادات مباشرة
	if value == 1 {
		_, err = database.Db.Exec(fmt.Sprintf(
			"UPDATE %s SET likes_count = likes_count + 1, dislikes_count = dislikes_count - 1 WHERE id = ?", tableName), targetID)
	} else {
		_, err = database.Db.Exec(fmt.Sprintf(
			"UPDATE %s SET dislikes_count = dislikes_count + 1, likes_count = likes_count - 1 WHERE id = ?", tableName), targetID)
	}
	return err
}
