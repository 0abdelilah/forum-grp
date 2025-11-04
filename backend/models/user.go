package models

type ErrorData struct {
	Error string
}
type LoginData struct {
	UserID       int
	Username     string
	PasswordHash string
	LoginInput   string
	Password     string
}
