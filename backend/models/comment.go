package models

type Comment struct {
	Id        string
	Username  string
	Content   string
	Likes     int
	Dislikes  int
	CreatedAt string
}
