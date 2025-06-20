package router

import (
	"net/http"

	"github.com/kylerqws/task-runner/internal/transport/http/handler"
	"github.com/kylerqws/task-runner/internal/transport/http/response"
)

// InitTaskRouter initializes HTTP routing for task-related endpoints.
// It registers routes for creating, retrieving, and deleting tasks.
func InitTaskRouter(taskHandler *handler.TaskHandler) http.Handler {
	mux := http.NewServeMux()

	// POST /tasks
	mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			taskHandler.Create(w, r)
			return
		}

		http.Error(w, response.ErrMethodNotAllowed, http.StatusMethodNotAllowed)
	})

	// GET /tasks/{id}, DELETE /tasks/{id}
	mux.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			taskHandler.Get(w, r)
			return
		}

		if r.Method == http.MethodDelete {
			taskHandler.Delete(w, r)
			return
		}

		http.Error(w, response.ErrMethodNotAllowed, http.StatusMethodNotAllowed)
	})

	return mux
}
