package dao

import (
	"errors"
	"github.com/so-technology-watch/go-backend-sample/model"

	"github.com/satori/go.uuid"
)

var _ TaskDAO = (*TaskDAOMock)(nil)

// MockedTask is the task returned by this mocked interface
var MockedTask = model.Task{
	Id:          uuid.NewV4().String(),
	Title:       "TestMock",
	Description: "TestMock",
}

// TaskDAOMock is the mocked implementation of the TaskDAO
type TaskDAOMock struct {
	storage map[string]*model.Task
}

// NewTaskDAOMock creates a new TaskDAO with a mocked implementation
func NewTaskDAOMock() TaskDAO {
	daoMock := &TaskDAOMock{
		storage: make(map[string]*model.Task),
	}

	// Adds some fake data
	daoMock.Upsert(&MockedTask)

	return daoMock
}

// Get return a task by its id
func (dao *TaskDAOMock) Get(id string) (*model.Task, error) {
	task, ok := dao.storage[id]
	if !ok {
		return nil, errors.New("task not found with id " + id)
	}
	return task, nil
}

// GetAll return all tasks
func (dao *TaskDAOMock) GetAll() ([]model.Task, error) {
	var tasks []model.Task
	for taskId := range dao.storage {
		task := dao.storage[taskId]
		tasks = append(tasks, *task)
	}
	return tasks, nil
}

// Upsert update or create a task
func (dao *TaskDAOMock) Upsert(task *model.Task) (*model.Task, error) {
	if task.Id == "" {
		task.Id = uuid.NewV4().String()
	}
	dao.save(task)
	return task, nil
}

// Delete delete a task by its id
func (dao *TaskDAOMock) Delete(id string) error {
	delete(dao.storage, id)
	return nil
}

// DeleteAll deletes all tasks
func (dao *TaskDAOMock) DeleteAll() error {
	for taskId := range dao.storage {
		delete(dao.storage, taskId)
	}
	return nil
}

// Exist check if the task exist
func (dao *TaskDAOMock) Exist(id string) (bool, error) {
	if dao.storage[id] != nil {
		return true, nil
	}
	return false, errors.New("task not found with id " + id)
}

// save the task
func (dao *TaskDAOMock) save(task *model.Task) *model.Task {
	dao.storage[task.Id] = task
	return task
}

// get a task by its id
func (dao *TaskDAOMock) get(id string) *model.Task {
	task := dao.storage[id]
	return task
}
