package models

type Post struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Username    string `json:"username"`
	Content     string `json:"content"`
	CommentsNum int    `json:"comments"`
	Likes       int    `json:"likes_count"`
	Dislikes    int    `json:"dislikes_count"`
	Category    string `json:"category"`
	CreatedAt   string `json:"created_at"`
	Categories  string
	Comments    []Comment
}
