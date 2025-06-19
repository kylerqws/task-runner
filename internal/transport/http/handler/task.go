package handler

import (
	"net/http"
	"strings"

	"github.com/kylerqws/task-runner/internal/domain/service"
	"github.com/kylerqws/task-runner/internal/transport/http/response"
)

// TaskHandler handles HTTP requests for task management operations.
type TaskHandler struct {
	Manager *service.TaskManager
}

// NewTaskHandler creates a new TaskHandler with the provided TaskManager.
func NewTaskHandler(manager *service.TaskManager) *TaskHandler {
	return &TaskHandler{Manager: manager}
}

// Create handles POST /tasks and creates a new task based on the given type.
func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	taskType := r.URL.Query().Get("type")
	if taskType == "" {
		http.Error(w, "task type is required", http.StatusBadRequest)
		return
	}

	task := h.Manager.CreateTask(taskType)
	response.RespondJSON(w, http.StatusCreated, task)
}

// Get handles GET /tasks/{id} and returns task details.
func (h *TaskHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/tasks/")
	task, ok := h.Manager.GetTask(id)
	if !ok {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}
	response.RespondJSON(w, http.StatusOK, task)
}

// Delete handles DELETE /tasks/{id} and removes a task if it's not running.
func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/tasks/")
	deleted, locked := h.Manager.DeleteTask(id)

	if locked {
		http.Error(w, "cannot delete running task", http.StatusConflict)
		return
	}
	if !deleted {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}
	response.RespondNoContent(w, http.StatusNoContent)
}
