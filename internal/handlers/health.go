package handlers

import (
	"fmt"
	"net/http"
)

// Health reports that the HTTP server has started successfully.
func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintln(w, "ok")
}
