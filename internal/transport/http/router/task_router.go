package router

import (
	"net/http"

	"github.com/kylerqws/task-runner/internal/transport/http/handler"
)

func InitTaskRouter(taskHandler *handler.TaskHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			taskHandler.Create(w, r)
			return
		}
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	})

	mux.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			taskHandler.Get(w, r)
			return
		}
		if r.Method == http.MethodDelete {
			taskHandler.Delete(w, r)
			return
		}
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	})

	return mux
}
