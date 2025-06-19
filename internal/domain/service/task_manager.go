package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"

	"github.com/kylerqws/task-runner/internal/domain/model"
	"github.com/kylerqws/task-runner/internal/domain/task"
)

type TaskManager struct {
	mu        sync.RWMutex
	tasks     map[string]*model.Task
	factories map[string]task.Factory
	queues    map[string]chan *model.Task
}

func NewTaskManager() *TaskManager {
	return &TaskManager{
		tasks:     make(map[string]*model.Task),
		factories: make(map[string]task.Factory),
		queues:    make(map[string]chan *model.Task),
	}
}

func (m *TaskManager) RegisterFactory(taskType string, factory task.Factory) {
	m.mu.Lock()
	m.factories[taskType] = factory
	m.mu.Unlock()

	if !m.queueExists(taskType) {
		m.createQueue(taskType, factory)
	}
}

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

func (m *TaskManager) GetTask(id string) (*model.Task, bool) {
	m.mu.RLock()
	t, ok := m.tasks[id]
	m.mu.RUnlock()

	return t, ok
}

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

func (m *TaskManager) generateID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic("failed to generate secure ID: " + err.Error())
	}

	return hex.EncodeToString(b)
}
