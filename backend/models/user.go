package models

type ErrorData struct {
	Error string
}
type LoginData struct {
	Email        string
	Username     string
	PasswordHash string
	LoginInput   string
	Password     string
}
