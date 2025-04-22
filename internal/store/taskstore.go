package store

import (
	"sync"
	. "todolist/internal/models"
)

type TaskStore struct {
	sync.RWMutex
	// Structure: projects[user][project][taskKey] = TaskRecord
	projects map[string]map[string]map[string]TaskRecord
}

func NewTaskStore() *TaskStore {
	return &TaskStore{
		projects: make(map[string]map[string]map[string]TaskRecord),
	}
}

func (s *TaskStore) ensureUserMap(user string) {
	if _, exists := s.projects[user]; !exists {
		s.projects[user] = make(map[string]map[string]TaskRecord)
	}
}

func (s *TaskStore) ensureProjectMap(user, project string) {
	s.ensureUserMap(user)
	if _, exists := s.projects[user][project]; !exists {
		s.projects[user][project] = make(map[string]TaskRecord)
	}
}

func (s *TaskStore) GetAllProjects(user string) []string {
	s.RLock()
	defer s.RUnlock()
	s.ensureUserMap(user)
	projects := make([]string, 0, len(s.projects[user]))
	for project := range s.projects[user] {
		projects = append(projects, project)
	}
	return projects
}

func (s *TaskStore) GetAllTasks(user, project string) map[string]TaskRecord {
	s.RLock()
	defer s.RUnlock()
	s.ensureUserMap(user)
	s.ensureProjectMap(user, project)
	tasksCpy := make(map[string]TaskRecord)
	for k, v := range s.projects[user][project] {
		tasksCpy[k] = v
	}
	return tasksCpy
}

func (s *TaskStore) RemoveTask(user, project, key string) {
	s.Lock()
	defer s.Unlock()
	s.ensureUserMap(user)
	delete(s.projects[user][project], key)
}

func (s *TaskStore) RemoveProject(user, project string) {
	s.Lock()
	defer s.Unlock()
	s.ensureUserMap(user)
	delete(s.projects[user], project)
}

func (s *TaskStore) RemoveUserTasks(user string) {
	s.Lock()
	defer s.Unlock()
	delete(s.projects, user)
}

func (s *TaskStore) GetTask(user, project, key string) (TaskRecord, bool) {
	s.RLock()
	defer s.RUnlock()
	s.ensureUserMap(user)
	task, exists := s.projects[user][project][key]
	return task, exists
}

func (s *TaskStore) HasTask(user, project, key string) bool {
	s.RLock()
	defer s.RUnlock()
	s.ensureUserMap(user)
	_, exists := s.projects[user][project][key]
	return exists
}

func (s *TaskStore) AddTask(user, project, key string, task TaskRecord) bool {
	s.Lock()
	defer s.Unlock()
	s.ensureUserMap(user)
	s.ensureProjectMap(user, project)
	if _, exists := s.projects[user][key]; exists {
		return false

	}
	s.projects[user][project][key] = task
	return true
}

func (s *TaskStore) UpdateTask(user, project, key string, task TaskRecord) bool {
	s.Lock()
	defer s.Unlock()
	s.ensureUserMap(user)
	s.ensureProjectMap(user, key)
	// Check if the task exists
	// If it doesn't exist, return false
	_, exists := s.projects[user][project][key]
	if !exists {
		return false
	}
	s.projects[user][project][key] = task
	return true
}
