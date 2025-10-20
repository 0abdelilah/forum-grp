package posts

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"forum/backend/database"
	"forum/backend/home"
)

func LoadPostsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Not used anymore")
}

func SeePostdetail(w http.ResponseWriter, r *http.Request) {
	PostsTemplete, err := template.ParseFiles("./frontend/templates/post-detail.html")
	if err != nil {
		log.Fatal(err)
	}
	if r.URL.Path != "/post-detail" {
		http.Error(w, "Page Not Found ", http.StatusNotFound)
		return
	}
	PostDatacontent := database.AllPageData(r, "postContent")
	// fmt.Println("Data from Post",PostDatacontent)
	PostsTemplete.Execute(w, PostDatacontent)
}

func CreatePostsHandler(w http.ResponseWriter, r *http.Request) {
	var PostData struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	username, err := home.GetUsernameFromCookie(r, "session_token")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"success": "false",
			"error":   "Unauthenticated",
		})
		log.Fatal(err)
		return
	}
	err = json.NewDecoder(r.Body).Decode(&PostData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"success": "false",
			"error":   "Invalid JSON",
		})
		log.Fatal(err)
		return
	}

	if len(PostData.Title) < 1 || len(PostData.Title) > 90 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"success": "false",
			"error":   "Title must be between 1 and 90 characters",
		})
		return
	}

	if len(PostData.Content) < 1 || len(PostData.Content) > 300 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"success": "false",
			"error":   "Content must be between 1 and 300 characters",
		})
		return
	}

	err = InsertPost(username, PostData.Title, PostData.Content)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"success": "false",
			"error":   "Internal Server error, try later",
		})
		log.Fatal(err)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"success": "true",
	})
}

func InsertPost(username, title, content string) error {
	_, err := database.Db.Exec(`
	INSERT INTO posts (username, title, content)
	VALUES (?, ?, ?)
	`, username, title, content)
	if err != nil {
		return err
	}
	return nil
}
