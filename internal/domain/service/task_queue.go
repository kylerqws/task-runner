package service

import (
	"fmt"
	"time"

	"github.com/kylerqws/task-runner/internal/domain/model"
	"github.com/kylerqws/task-runner/internal/domain/task"
)

const (
	taskQueueBufferSize        = 100                    // Max number of tasks in the queue
	taskDurationUpdateInterval = 500 * time.Millisecond // Duration update interval
)

// workerLoop processes tasks from the queue in order for a given type.
func (m *TaskManager) workerLoop(taskType string) {
	for {
		m.mu.Lock()
		queue := m.queues[taskType]
		factory := m.factories[taskType]

		if len(queue) == 0 {
			m.mu.Unlock()
			time.Sleep(100 * time.Millisecond)
			continue
		}

		t := queue[0]
		m.queues[taskType] = queue[1:]
		m.mu.Unlock()

		exec := factory.New(t)
		m.runExecutableTask(t, exec)

		m.mu.Lock()
		m.active[t.Type]--
		m.mu.Unlock()
	}
}

// runExecutableTask runs the task and finalizes its result.
func (m *TaskManager) runExecutableTask(t *model.Task, exec task.ExecutableTask) {
	t.Status = model.TaskStatusRunning

	start := time.Now()
	stop := m.trackDuration(t, start)

	err := exec.Run()
	stop()

	m.finalizeTask(t, err)
}

// trackDuration updates task duration while it's running.
func (m *TaskManager) trackDuration(t *model.Task, start time.Time) func() {
	ticker := time.NewTicker(taskDurationUpdateInterval)
	done := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				m.updateDuration(t, start)
			case <-done:
				m.updateDuration(t, start)
				ticker.Stop()
				return
			}
		}
	}()

	return func() { close(done) }
}

// finalizeTask sets task status and result after execution.
func (m *TaskManager) finalizeTask(t *model.Task, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if err != nil {
		t.Status = model.TaskStatusFailed
		t.Result = fmt.Sprintf("Task execution failed: %v", err)
		return
	}

	t.Status = model.TaskStatusDone
	t.Result = "Task completed successfully"
}

// updateDuration sets how long the task has been running.
func (m *TaskManager) updateDuration(t *model.Task, start time.Time) {
	m.mu.Lock()
	defer m.mu.Unlock()

	t.Duration = time.Since(start).Truncate(time.Second).String()
}

// enqueueTask adds a task to the queue and updates the counter.
// WARNING: Must be called with m.mu.Lock held.
func (m *TaskManager) enqueueTask(t *model.Task) {
	m.queues[t.Type] = append(m.queues[t.Type], t)
	m.active[t.Type]++
}

// removeFromQueue deletes a task from the queue and updates the counter.
// WARNING: Must be called with m.mu.Lock held.
func (m *TaskManager) removeFromQueue(t *model.Task) {
	q := m.queues[t.Type]

	for i := range q {
		if q[i].ID == t.ID {
			m.queues[t.Type] = append(q[:i], q[i+1:]...)
			m.active[t.Type]--
			break
		}
	}
}
