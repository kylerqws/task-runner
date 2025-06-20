package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"

	"github.com/kylerqws/task-runner/internal/domain/model"
	"github.com/kylerqws/task-runner/internal/domain/task"
)

// TaskManager coordinates the lifecycle of tasks: creation, execution, lookup, and deletion.
// Each task type has its own queue and active counter.
type TaskManager struct {
	mu        sync.RWMutex                // Synchronizes access to all internal maps
	tasks     map[string]*model.Task      // Stores all tasks by their unique ID
	factories map[string]task.Factory     // Registered factories for each task type
	queues    map[string]chan *model.Task // Execution queues per task type
	active    map[string]int              // Number of active tasks per task type
}

// NewTaskManager initializes and returns a new TaskManager instance.
func NewTaskManager() *TaskManager {
	return &TaskManager{
		tasks:     make(map[string]*model.Task),
		factories: make(map[string]task.Factory),
		queues:    make(map[string]chan *model.Task),
		active:    make(map[string]int),
	}
}

// RegisterFactory binds a task type to a factory,
// and creates a queue for the type if not already present.
func (m *TaskManager) RegisterFactory(taskType string, factory task.Factory) {
	m.mu.Lock()
	m.factories[taskType] = factory
	m.mu.Unlock()

	if !m.queueExists(taskType) {
		m.createQueue(taskType, factory)
	}
}

// CreateTask creates and queues a new task of the given type.
// Returns an error if the type is unknown or the active queue is full.
func (m *TaskManager) CreateTask(taskType string) (*model.Task, error) {
	if taskType == "" {
		return nil, fmt.Errorf("cannot create task: %w", ErrTaskUnknownType)
	}

	m.mu.RLock()
	_, ok := m.factories[taskType]
	m.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("cannot create task with type %q: %w", taskType, ErrTaskUnknownType)
	}

	m.mu.RLock()
	active, _ := m.active[taskType]
	m.mu.RUnlock()

	if active >= taskQueueBufferSize {
		return nil, fmt.Errorf("cannot create task with type %q: %w", taskType, ErrTaskQueueLimitReached)
	}

	id := m.generateID()

	m.mu.RLock()
	_, ok = m.tasks[id]
	m.mu.RUnlock()

	if ok {
		return nil, fmt.Errorf("cannot create task with ID %q: %w", id, ErrTaskAlreadyExists)
	}

	t := model.NewTask(id, taskType)

	m.mu.Lock()
	m.tasks[t.ID] = t
	m.queues[taskType] <- t
	m.active[taskType]++
	m.mu.Unlock()

	return t, nil
}

// GetTask returns a task by its ID from the internal registry.
// Returns an error if the task does not exist or was not previously created.
func (m *TaskManager) GetTask(id string) (*model.Task, error) {
	m.mu.RLock()
	t, ok := m.tasks[id]
	m.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("cannot find task with ID %q: %w", id, ErrTaskNotFound)
	}
	return t, nil
}

// DeleteTask removes a task by ID if it is not currently running.
// Returns an error if the task does not exist or is in progress.
func (m *TaskManager) DeleteTask(id string) error {
	m.mu.RLock()
	t, ok := m.tasks[id]
	m.mu.RUnlock()

	if !ok {
		return fmt.Errorf("cannot delete task with ID %q: %w", id, ErrTaskNotFound)
	}

	if t.Status == model.TaskStatusRunning {
		return fmt.Errorf("cannot delete task with ID %q: %w", id, ErrTaskIsProgress)
	}

	if t.Status == model.TaskStatusPending { // != TaskStatusDone and TaskStatusFailed
		m.mu.Lock()
		m.active[t.Type]--
		m.mu.Unlock()
	}

	m.mu.Lock()
	delete(m.tasks, id)
	m.mu.Unlock()

	return nil
}

// generateID creates a secure 128-bit hexadecimal task ID.
func (m *TaskManager) generateID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic(fmt.Sprintf("failed to generate secure ID: %v", err))
	}

	return hex.EncodeToString(b)
}
