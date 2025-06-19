package task

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/kylerqws/task-runner/internal/domain/model"
)

type DefaultTask struct {
	meta  *model.Task
	rng   *rand.Rand
	delay time.Duration
}

func NewDefaultTask(meta *model.Task, rng *rand.Rand, delay time.Duration) *DefaultTask {
	return &DefaultTask{meta: meta, rng: rng, delay: delay}
}

func (t *DefaultTask) Run() error {
	time.Sleep(t.delay)
	if t.rng.Intn(100) >= 60 {
		return fmt.Errorf("simulated task failure")
	}
	return nil
}
