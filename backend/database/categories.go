package databasecreate

import (
	"log"
	"strings"

	"forum/backend/models"
)

func GetAllCategories() []models.Category {
	var categories []models.Category

	Db := Open()

	rows, err := Db.Query(`SELECT id, name FROM categories`)
	if err != nil {
		log.Printf("failed to query categories: %v", err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.Id, &c.Category); err != nil {
			log.Printf("error scanning category: %v", err)
			continue
		}
		categories = append(categories, c)
	}

	if err := rows.Err(); err != nil {
		log.Printf("row iteration error: %v", err)
	}
	return categories
}

func getPostCategories(postID int) (string, error) {
	var categories []string
	Db := Open()
	rows, err := Db.Query(`
		SELECT c.name
		FROM categories c
		JOIN post_categories pc ON c.id = pc.category_id
		WHERE pc.post_id = ?
	`, postID)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return "", err
		}
		categories = append(categories, name)
	}

	return strings.Join(categories, ", "), nil
}
