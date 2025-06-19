package service

import (
	"fmt"
	"time"

	"github.com/kylerqws/task-runner/internal/domain/model"
	"github.com/kylerqws/task-runner/internal/domain/task"
)

const (
	taskQueueBufferSize        = 100                    // Maximum number of pending tasks per type
	taskDurationUpdateInterval = 500 * time.Millisecond // Interval to update task duration
)

// queueExists checks whether a queue already exists for the given task type.
func (m *TaskManager) queueExists(taskType string) bool {
	m.mu.RLock()
	_, ok := m.queues[taskType]
	m.mu.RUnlock()

	return ok
}

// createQueue initializes a queue and a background worker for the given task type.
func (m *TaskManager) createQueue(taskType string, factory task.Factory) {
	m.mu.Lock()
	m.queues[taskType] = make(chan *model.Task, taskQueueBufferSize)
	m.mu.Unlock()

	go func() {
		for t := range m.queues[taskType] {
			exec := factory.New(t)
			m.runExecutableTask(t, exec)
		}
	}()
}

// runExecutableTask updates task status, tracks execution time, runs the task,
// and finalizes its result.
func (m *TaskManager) runExecutableTask(t *model.Task, exec task.ExecutableTask) {
	t.Status = model.TaskStatusRunning

	start := time.Now()
	stop := m.trackDuration(t, start)

	err := exec.Run()
	stop()

	m.finalizeTask(t, err)
}

// trackDuration updates the task's duration field every interval.
// Returns a function to stop the tracker when the task is done.
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

// finalizeTask sets the final task status and result message based on execution outcome.
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

// updateDuration updates the task duration field based on elapsed time.
func (m *TaskManager) updateDuration(t *model.Task, start time.Time) {
	m.mu.Lock()
	defer m.mu.Unlock()

	t.Duration = time.Since(start).Truncate(time.Second).String()
}
