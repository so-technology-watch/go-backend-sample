package dao

import (
	"go-backend-sample/model"
)

type TaskDAO interface {

	Get(id string) (*model.Task, error)

	GetAll() ([]model.Task, error)

	Upsert(task *model.Task) (*model.Task, error)

	Delete(id string) error

	DeleteAll() error

	Exist(id string) (bool, error)
}
