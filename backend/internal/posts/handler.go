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
    "reactions": {
        "likes": ["abelkhadir", "iallaoui"],
        "dislikes": ["nabil"]
    },
    "comments": 4,
    "category": "cybersecurity",
    "created_at": "2025-10-13T18:00:00Z"
  },
  {
    "id": 2,
    "title": "Top 5 Tools for Web Pentesting",
    "author": "sokayna",
    "content": "A quick overview of the most effective tools for modern web app pentesting.",
    "reactions": {
        "likes": ["abelkhadir", "iallaoui"],
        "dislikes": ["nabil"]
    },
    "comments": 3,
    "category": "programming",
    "created_at": "2025-10-12T16:30:00Z"
  }
]
`

type Reactions struct {
	Likes    []string `json:"likes"`
	Dislikes []string `json:"dislikes"`
}

type Post struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Content     string    `json:"content"`
	CommentsNum int       `json:"comments"`
	Reactions   Reactions `json:"reactions"`
	Category    string    `json:"category"`
	CreatedAt   string    `json:"created_at"`
}

func LoadPostsHandler(w http.ResponseWriter, r *http.Request) {

	// SELECT id, title, author, content, comments, created_at FROM posts

	var Posts []Post

	// example
	err := json.Unmarshal([]byte(jsonExample), &Posts)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"success": "false",
			"error":   "Invalid JSON",
		})
		return
	}

	// return success true and comments
	json.NewEncoder(w).Encode(map[string]any{
		"success": "true",
		"posts":   Posts,
	})
}

func CreatePostsHandler(w http.ResponseWriter, r *http.Request) {

}
