package models

type PageData struct {
	Username     string
	Categories   []Category
	AllPosts     []Post
	AllPostLikes []Likes
	AllComments  []Comment
	PostContent  Post
}
