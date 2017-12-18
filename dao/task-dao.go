package dao

import (
	"go-backend-sample/model"
)

// TaskDAO is the DAO interface to work with tasks
type TaskDAO interface {

	// Get returns a task by its id
	Get(id string) (*model.Task, error)

	// GetAll returns all tasks
	GetAll() ([]model.Task, error)

	// Upsert updates or creates a task, returns true if updated, false otherwise or on error
	Upsert(task *model.Task) (*model.Task, error)

	// Delete deletes a task by its ID
	Delete(id string) error

	// DeleteAll deletes all tasks
	DeleteAll() error

	// Exist checks if the task exist
	Exist(id string) (bool, error)
}
