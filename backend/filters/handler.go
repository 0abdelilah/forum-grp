package filters

import (
	"fmt"
	"forum/backend/database"
	"forum/backend/models"
	"html/template"
	"log"
	"net/http"
	"strings"
)

// func showfelters()
func FelterbyCategory(m models.PageData, value string) []models.Post {
	var newPosts []models.Post
	for i := 0; i < len(m.AllPosts); i++ {
		if contains(m.AllPosts[i].Categories, value) {
			newPosts = append(newPosts, m.AllPosts[i])
		}
	}
	return newPosts
}
func contains(elems []models.Category, v string) bool {
	for _, s := range elems {
		if v == s.Category {
			return true
		}
	}
	return false
}

func FilterPostsHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the form
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get selected categories
	categories := r.Form["categories[]"]

	// Example: query posts matching these categories
	var placeholders []string
	var args []interface{}
	for i, cat := range categories {
		placeholders = append(placeholders, "?")
		args = append(args, cat)
		_ = i
	}

	query := "SELECT id, title, content FROM posts"
	if len(categories) > 0 {
		query += " WHERE category IN (" + strings.Join(placeholders, ",") + ")"
	}

	rows, err := database.Db.Query(query, args...)
	if err != nil {
		http.Error(w, "DB query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
		if err := rows.Scan(&p.Id, &p.Title, &p.Content); err != nil {
			http.Error(w, "Failed to scan row", http.StatusInternalServerError)
			return
		}
		posts = append(posts, p)
	}

	tmpl, err := template.ParseFiles("./frontend/templates/index.html")
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	PageData := database.AllPageData(r, "HomeData")
	PageData.AllPosts = posts

	if err := tmpl.Execute(w, PageData); err != nil {
		log.Printf("template execution error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
