package dao

import (
	"errors"
	"go-backend-sample/model"
	"github.com/satori/go.uuid"
)

var _ TaskDAO = (*TaskDAOMock)(nil)

// TaskDAOMock is the mocked implementation of the TaskDAO
type TaskDAOMock struct {
	storage map[string]*model.Task
}

// NewTaskDAOMock creates a new TaskDAO with a mocked implementation
func NewTaskDAOMock() TaskDAO {
	daoMock := &TaskDAOMock{
		storage: make(map[string]*model.Task),
	}

	return daoMock
}

// Get return a task by its id
func (s *TaskDAOMock) Get(id string) (*model.Task, error) {
	task, ok := s.storage[id]
	if !ok {
		return nil, errors.New("task not found with id " + id)
	}
	return task, nil
}

// Upsert update or create a task
func (s *TaskDAOMock) Upsert(task *model.Task) (*model.Task, error) {
	if task.Id == "" {
		task.Id = uuid.NewV4().String()
	}
	s.save(task)
	return task, nil
}

// Delete delete a task by its id
func (s *TaskDAOMock) Delete(id string) error {
	delete(s.storage, id)
	return nil
}

// Exist check if the task exist
func (s *TaskDAOMock) Exist(id string) (bool, error) {
	if s.storage[id] != nil {
		return true, nil
	}
	return false, errors.New("task not found with id " + id)
}

// save the task
func (s *TaskDAOMock) save(task *model.Task) *model.Task {
	s.storage[task.Id] = task
	return task
}
