package model

import (
	"errors"
	"github.com/satori/go.uuid"
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
	Id               string     `json:"id,omitempty" bson:"id"`
	Title            string     `json:"title" bson:"title"`
	Description      string     `json:"description" bson:"description"`
	Status           TaskStatus `json:"status" bson:"status"`
	CreationDate     time.Time  `json:"creationDate" bson:"creationDate"`
	ModificationDate time.Time  `json:"modificationDate" bson:"modificationDate"`
}

// Task constructor
func NewTask() *Task {
	return &Task{
		Id:           uuid.NewV4().String(),
		CreationDate: time.Now(),
		Status:       StatusTodo,
	}
}

// Validation of a task structure
func (task Task) Valid() error {
	if task.Title == "" {
		return errors.New("title is mandatory")
	}
	return nil
}

// Equal compares a task to another and returns true if they are equal false otherwise
func (t Task) Equal(task Task) bool {
	return t.Id == task.Id &&
		t.Title == task.Title &&
		t.Description == task.Description &&
		t.Status == task.Status &&
		t.CreationDate.Equal(task.CreationDate) &&
		t.ModificationDate.Equal(task.ModificationDate)
}
