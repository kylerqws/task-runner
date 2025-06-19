package task

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/kylerqws/task-runner/internal/domain/model"
)

// DefaultTask simulates a task with a fixed delay and random failure chance.
type DefaultTask struct {
	meta  *model.Task
	rng   *rand.Rand
	delay time.Duration
}

// NewDefaultTask creates a new DefaultTask with the given metadata, RNG, and delay.
func NewDefaultTask(meta *model.Task, rng *rand.Rand, delay time.Duration) *DefaultTask {
	return &DefaultTask{meta: meta, rng: rng, delay: delay}
}

// Run simulates task execution by sleeping for a predefined delay.
// It randomly returns an error to mimic failure in ~40% of cases.
func (t *DefaultTask) Run() error {
	time.Sleep(t.delay)
	if t.rng.Intn(100) >= 60 {
		return fmt.Errorf("simulated task failure")
	}
	return nil
}
