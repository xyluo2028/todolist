package services

import (
	"errors"
	"todolist/internal/models"
	"todolist/internal/repository"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (svc *UserService) RegisterUser(username, password string) error {
	if username == "" || password == "" {
		return errors.New("username and password cannot be empty")
	}

	user := models.User{
		Username: username,
		Password: password, // hash this in real-world scenarios
		Active:   true,
	}

	return svc.repo.AddUser(user)
}

func (svc *UserService) AuthenticateUser(username, password string) bool {
	return svc.repo.Authenticate(username, password)
}

func (svc *UserService) DeactivateUser(username string, taskSvc *TaskService) error {
	if err := svc.repo.DeactivateUser(username); err != nil {
		return err
	}
	// Remove user's tasks clearly
	return taskSvc.RemoveUserTasks(username)
}

func (svc *UserService) GetUser(username string) (models.User, error) {
	return svc.repo.GetUser(username)
}
