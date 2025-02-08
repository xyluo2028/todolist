package store

import "sync"

type TaskStore struct {
	sync.RWMutex
	tasks map[string]map[string]string
}

func NewTaskStore() *TaskStore {
	return &TaskStore{
		tasks: make(map[string]map[string]string),
	}
}

func (s *TaskStore) ensureUserMap(user string) {
	if _, exists := s.tasks[user]; !exists {
		s.tasks[user] = make(map[string]string)
	}
}

func (s *TaskStore) GetAllTasks(user string) map[string]string {
	s.RLock()
	defer s.RUnlock()
	s.ensureUserMap(user)
	tasksCpy := make(map[string]string)
	for k, v := range s.tasks[user] {
		tasksCpy[k] = v
	}
	return tasksCpy
}

func (s *TaskStore) RemoveTask(user, key string) {
	s.Lock()
	defer s.Unlock()
	s.ensureUserMap(user)
	delete(s.tasks[user], key)
}

func (s *TaskStore) RemoveUserTasks(user string) {
	s.Lock()
	defer s.Unlock()
	delete(s.tasks, user)
}

func (s *TaskStore) GetTask(user, key string) (string, bool) {
	s.RLock()
	defer s.RUnlock()
	s.ensureUserMap(user)
	task, exists := s.tasks[user][key]
	return task, exists
}

func (s *TaskStore) HasTask(user, key string) bool {
	s.RLock()
	defer s.RUnlock()
	s.ensureUserMap(user)
	_, exists := s.tasks[user][key]
	return exists
}

func (s *TaskStore) AddTask(user string, key string, task string) bool {
	s.Lock()
	defer s.Unlock()
	s.ensureUserMap(user)
	if _, exists := s.tasks[user][key]; exists {
		return false

	}
	s.tasks[user][key] = task
	return true
}

func (s *TaskStore) UpdateTask(user string, key string, task string) bool {
	s.Lock()
	defer s.Unlock()
	s.ensureUserMap(user)
	_, exists := s.tasks[user][key]
	if !exists {
		return false
	}
	s.tasks[user][key] = task
	return true
}
