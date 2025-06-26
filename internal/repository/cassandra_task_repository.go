// filepath: /home/luoxinyu01/github/todolist/internal/repository/cassandra_task_repository.go
package repository

import (
	"fmt"
	"log"
	"time"
	"todolist/internal/models"

	"github.com/gocql/gocql"
)

type CassandraTaskRepository struct {
	session *gocql.Session
}

func NewCassandraTaskRepository(session *gocql.Session) *CassandraTaskRepository {
	return &CassandraTaskRepository{session: session}
}

func (repo *CassandraTaskRepository) CreateProject(username, project string) error {
	query := "INSERT INTO projects (username, project) VALUES (?, ?)"
	err := repo.session.Query(query, username, project).Exec()
	if err != nil {
		return fmt.Errorf("error creating project %s for user %s: %w", project, username, err)
	}
	return nil
}

func (repo *CassandraTaskRepository) CreateTask(username, project string, task models.Task) error {
	var existing string
	if err := repo.session.Query(
		`SELECT project FROM todolist.projects WHERE project = ? ALLOW FILTERING`,
		project,
	).Scan(&existing); err != nil {
		if err == gocql.ErrNotFound {
			return fmt.Errorf("project %q does not exist", project)
		}
		return fmt.Errorf("failed to verify project existence: %w", err)
	}
	query := "INSERT INTO tasks (username, project, id, content, priority, updated_time, due, completed) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	err := repo.session.Query(query, username, project, task.ID, task.Content, task.Priority, time.Now(), task.Due, task.Completed).Exec()
	return err
}

func (repo *CassandraTaskRepository) ListTasks(username, project string) ([]models.Task, error) {
	var tasks []models.Task
	query := "SELECT id, content, priority, updated_time, due, completed FROM tasks WHERE username = ? AND project = ?"
	iter := repo.session.Query(query, username, project).Iter()
	defer iter.Close()

	var task models.Task
	for iter.Scan(&task.ID, &task.Content, &task.Priority, &task.UpdatedTime, &task.Due, &task.Completed) {
		tasks = append(tasks, task)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (repo *CassandraTaskRepository) UpdateTask(username, project string, task models.Task) error {
	query := "UPDATE tasks SET content = ?, priority = ?, updated_time = ?, due = ?, completed = ? WHERE username = ? AND project = ? AND id = ?"
	err := repo.session.Query(query, task.Content, task.Priority, time.Now(), task.Due, task.Completed, username, project, task.ID).Exec()
	return err
}

func (repo *CassandraTaskRepository) DeleteTask(username, project, taskID string) error {
	query := "DELETE FROM tasks WHERE username = ? AND project = ? AND id = ?"
	err := repo.session.Query(query, username, project, taskID).Exec()
	return err
}

func (repo *CassandraTaskRepository) CompleteTask(username, project, taskID string) error {
	query := "UPDATE tasks SET completed = true WHERE username = ? AND project = ? AND id = ?"
	err := repo.session.Query(query, username, project, taskID).Exec()
	return err
}

func (repo *CassandraTaskRepository) GetTask(username, project, taskID string) (models.Task, bool) {
	var task models.Task
	query := "SELECT id, content, priority, updated_time, due, completed FROM tasks WHERE username = ? AND project = ? AND id = ? ALLOW FILTERING"
	err := repo.session.Query(query, username, project, taskID).Scan(&task.ID, &task.Content, &task.Priority, &task.UpdatedTime, &task.Due, &task.Completed)
	if err != nil {
		if err == gocql.ErrNotFound {
			return task, false
		}
		log.Println("Error retrieving task:", err)
		return task, false
	}
	return task, true
}

func (repo *CassandraTaskRepository) ListProjects(username string) ([]string, error) {
	var projects []string

	query := "SELECT project FROM projects WHERE username = ?"
	iter := repo.session.Query(query, username).Iter()

	var projectName string
	for iter.Scan(&projectName) {
		projects = append(projects, projectName)
	}

	if err := iter.Close(); err != nil {
		log.Printf("Error iterating over projects for user %s: %v", username, err)
		return nil, fmt.Errorf("error listing projects for user %s: %w", username, err)
	}

	return projects, nil
}

func (repo *CassandraTaskRepository) DeleteProject(username, project string) error {
	query := "DELETE FROM tasks WHERE username = ? AND project = ?"
	err := repo.session.Query(query, username, project).Exec()
	if err != nil {
		log.Printf("Error deleting tasks in project %s for user %s: %v", project, username, err)
		return fmt.Errorf("error deleting tasks in project %s for user %s: %w", project, username, err)
	}
	query = "DELETE FROM projects WHERE username = ? AND project = ?"
	err = repo.session.Query(query, username, project).Exec()
	if err != nil {
		log.Printf("Error deleting project entry %s for user %s: %v", project, username, err)
		return fmt.Errorf("error deleting project entry %s for user %s: %w", project, username, err)
	}
	return nil
}

func (repo *CassandraTaskRepository) DeleteUserTasks(username string) error {
	query := "DELETE FROM tasks WHERE username = ?"
	err := repo.session.Query(query, username).Exec()
	return err
}
