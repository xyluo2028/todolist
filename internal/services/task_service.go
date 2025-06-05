package services

import (
	"errors"
	"fmt"
	"time"
	"todolist/internal/models"
	"todolist/internal/repository"

	"github.com/google/uuid"
)

var DefaultTimestamp = time.Date(2099, 12, 31, 23, 59, 59, 0, time.UTC)

type TaskService struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (svc *TaskService) WriteTask(user, project string, task models.Task) (models.Task, error) {
	if task.Content == "" {
		return models.Task{}, errors.New("task content cannot be empty")
	}
	if task.ID == "" {
		task.ID = fmt.Sprintf("task_%s", uuid.New().String())
	}
	if task.Due.IsZero() {
		task.Due = DefaultTimestamp
	}
	if _, exist := svc.repo.GetTask(user, project, task.ID); !exist {
		err := svc.repo.CreateTask(user, project, task)
		if err != nil {
			return models.Task{}, err
		}
	} else {
		err := svc.repo.UpdateTask(user, project, task)
		if err != nil {
			return models.Task{}, err
		}
	}
	updatedTask, _ := svc.repo.GetTask(user, project, task.ID)
	return updatedTask, nil
}

func (svc *TaskService) MarkTaskComplete(user, project, taskID string) error {
	if _, exist := svc.repo.GetTask(user, project, taskID); !exist {
		return ErrTaskNotFound
	}
	return svc.repo.CompleteTask(user, project, taskID)
}

func (svc *TaskService) GetTasks(user, project string) ([]models.Task, error) {
	return svc.repo.ListTasks(user, project)
}

func (svc *TaskService) GetProjects(user string) ([]string, error) {
	return svc.repo.ListProjects(user)
}

func (svc *TaskService) RemoveProject(user, project string) error {
	return svc.repo.DeleteProject(user, project)
}

func (svc *TaskService) RemoveTask(user, project, taskID string) error {
	if _, exist := svc.repo.GetTask(user, project, taskID); !exist {
		return ErrTaskNotFound
	}
	return svc.repo.DeleteTask(user, project, taskID)
}

func (svc *TaskService) RemoveUserTasks(user string) error {
	return svc.repo.DeleteUserTasks(user)
}
