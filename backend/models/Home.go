package models

type PageData struct {
	Username       string
	Categories     []Category
	CategoryChoice []Category
	AllPosts       []Post
	AllPostLikes   []Likes
	AllComments    []Comment
	PostContent    Post
}
