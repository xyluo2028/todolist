package handlers

import (
	"encoding/json"
	"fmt"
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
	if err != nil {
		http.Error(w, "Error retrieving tasks", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if len(tasks) == 0 {
		fmt.Fprintln(w, "No tasks found!")
		return
	}

	fmt.Fprintln(w, "Todos: ")
	tasknum := 0
	for _, task := range tasks {
		data, err := json.Marshal(task)
		if err != nil {
			http.Error(w, "Error serializing task", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Task %d: %s \n", tasknum, string(data))
		tasknum++
	}
}

func (h *TaskHandler) GetAllProjectsHttp(w http.ResponseWriter, r *http.Request) {
	user, _, _ := r.BasicAuth()
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

	fmt.Fprintln(w, "Projects: ")
	for _, project := range projects {
		fmt.Fprintf(w, "%s \n", project)
	}
}

func (h *TaskHandler) WriteTaskHttp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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
	err := h.svc.MarkTaskComplete(user, project, key)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error marking task as complete: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "task: %s completed", key)
}

func (h *TaskHandler) RemoveTaskHttp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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
	err := h.svc.RemoveTask(user, project, key)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error removing task: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "task: %s removed", key)
}

func (h *TaskHandler) RemoveProjectHttp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	project := r.URL.Query().Get("pjt")
	user, _, _ := r.BasicAuth()
	if project == "" {
		http.Error(w, "Project query parameter 'pjt' is required", http.StatusBadRequest)
		return
	}
	err := h.svc.RemoveProject(user, project)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error removing project: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "project: %s removed", project)
}
