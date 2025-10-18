package models

type PageData struct {
	Categories     []Category
	AllPosts       []Post
	AllPostLikes   []Reactions
	AllComments    []Comment
}
