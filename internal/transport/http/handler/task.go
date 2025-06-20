package handler

import (
	"errors"
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
	task, err := h.Manager.CreateTask(taskType)

	if err != nil {
		switch {
		case errors.Is(err, service.ErrTaskUnknownType):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, service.ErrTaskAlreadyExists):
			http.Error(w, err.Error(), http.StatusConflict)
		case errors.Is(err, service.ErrTaskQueueLimitReached):
			http.Error(w, err.Error(), http.StatusTooManyRequests)
		default:
			http.Error(w, response.ErrInternalServer, http.StatusInternalServerError)
		}
		return
	}

	response.RespondJSON(w, http.StatusCreated, task)
}

// Get handles GET /tasks/{id} and returns task details.
func (h *TaskHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/tasks/")
	task, err := h.Manager.GetTask(id)

	if err != nil {
		switch {
		case errors.Is(err, service.ErrTaskNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, response.ErrInternalServer, http.StatusInternalServerError)
		}
		return
	}

	response.RespondJSON(w, http.StatusOK, task)
}

// Delete handles DELETE /tasks/{id} and removes a task if it's not running.
func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/tasks/")
	err := h.Manager.DeleteTask(id)

	if err != nil {
		switch {
		case errors.Is(err, service.ErrTaskNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		case errors.Is(err, service.ErrTaskInProgress):
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			http.Error(w, response.ErrInternalServer, http.StatusInternalServerError)
		}
		return
	}

	response.RespondNoContent(w, http.StatusNoContent)
}
