package task

import (
	"math/rand"
	"time"

	"github.com/kylerqws/task-runner/internal/domain/model"
)

type DefaultTaskFactory struct {
	Rng   *rand.Rand
	Delay time.Duration
}

func (f *DefaultTaskFactory) New(task *model.Task) ExecutableTask {
	return NewDefaultTask(task, f.Rng, f.Delay)
}
