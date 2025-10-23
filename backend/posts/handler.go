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
		fmt.Println(err)
		return
	}
	err = json.NewDecoder(r.Body).Decode(&PostData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"success": "false",
			"error":   "Invalid JSON",
		})
		fmt.Println(err)
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
		fmt.Println(err)
		return
	}

	home.HomeHandler(w, r)
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
