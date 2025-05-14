package repository

import (
	"errors"
	"sync"
	"todolist/internal/models"
)

type InMemUserRepository struct {
	mu    sync.RWMutex
	users map[string]models.User
}

func NewInMemUserRepository() *InMemUserRepository {
	return &InMemUserRepository{
		users: make(map[string]models.User),
	}
}

func (repo *InMemUserRepository) AddUser(user models.User) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.users[user.Username]; exists {
		return errors.New("user already exists")
	}

	repo.users[user.Username] = user
	return nil
}

func (repo *InMemUserRepository) GetUser(username string) (models.User, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	user, exists := repo.users[username]
	if !exists {
		return models.User{}, errors.New("user not found")
	}

	return user, nil
}

func (repo *InMemUserRepository) Authenticate(username, password string) bool {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	user, exists := repo.users[username]
	return exists && user.Password == password && user.Active
}

func (repo *InMemUserRepository) DeactivateUser(username string) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	user, exists := repo.users[username]
	if !exists {
		return errors.New("user not found")
	}

	user.Active = false
	repo.users[username] = user
	return nil
}
