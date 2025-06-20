package service

import "errors"

// Common task-related errors used in the TaskManager service layer.
var (
	ErrTaskNotFound          = errors.New("task not found")
	ErrTaskIsProgress        = errors.New("task in progress")
	ErrTaskAlreadyExists     = errors.New("task already exists")
	ErrTaskQueueLimitReached = errors.New("task queue limit reached")
	ErrTaskUnknownType       = errors.New("task unknown type")
)
