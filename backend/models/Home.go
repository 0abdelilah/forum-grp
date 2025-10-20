package models

type PageData struct {
	Username       string
	IsLoggedIn     bool
	Categories     []Category
	CategoryChoice []Category
	AllPosts       []Post
	AllPostLikes   []Reactions
	AllComments    []Comment
	Postcontent    Post
}
