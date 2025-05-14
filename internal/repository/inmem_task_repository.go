package repository

import (
	"errors"
	"fmt"
	"sync"
	"time"
	"todolist/internal/models"
)

type InMemTaskRepository struct {
	mu    sync.RWMutex
	tasks map[string]map[string]map[string]models.Task
}

func NewInMemTaskRepository() *InMemTaskRepository {
	return &InMemTaskRepository{
		tasks: make(map[string]map[string]map[string]models.Task),
	}
}

func (repo *InMemTaskRepository) CreateTask(username, project string, task models.Task) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if _, exists := repo.tasks[username]; !exists {
		repo.tasks[username] = make(map[string]map[string]models.Task)
	}

	if _, exists := repo.tasks[username][project]; !exists {
		repo.tasks[username][project] = make(map[string]models.Task)
	}
	task.UpdatedTime = time.Now()
	repo.tasks[username][project][task.ID] = task

	return nil
}

func (repo *InMemTaskRepository) ListProjects(username string) ([]string, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	if _, exists := repo.tasks[username]; !exists {
		fmt.Println("User not found")
		return []string{}, nil
	}
	projects := make([]string, 0, len(repo.tasks[username]))
	for project := range repo.tasks[username] {
		projects = append(projects, project)
	}
	return projects, nil
}

func (repo *InMemTaskRepository) ListTasks(username, project string) ([]models.Task, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()
	taskMap, exists := repo.tasks[username][project]
	if !exists {
		return []models.Task{}, nil
	}
	tasks := make([]models.Task, 0, len(taskMap))
	for _, task := range taskMap {
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (repo *InMemTaskRepository) UpdateTask(username, project string, task models.Task) error {
	repo.mu.RLock()
	defer repo.mu.RUnlock()
	// Check if the task exists
	if _, exists := repo.tasks[username][project][task.ID]; !exists {
		return errors.New("task not found")
	}
	task.UpdatedTime = time.Now()
	repo.tasks[username][project][task.ID] = task
	return nil
}

func (repo *InMemTaskRepository) DeleteTask(username, project, taskID string) error {
	repo.mu.RLock()
	defer repo.mu.RUnlock()
	delete(repo.tasks[username][project], taskID)
	return nil
}

func (repo *InMemTaskRepository) DeleteProject(username, project string) error {
	repo.mu.RLock()
	defer repo.mu.RUnlock()
	delete(repo.tasks[username], project)
	return nil
}

func (repo *InMemTaskRepository) DeleteUserTasks(username string) error {
	repo.mu.RLock()
	defer repo.mu.RUnlock()
	delete(repo.tasks, username)
	return nil
}

func (repo *InMemTaskRepository) GetTask(username, project, taskID string) (models.Task, bool) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()
	task, exists := repo.tasks[username][project][taskID]
	return task, exists
}

func (repo *InMemTaskRepository) CompleteTask(username, project, taskID string) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	task, exists := repo.tasks[username][project][taskID]
	if !exists {
		return errors.New("task not found")
	}

	task.Completed = true
	repo.tasks[username][project][taskID] = task
	return nil
}
