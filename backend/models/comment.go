package models

type Comment struct {
	Id           string
	Username     string
	PostId       string
	UserId       string
	Content      string
	CreatedAt    string
	Likes        string
	CurrUserRate string
}
