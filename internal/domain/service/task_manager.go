package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"

	"github.com/kylerqws/task-runner/internal/domain/model"
	"github.com/kylerqws/task-runner/internal/domain/task"
)

// TaskManager manages task creation, execution, lookup, and deletion.
type TaskManager struct {
	mu        sync.RWMutex
	tasks     map[string]*model.Task   // All tasks by ID
	factories map[string]task.Factory  // Task type -> factory
	queues    map[string][]*model.Task // Task type -> task queue
	active    map[string]int           // Task type -> active count
}

// NewTaskManager returns a new instance with empty internal maps.
func NewTaskManager() *TaskManager {
	return &TaskManager{
		tasks:     make(map[string]*model.Task),
		factories: make(map[string]task.Factory),
		queues:    make(map[string][]*model.Task),
		active:    make(map[string]int),
	}
}

// RegisterFactory sets up a task type with its factory and starts the worker.
func (m *TaskManager) RegisterFactory(taskType string, factory task.Factory) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.factories[taskType]; !ok {
		m.factories[taskType] = factory
		m.queues[taskType] = []*model.Task{}
		go m.workerLoop(taskType)
	}
}

// CreateTask adds a new task to the queue if the type is known and not full.
func (m *TaskManager) CreateTask(taskType string) (*model.Task, error) {
	if taskType == "" {
		return nil, fmt.Errorf("cannot create task: %w", ErrTaskUnknownType)
	}

	m.mu.RLock()
	_, typeExists := m.factories[taskType]
	activeCount := m.active[taskType]
	m.mu.RUnlock()

	if !typeExists {
		return nil, fmt.Errorf("cannot create task with type %q: %w", taskType, ErrTaskUnknownType)
	}
	if activeCount >= taskQueueBufferSize {
		return nil, fmt.Errorf("cannot create task with type %q: %w", taskType, ErrTaskQueueLimitReached)
	}

	id := m.generateID()

	m.mu.RLock()
	_, taskExists := m.tasks[id]
	m.mu.RUnlock()

	if taskExists {
		return nil, fmt.Errorf("cannot create task with ID %q: %w", id, ErrTaskAlreadyExists)
	}

	t := model.NewTask(id, taskType)

	m.mu.Lock()
	m.tasks[t.ID] = t
	m.enqueueTask(t)
	m.mu.Unlock()

	return t, nil
}

// GetTask returns a task by ID or an error if not found.
func (m *TaskManager) GetTask(id string) (*model.Task, error) {
	m.mu.RLock()
	t, taskExists := m.tasks[id]
	m.mu.RUnlock()

	if !taskExists {
		return nil, fmt.Errorf("cannot find task with ID %q: %w", id, ErrTaskNotFound)
	}
	return t, nil
}

// DeleteTask removes a task if it's not running.
func (m *TaskManager) DeleteTask(id string) error {
	m.mu.RLock()
	t, taskExists := m.tasks[id]
	m.mu.RUnlock()

	if !taskExists {
		return fmt.Errorf("cannot delete task with ID %q: %w", id, ErrTaskNotFound)
	}
	if t.Status == model.TaskStatusRunning {
		return fmt.Errorf("cannot delete task with ID %q: %w", id, ErrTaskInProgress)
	}

	m.mu.Lock()
	delete(m.tasks, id)
	m.removeFromQueue(t)
	m.mu.Unlock()

	return nil
}

// generateID returns a secure random 128-bit hex string.
func (m *TaskManager) generateID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic(fmt.Sprintf("failed to generate secure ID: %v", err))
	}

	return hex.EncodeToString(b)
}
