package handlers

import (
	"fmt"
	"net/http"
	"time"
	"todolist/internal/store"
)

type TaskHandler struct {
	store *store.TaskStore
}

func NewTaskHandler(store *store.TaskStore) *TaskHandler {
	return &TaskHandler{
		store: store,
	}
}

func (h *TaskHandler) PrintTasksHttp(w http.ResponseWriter, r *http.Request) {
	user, _, _ := r.BasicAuth()
	tasks := h.store.GetAllTasks(user)
	w.WriteHeader(http.StatusOK)
	if len(tasks) == 0 {
		fmt.Fprintln(w, "No tasks found!")
		return
	}

	fmt.Fprintln(w, "Todos: ")
	for key, task := range tasks {
		fmt.Fprintf(w, "Key: %s, Task: %s \n", key, task)
	}
}

func (h *TaskHandler) AddTaskHttp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get task from query parameter
	task := r.URL.Query().Get("q")
	user, _, _ := r.BasicAuth()
	if task == "" {
		http.Error(w, "Task query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	// Add to global tasks slice
	key := fmt.Sprintf("task_%d", time.Now().UnixNano())
	if success := h.store.AddTask(user, key, task); !success {
		http.Error(w, "Task already exists", http.StatusBadRequest)
		return
	}

	// Return success message
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Added task with key: %s", key)
}

func (h *TaskHandler) RemoveTaskHttp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	key := r.URL.Query().Get("key")
	user, _, _ := r.BasicAuth()
	if key == "" {
		http.Error(w, "Key query parameter 'key' is required", http.StatusBadRequest)
		return
	}
	if !h.store.HasTask(user, key) {
		http.Error(w, "Task does not exist", http.StatusBadRequest)
		return
	}
	task, _ := h.store.GetTask(user, key)
	h.store.RemoveTask(user, key)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "task: %s removed", task)
}

func (h *TaskHandler) UpdateTaskHttp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	key := r.URL.Query().Get("key")
	user, _, _ := r.BasicAuth()
	task := r.URL.Query().Get("q")
	if key == "" {
		http.Error(w, "Key query parameter 'key' is required", http.StatusBadRequest)
		return
	}

	if task == "" {
		http.Error(w, "Task query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	if !h.store.HasTask(user, key) {
		http.Error(w, "Task does not exist", http.StatusBadRequest)
		return
	}

	h.store.UpdateTask(user, key, task)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "task updated")
}
