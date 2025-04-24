package models

type User struct {
	ID    string `json:"gid"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
