package models

type PageData struct {
	Error       string
	Username    string
	Categories  []Category
	AllPosts    []Post
	PostContent Post
}
