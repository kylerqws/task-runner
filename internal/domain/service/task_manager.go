package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"

	"github.com/kylerqws/task-runner/internal/domain/model"
	"github.com/kylerqws/task-runner/internal/domain/task"
)

// TaskManager manages the lifecycle, execution, and tracking of tasks.
// It holds registered task factories, task instances, and execution queues.
type TaskManager struct {
	mu        sync.RWMutex
	tasks     map[string]*model.Task
	factories map[string]task.Factory
	queues    map[string]chan *model.Task
}

// NewTaskManager creates a new TaskManager with empty maps for tasks,
// factories, and per-type queues.
func NewTaskManager() *TaskManager {
	return &TaskManager{
		tasks:     make(map[string]*model.Task),
		factories: make(map[string]task.Factory),
		queues:    make(map[string]chan *model.Task),
	}
}

// RegisterFactory registers a task factory for a specific task type,
// and initializes its execution queue if not already present.
func (m *TaskManager) RegisterFactory(taskType string, factory task.Factory) {
	m.mu.Lock()
	m.factories[taskType] = factory
	m.mu.Unlock()

	if !m.queueExists(taskType) {
		m.createQueue(taskType, factory)
	}
}

// CreateTask creates a new task of the given type and pushes it into the corresponding queue.
// If the type is unknown, a failed task is returned.
func (m *TaskManager) CreateTask(taskType string) *model.Task {
	id := m.generateID()
	t := model.NewTask(id)

	m.mu.Lock()
	m.tasks[id] = t
	queue, ok := m.queues[taskType]
	m.mu.Unlock()

	if !ok {
		t.Status = model.TaskStatusFailed
		t.Result = fmt.Sprintf("Unknown task type: %q", taskType)
		return t
	}

	queue <- t
	return t
}

// GetTask returns the task by ID and a boolean indicating whether it was found.
func (m *TaskManager) GetTask(id string) (*model.Task, bool) {
	m.mu.RLock()
	t, ok := m.tasks[id]
	m.mu.RUnlock()

	return t, ok
}

// DeleteTask removes a task by ID unless it's currently running.
// It returns a deletion success flag and a locked status flag.
func (m *TaskManager) DeleteTask(id string) (deleted bool, locked bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	t, ok := m.tasks[id]
	if !ok {
		return false, false
	}
	if t.Status == model.TaskStatusRunning {
		return false, true
	}

	delete(m.tasks, id)
	return true, false
}

// generateID produces a random 128-bit hexadecimal task ID.
func (m *TaskManager) generateID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic("failed to generate secure ID: " + err.Error())
	}

	return hex.EncodeToString(b)
}
