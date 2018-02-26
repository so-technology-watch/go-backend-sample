package dao

import (
	"github.com/so-technology-watch/go-backend-sample/model"
)

// TaskDAO is the interface for Task
type TaskDAO interface {
	Get(id string) (*model.Task, error)

	GetAll() ([]model.Task, error)

	Upsert(task *model.Task) (*model.Task, error)

	Delete(id string) error

	DeleteAll() error

	Exist(id string) (bool, error)
}
