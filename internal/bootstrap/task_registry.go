package bootstrap

import (
	"math/rand"
	"time"

	"github.com/kylerqws/task-runner/internal/domain/service"
	"github.com/kylerqws/task-runner/internal/domain/task"
)

// RegisterTaskFactories registers all available task factories
// to the provided TaskManager instance.
func RegisterTaskFactories(m *service.TaskManager) {
	m.RegisterFactory(task.DefaultTaskType, newDefaultTaskFactory())
}

// newDefaultTaskFactory returns a Factory for the "default" task type,
// preconfigured with a random delay between 3 and 5 minutes.
func newDefaultTaskFactory() task.Factory {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	delay := time.Duration(3+rng.Intn(3)) * time.Minute

	return &task.DefaultTaskFactory{Rng: rng, Delay: delay}
}
