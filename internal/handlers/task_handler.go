package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"todolist/internal/models"
	"todolist/internal/services"
)

type TaskHandler struct {
	svc *services.TaskService
}

func NewTaskHandler(svc *services.TaskService) *TaskHandler {
	return &TaskHandler{svc: svc}
}

func (h *TaskHandler) GetAllTasksFromPjtHttp(w http.ResponseWriter, r *http.Request) {
	user, _, _ := r.BasicAuth()
	project := r.URL.Query().Get("pjt")
	tasks, err := h.svc.GetTasks(user, project)
	log.Printf("Retrieving tasks for user '%s', URI = '%s', method = '%s', project = '%s'", user, r.RequestURI, r.Method, project)

	if err != nil {
		http.Error(w, "Error retrieving tasks", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if len(tasks) == 0 {
		fmt.Fprintln(w, "No tasks found!")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		http.Error(w, "Error serializing tasks", http.StatusInternalServerError)
		return
	}
}

func (h *TaskHandler) GetAllProjectsHttp(w http.ResponseWriter, r *http.Request) {
	user, _, _ := r.BasicAuth()
	log.Printf("Retrieving projects for user '%s', URI = '%s', method = '%s'", user, r.RequestURI, r.Method)
	projects, err := h.svc.GetProjects(user)
	if err != nil {
		http.Error(w, "Error retrieving projects", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if len(projects) == 0 {
		fmt.Fprintln(w, "No projects found!")
		return
	}

	for _, project := range projects {
		fmt.Fprintf(w, "%s \n", project)
	}
}

func (h *TaskHandler) CreateProjectHttp(w http.ResponseWriter, r *http.Request) {
	user, _, _ := r.BasicAuth()
	project := r.URL.Query().Get("pjt")
	if project == "" {
		http.Error(w, "Project query parameter 'pjt' is required", http.StatusBadRequest)
		return
	}
	log.Printf("Creating project '%s' for user '%s', URI= '%s', method= '%s'", project, user, r.RequestURI, r.Method)

	if err := h.svc.CreateProject(user, project); err != nil {
		http.Error(w, fmt.Sprintf("Error creating project: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Project '%s' created successfully", project)
}

func (h *TaskHandler) WriteTaskHttp(w http.ResponseWriter, r *http.Request) {
	project := r.URL.Query().Get("pjt")
	if project == "" {
		http.Error(w, "Project query parameter 'pjt' is required", http.StatusBadRequest)
		return
	}
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}
	log.Printf("Writing task for project '%s', URI= '%s', method= '%s', task= '%s'", project, r.RequestURI, r.Method, task)

	if valid, err := h.validTask(task); !valid {
		http.Error(w, fmt.Sprintf("Invalid task: %v", err), http.StatusBadRequest)
		return
	}
	user, _, _ := r.BasicAuth()

	task, err := h.svc.WriteTask(user, project, task)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error writing task: %v", err), http.StatusInternalServerError)
		return
	}
	// Return success message
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Write task with key: %s at: %v", task.ID, task.UpdatedTime)
}

func (h *TaskHandler) CompleteTaskHttp(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	user, _, _ := r.BasicAuth()
	project := r.URL.Query().Get("pjt")
	if project == "" {
		http.Error(w, "Project query parameter 'pjt' is required", http.StatusBadRequest)
		return
	}
	if key == "" {
		http.Error(w, "Key query parameter 'key' is required", http.StatusBadRequest)
		return
	}
	log.Printf("Completing task '%s' for project '%s', URI= '%s', method= '%s'", key, project, r.RequestURI, r.Method)

	err := h.svc.MarkTaskComplete(user, project, key)
	if err != nil {
		if err == services.ErrTaskNotFound {
			http.Error(w, fmt.Sprintf("Task in project %s with key %s not found", project, key), http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Error marking task as complete: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "task: %s completed", key)
}

func (h *TaskHandler) RemoveTaskHttp(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	user, _, _ := r.BasicAuth()
	project := r.URL.Query().Get("pjt")
	if project == "" {
		http.Error(w, "Project query parameter 'pjt' is required", http.StatusBadRequest)
		return
	}
	if key == "" {
		http.Error(w, "Key query parameter 'key' is required", http.StatusBadRequest)
		return
	}
	log.Printf("Removing task '%s' for project '%s', URI= '%s', method= '%s'", key, project, r.RequestURI, r.Method)
	err := h.svc.RemoveTask(user, project, key)
	if err != nil {
		if err == services.ErrTaskNotFound {
			http.Error(w, fmt.Sprintf("Task in project %s with key %s not found", project, key), http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Error removing task: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "task: %s removed", key)
}

func (h *TaskHandler) RemoveProjectHttp(w http.ResponseWriter, r *http.Request) {
	project := r.URL.Query().Get("pjt")
	user, _, _ := r.BasicAuth()
	if project == "" {
		http.Error(w, "Project query parameter 'pjt' is required", http.StatusBadRequest)
		return
	}
	log.Printf("Removing project '%s' for user '%s', URI= '%s', method= '%s'", project, user, r.RequestURI, r.Method)
	err := h.svc.RemoveProject(user, project)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error removing project: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "project: %s removed", project)
}

func (h *TaskHandler) validTask(task models.Task) (bool, error) {
	if task.Content == "" {
		return false, fmt.Errorf("task content cannot be empty")
	}
	if task.Priority < 0 || task.Priority > 10 {
		return false, fmt.Errorf("task priority must be between 0 and 10")
	}
	return true, nil
}
