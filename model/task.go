package model

import (
	"time"
)

// TaskStatus is the current processing status of a task
type TaskStatus int

const (
	// StatusTodo is used for incomplete tasks
	StatusTodo TaskStatus = iota
	// StatusInProgress is used for tasks in progress
	StatusInProgress
	// StatusDone is used for completed tasks
	StatusDone
)

// Structure of a task
type Task struct {
	Id               string     `json:"id,omitempty"`
	Title            string     `json:"title"`
	Description      string     `json:"description"`
	Status           TaskStatus `json:"status"`
	CreationDate     time.Time  `json:"creationDate"`
	ModificationDate time.Time  `json:"modificationDate"`
}