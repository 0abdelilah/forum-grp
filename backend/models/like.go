package models

type Likes struct {
	UserId      int
	Target_type string
	Target_id   int
	Value       int
}
type LikesID struct {
	UserId    int
	Value     int
	Target_id string
	Id        int
}
