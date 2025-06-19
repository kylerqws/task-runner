package task

import (
	"math/rand"
	"time"

	"github.com/kylerqws/task-runner/internal/domain/model"
)

// DefaultTaskFactory creates instances of DefaultTask with configured RNG and delay.
type DefaultTaskFactory struct {
	Rng   *rand.Rand    // Random number generator
	Delay time.Duration // Execution delay duration
}

// New creates a new DefaultTask using the factory's configuration.
func (f *DefaultTaskFactory) New(task *model.Task) ExecutableTask {
	return NewDefaultTask(task, f.Rng, f.Delay)
}
