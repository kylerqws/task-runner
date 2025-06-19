package bootstrap

import (
	"math/rand"
	"time"

	"github.com/kylerqws/task-runner/internal/domain/service"
	"github.com/kylerqws/task-runner/internal/domain/task"
)

func RegisterTaskFactories(m *service.TaskManager) {
	m.RegisterFactory(task.DefaultTaskType, newDefaultTaskFactory())
}

func newDefaultTaskFactory() task.Factory {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	delay := time.Duration(3+rng.Intn(3)) * time.Minute

	return &task.DefaultTaskFactory{Rng: rng, Delay: delay}
}
