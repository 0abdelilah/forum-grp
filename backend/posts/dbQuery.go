package posts

import databasecreate "forum/backend/database"

func addPostCategory(postID, categoryID int) error {
	Db := databasecreate.Open()
	_, err := Db.Exec(
		`INSERT OR IGNORE INTO post_categories (post_id, category_id) VALUES (?, ?)`,
		postID, categoryID,
	)
	return err
}

func getCategoryID(name string) (int, error) {
	Db := databasecreate.Open()
	var id int
	err := Db.QueryRow(`SELECT id FROM categories WHERE name = ?`, name).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
