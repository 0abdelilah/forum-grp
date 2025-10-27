package database

import (
	"database/sql"
	"forum/backend/models"
)

func GetAlllike(db *sql.DB, target string, userID int) ([]models.Likes, error) {
	rows, err := db.Query(`
        SELECT id, user_id, target_id, value
        FROM likes
        WHERE user_id = ?
		WHERE target_type = ?
        ORDER BY created_at ASC
    `, userID, target)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var likes []models.Likes
	for rows.Next() {
		var l models.Likes
		if err := rows.Scan(&l.UserId, &l.Target_id, &l.Target_type, &l.Value); err != nil {
			return nil, err

		}
		likes = append(likes, l)

	}
	return likes, nil

}

func GetAllliketarget(db *sql.DB, Target_id int) ([]models.LikesID, error) {

	rows, err := db.Query(`
        SELECT user_id, target_id, value,id
        FROM likes
      WHERE target_id = ?
        ORDER BY created_at ASC
    `, Target_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var likes []models.LikesID
	for rows.Next() {
		var l models.LikesID
		if err := rows.Scan(&l.UserId, &l.Target_id, &l.Value, &l.Id); err != nil {
			return nil, err

		}
		likes = append(likes, l)

	}
	return likes, nil

}
