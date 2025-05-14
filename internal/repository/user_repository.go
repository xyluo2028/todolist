package repository

import "todolist/internal/models"

type UserRepository interface {
	AddUser(user models.User) error
	GetUser(username string) (models.User, error)
	Authenticate(username, password string) bool
	DeactivateUser(username string) error
}
