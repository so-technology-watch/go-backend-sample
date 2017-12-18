package dao

import (
	"go-backend-sample/model"
)

// TaskDAO is the DAO interface to work with tasks
type TaskDAO interface {

	// Get return a task by its id
	Get(id string) (*model.Task, error)

	// GetAll return all tasks
	GetAll() ([]model.Task, error)

	// Upsert update or create a task, returns true if updated, false otherwise or on error
	Upsert(task *model.Task) (*model.Task, error)

	// Delete delete a task by its ID
	Delete(id string) error

	// DeleteAll delete all tasks
	DeleteAll() error

	// Exist check if the task exist
	Exist(id string) (bool, error)
}
