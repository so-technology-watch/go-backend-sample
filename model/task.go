package model

import (
	"errors"
	"time"
)

// TaskStatus type of task status
type TaskStatus int

const (
	StatusTodo TaskStatus = iota
	StatusInProgress
	StatusDone
)

// Task structure of a task
type Task struct {
	Id               string     `json:"id,omitempty" bson:"id"`
	Title            string     `json:"title" bson:"title"`
	Description      string     `json:"description" bson:"description"`
	Status           TaskStatus `json:"status" bson:"status"`
	CreationDate     time.Time  `json:"creationDate" bson:"creationDate"`
	ModificationDate time.Time  `json:"modificationDate" bson:"modificationDate"`
}

// Valid check if task is valid
func (task Task) Valid() error {
	if task.Title == "" {
		return errors.New("title is mandatory")
	}
	if task.Description == "" {
		return errors.New("description is mandatory")
	}
	return nil
}

// Equal check if task is equal
func (t Task) Equal(task Task) bool {
	return t.Id == task.Id &&
		t.Title == task.Title &&
		t.Description == task.Description &&
		t.Status == task.Status &&
		t.CreationDate.Equal(task.CreationDate) &&
		t.ModificationDate.Equal(task.ModificationDate)
}
