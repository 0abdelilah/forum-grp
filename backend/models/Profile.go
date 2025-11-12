package models

type Profile struct {
	Username   string
	AllPosts   []Post
	Categories []Category
}
