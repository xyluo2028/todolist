package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"todolist/internal/services"
)

type UserHandler struct {
	userSvc *services.UserService
	taskSvc *services.TaskService
}

func NewUserHandler(userSvc *services.UserService, taskSvc *services.TaskService) *UserHandler {
	return &UserHandler{
		userSvc: userSvc,
		taskSvc: taskSvc,
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username and password required", http.StatusBadRequest)
		return
	}

	if err := h.userSvc.RegisterUser(req.Username, req.Password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("User %s registered successfully", req.Username)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User registered successfully",
	})
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	username, _, ok := r.BasicAuth()
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err := h.userSvc.DeactivateUser(username, h.taskSvc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	errClearAllTasks := h.taskSvc.RemoveUserTasks(username)
	if errClearAllTasks != nil {
		http.Error(w, errClearAllTasks.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("User %s deleted successfully", username)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User deleted successfully",
	})
}
