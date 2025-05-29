package middleware

import (
	"log"
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
		username, password, ok := r.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="Todo App"`)
			http.Error(w, "Unauthorized: Basic auth required", http.StatusUnauthorized)
			return
		}
		log.Printf("Authenticating user %s ...", username)
		if !m.userService.AuthenticateUser(username, password) {
			// Authentication failed
			log.Printf("User %s authentication failed", username)
			w.Header().Set("WWW-Authenticate", `Basic realm="Todo App"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		log.Printf("User %s authenticated successfully", username)
		next(w, r)
	}
}
