package models

type PageData struct {
	Username       string
	IsLoggedIn     bool
	Categories     []Category
	CategoryChoice []Category
	AllPosts       []Post
	AllPostLikes   []Likes
	AllComments    []Comment
	Postcontent    Post
}
