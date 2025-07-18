package repository

import "todolist/internal/models"

type TaskRepository interface {
	CreateTask(username, project string, task models.Task) error
	CreateProject(username, project string) error
	ListTasks(username, project string) ([]models.Task, error)
	ListProjects(username string) ([]string, error)
	GetTask(username, project, taskID string) (models.Task, bool)
	CompleteTask(username, project, taskID string) error
	UpdateTask(username, project string, task models.Task) error
	DeleteTask(username, project, taskID string) error
	DeleteProject(username, project string) error
	DeleteUserTasks(username string) error
}
