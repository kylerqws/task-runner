package task

import "github.com/kylerqws/task-runner/internal/domain/model"

type ExecutableTask interface {
	Run() error
}

type Factory interface {
	New(task *model.Task) ExecutableTask
}
