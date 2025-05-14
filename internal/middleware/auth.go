package middleware

import (
	"fmt"
	"net/http"
	"todolist/internal/services"
)

type AuthMiddleware struct {
	userService *services.UserService
}

func NewAuthMiddleware(userService *services.UserService) *AuthMiddleware {
	return &AuthMiddleware{
		userService: userService,
	}
}

func (m *AuthMiddleware) Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, _ := r.BasicAuth()
		fmt.Fprintf(w, "Authenticating user %s ... \n \n", username)
		if !m.userService.AuthenticateUser(username, password) {
			// Authentication failed
			w.Header().Set("WWW-Authenticate", `Basic realm="Todo App"`)
			fmt.Fprintf(w, "User %s authentication failed \n", username)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		fmt.Fprintf(w, "User %s authenticated successfully \n \n", username)
		next(w, r)
	}
}
