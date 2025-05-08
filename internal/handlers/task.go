package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	. "todolist/internal/models"
	. "todolist/internal/store"

	"github.com/google/uuid"
)

type TaskHandler struct {
	store *TaskStore
}

func NewTaskHandler(store *TaskStore) *TaskHandler {
	return &TaskHandler{
		store: store,
	}
}

func (h *TaskHandler) GetAllTasksFromPjtHttp(w http.ResponseWriter, r *http.Request) {
	user, _, _ := r.BasicAuth()
	project := r.URL.Query().Get("pjt")
	tasks := h.store.GetAllTasks(user, project)
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
	projects := h.store.GetAllProjects(user)
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
	var task TaskRecord
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}
	user, _, _ := r.BasicAuth()
	// Add to global tasks slice
	if task.ID == "" {
		task.ID = fmt.Sprintf("task_%s", uuid.New().String())
	}

	task.UpdatedTime = time.Now()
	taskExists := h.store.HasTask(user, project, task.ID)
	if taskExists {
		if success := h.store.UpdateTask(user, project, task.ID, task); !success {
			http.Error(w, "Failed to update task", http.StatusInternalServerError)
			return
		}
	} else {
		if success := h.store.AddTask(user, project, task.ID, task); !success {
			http.Error(w, "Failed to add task", http.StatusInternalServerError)
			return
		}
	}

	// Return success message
	w.WriteHeader(http.StatusOK)
	if taskExists {
		fmt.Fprintf(w, "Updated task with key: %s at: %v", task.ID, task.UpdatedTime)
		return
	}
	fmt.Fprintf(w, "Added task with key: %s at: %v", task.ID, task.UpdatedTime)
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
	if !h.store.HasTask(user, project, key) {
		http.Error(w, "Task does not exist", http.StatusBadRequest)
		return
	}
	task, _ := h.store.GetTask(user, project, key)
	task.Completed = true
	h.store.UpdateTask(user, project, key, task)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "task: %s completed", task.ID)
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
	if !h.store.HasTask(user, project, key) {
		http.Error(w, "Task does not exist", http.StatusBadRequest)
		return
	}
	task, _ := h.store.GetTask(user, project, key)
	h.store.RemoveTask(user, project, key)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "task: %s removed", task.ID)
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
	h.store.RemoveProject(user, project)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "project: %s removed", project)
}
