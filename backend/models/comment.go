package models

type Comment struct {
	Id           string
	PostId       string
	UserId       string
	Content      string
	Created      string
	Likes        string
	CurrUserRate string
}
