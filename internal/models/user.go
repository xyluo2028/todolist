package models

type User struct {
	Username string `json:"username"`
	Password string `json:"password"` // store hashed passwords in real scenarios
	Active   bool   `json:"active"`
}
