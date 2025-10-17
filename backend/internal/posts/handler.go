package posts

import (
	"encoding/json"
	"net/http"
)

var jsonExample = `
[
  {
    "id": 1,
    "title": "Getting Started with Bug Bounty",
    "author": "Gadr",
    "content": "Learn how to find your first vulnerability and report it responsibly.",
    "comments": 4,
    "created_at": "2025-10-13T18:00:00Z"
  },
  {
    "id": 2,
    "title": "Top 5 Tools for Web Pentesting",
    "author": "Gadr",
    "content": "A quick overview of the most effective tools for modern web app pentesting.",
    "comments": 3,
	"created_at": "2025-10-12T16:30:00Z"
  }
]
`

type Post struct {
	Id         int
	Title      string
	Author     string
	Content    string
	Comments   int
	Created_at string
}

func LoadPostsHandler(w http.ResponseWriter, r *http.Request) {

	// SELECT id, title, author, content, comments, created_at FROM posts

	var Posts []Post

	// example
	err := json.NewDecoder(r.Body).Decode(&Posts)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"success": "false",
			"error":   "Invalid JSON",
		})
		return
	}

}

func CreatePostsHandler(w http.ResponseWriter, r *http.Request) {

}
