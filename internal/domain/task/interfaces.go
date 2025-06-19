package task

import "github.com/kylerqws/task-runner/internal/domain/model"

// ExecutableTask defines the behavior of a task that can be executed.
type ExecutableTask interface {
	// Run executes the task logic and returns an error if it fails.
	Run() error
}

// Factory defines an interface for creating tasks of a specific type.
type Factory interface {
	// New creates a new ExecutableTask based on the provided Task metadata.
	New(task *model.Task) ExecutableTask
}
