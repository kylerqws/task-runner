package model

import "time"

type TaskStatus string

const (
	TaskStatusPending TaskStatus = "pending"
	TaskStatusRunning TaskStatus = "running"
	TaskStatusDone    TaskStatus = "done"
	TaskStatusFailed  TaskStatus = "failed"
)

type Task struct {
	ID        string     `json:"id"`
	Status    TaskStatus `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	Duration  string     `json:"duration,omitempty"`
	Result    string     `json:"result,omitempty"`
}

func NewTask(id string) *Task {
	return &Task{ID: id, Status: TaskStatusPending, CreatedAt: time.Now()}
}
