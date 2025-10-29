package posts

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"forum/backend/database"
	"forum/backend/home"
)

func SeePostdetail(w http.ResponseWriter, r *http.Request) {
	postStr := r.URL.Query().Get("postid")
	n, err := strconv.Atoi(postStr)
	if n == 0 || err != nil {
		home.PageNotFound(w)
	}
	PostsTemplete, err := template.ParseFiles("./frontend/templates/post-detail.html")

	if err != nil {
		home.PageNotFound(w)
		fmt.Println(err)
	}

	PageData := database.AllPageData(r, "postContent")
	if PageData.PostContent.Id == 0 {
		home.PageNotFound(w)
		return
	}

	PageData.Username, _ = home.GetUsernameFromCookie(r, "session_token")

	PostsTemplete.Execute(w, PageData)
}

func verifyCategories(categories []string) (bool, error) {
	if len(categories) == 0 {
		return false, fmt.Errorf("please choose at least one category")
	}

	defaults := []string{"All", "Programming", "Cybersecurity", "Gadgets & Hardware", "Web Development"}

	for _, cat := range categories {
		valid := false
		for _, allowed := range defaults {
			if cat == allowed {
				valid = true
				break
			}
		}
		if !valid {
			return false, fmt.Errorf("invalid category: %s", cat)
		}
	}

	return true, nil
}

func CreatePostsHandler(w http.ResponseWriter, r *http.Request) {
	username, err := home.GetUsernameFromCookie(r, "session_token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		fmt.Println("Coulndt get user", err)
		return
	}

	err = r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Failed to parse form")
		return
	}

	categories := r.Form["categories"]
	title := r.FormValue("title")
	content := r.FormValue("content")

	fmt.Println(categories)

	if len(title) < 1 || len(content) > 90 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"success": "false",
			"error":   "Title must be between 1 and 90 characters",
		})
		return
	}

	if len(content) < 1 || len(content) > 300 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"success": "false",
			"error":   "Content must be between 1 and 300 characters",
		})
		return
	}

	// test if working
	if valid, err := verifyCategories(categories); !valid {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"success": "false",
			"error":   err,
		})
		return
	}

	err = InsertPost(username, title, content, categories)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"success": "false",
			"error":   "Internal Server error, try later",
		})
		fmt.Println(err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

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

func InsertPost(username, title, content string, categories []string) error {
	res, err := database.Db.Exec(`
		INSERT INTO posts (username, title, content)
		VALUES (?, ?, ?)
	`, username, title, content)
	if err != nil {
		return err
	}

	postID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	// Add all categories
	for _, catName := range categories {
		catID, err := getCategoryID(catName)
		if err != nil {
			fmt.Println("Error getting category ID:", err)
			continue
		}

		fmt.Println("Post ID:", postID, "Category ID:", catID)

		err = addPostCategory(int(postID), catID) // working
		if err != nil {
			fmt.Println("Error adding category to post:", err)
		}
	}

	return nil
}
