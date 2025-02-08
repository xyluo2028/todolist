package handlers

import (
	"fmt"
	"net/http"
)

type WelcomeHandler struct{}

func (h *WelcomeHandler) Welcome(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintf(w, "Welcome to simple todo list!")
}

func NewWelcomeHandler() *WelcomeHandler {
	return &WelcomeHandler{}
}
