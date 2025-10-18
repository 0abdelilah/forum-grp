package models

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
