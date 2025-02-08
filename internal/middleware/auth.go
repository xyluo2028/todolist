package middleware

import (
	"fmt"
	"net/http"
	"todolist/internal/store"
)

type AuthMiddleware struct {
	userStore *store.UserStore
}

func NewAuthMiddleware(userStore *store.UserStore) *AuthMiddleware {
	return &AuthMiddleware{
		userStore: userStore,
	}
}

func (m *AuthMiddleware) Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, _ := r.BasicAuth()
		fmt.Fprintf(w, "Authenticating user %s ... \n \n", username)
		if !m.userStore.Authenticate(username, password) {
			w.Header().Set("WWW-Authenticate", `Basic realm="Todo App"`)
			fmt.Fprintf(w, "User %s authentication failed \n", username)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		fmt.Fprintf(w, "User %s authenticated successfully \n \n", username)
		next(w, r)
	}
}
