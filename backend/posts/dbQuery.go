package posts

import "forum/backend/database"

func addPostCategory(postID, categoryID int) error {
	_, err := database.Db.Exec(
		`INSERT OR IGNORE INTO post_categories (post_id, category_id) VALUES (?, ?)`,
		postID, categoryID,
	)
	return err
}

func getCategoryID(name string) (int, error) {
	var id int
	err := database.Db.QueryRow(`SELECT id FROM categories WHERE name = ?`, name).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
