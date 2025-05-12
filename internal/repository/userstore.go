package repository

import "sync"

type UserStore struct {
	sync.RWMutex
	users map[string]string // username -> password
}

func NewUserStore() *UserStore {
	return &UserStore{
		users: map[string]string{
			"admin": "admin123",
		},
	}
}

func (s *UserStore) AddUser(username, password string) bool {
	s.Lock()
	defer s.Unlock()

	if _, exists := s.users[username]; exists {
		return false
	}

	s.users[username] = password
	return true
}

func (s *UserStore) Authenticate(username, password string) bool {
	s.RLock()
	defer s.RUnlock()
	pwd, exists := s.users[username]
	return exists && pwd == password
}

func (s *UserStore) RemoveUser(username string) bool {
	s.Lock()
	defer s.Unlock()

	if _, exists := s.users[username]; !exists {
		return false
	}
	delete(s.users, username)
	return true
}
