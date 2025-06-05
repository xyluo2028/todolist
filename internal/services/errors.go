package services

import "errors"

// ErrTaskNotFound is returned when a task is not found.
var ErrTaskNotFound = errors.New("task not found")
