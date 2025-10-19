package models

type PageData struct {
	Username     string
	IsLoggedIn   bool
	Categories   []Category
	AllPosts     []Post
	AllPostLikes []Reactions
	AllComments  []Comment
}
