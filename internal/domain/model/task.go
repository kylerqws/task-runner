package model

import "time"

// TaskStatus represents the current status of a task.
type TaskStatus string

const (
	TaskStatusPending TaskStatus = "pending"
	TaskStatusRunning TaskStatus = "running"
	TaskStatusDone    TaskStatus = "done"
	TaskStatusFailed  TaskStatus = "failed"
)

// Task holds metadata about an asynchronous task's lifecycle and result.
type Task struct {
	ID        string     `json:"id"`                 // Unique task identifier
	Type      string     `json:"type"`               // Type of the task (e.g. "default", etc.)
	Status    TaskStatus `json:"status"`             // Current task status
	CreatedAt time.Time  `json:"created_at"`         // Task creation timestamp
	Duration  string     `json:"duration,omitempty"` // Total execution time (if available)
	Result    string     `json:"result,omitempty"`   // Result message or error
}

// NewTask creates and returns a new Task with default status and creation time.
func NewTask(id, taskType string) *Task {
	return &Task{ID: id, Type: taskType, Status: TaskStatusPending, CreatedAt: time.Now()}
}
