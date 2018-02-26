package dao

import (
	"errors"
	"strconv"

	"github.com/so-technology-watch/go-backend-sample/model"

	"cloud.google.com/go/datastore"
	"golang.org/x/net/context"
)

const entityName = "Task"

var _ TaskDAO = (*TaskDAODatastore)(nil)

// TaskDAODatastore is the datastore implementation of the TaskDAO
type TaskDAODatastore struct {
	client *datastore.Client
}

// NewTaskDAODatastore creates a new TaskDAO with a datastore implementation
func NewTaskDAODatastore(client *datastore.Client) TaskDAO {
	return &TaskDAODatastore{
		client: client,
	}
}

// Get return a task by its id
func (dao *TaskDAODatastore) Get(id string) (*model.Task, error) {
	ctx := context.Background()
	key := dao.datastoreKey(id)
	task := &model.Task{}
	if err := dao.client.Get(ctx, key, task); err != nil {
		return nil, err
	}
	task.Id = id
	return task, nil
}

// GetAll return all tasks
func (dao *TaskDAODatastore) GetAll() ([]model.Task, error) {
	ctx := context.Background()
	tasks := make([]model.Task, 0)
	query := datastore.NewQuery(entityName).Order("Id")

	keys, err := dao.client.GetAll(ctx, query, &tasks)
	if err != nil {
		return nil, err
	}

	for i, key := range keys {
		tasks[i].Id = strconv.FormatInt(key.ID, 10)
	}

	return tasks, nil
}

// Upsert update or create a task
func (dao *TaskDAODatastore) Upsert(task *model.Task) (*model.Task, error) {
	ctx := context.Background()
	var key *datastore.Key
	if task.Id == "" {
		key = datastore.IncompleteKey(entityName, nil)
	} else {
		key = dao.datastoreKey(task.Id)
	}

	key, err := dao.client.Put(ctx, key, task)
	if err != nil {
		return nil, err
	}
	task.Id = strconv.FormatInt(key.ID, 10)
	return task, nil
}

// Delete delete a task by its id
func (dao *TaskDAODatastore) Delete(id string) error {
	ctx := context.Background()
	key := dao.datastoreKey(id)
	return dao.client.Delete(ctx, key)
}

// DeleteAll deletes all tasks
func (dao *TaskDAODatastore) DeleteAll() error {
	ctx := context.Background()
	tasks := make([]model.Task, 0)
	query := datastore.NewQuery(entityName).KeysOnly()
	keys, err := dao.client.GetAll(ctx, query, &tasks)
	if err != nil {
		return err
	}
	return dao.client.DeleteMulti(ctx, keys)
}

// Exist check if the task exist
func (dao *TaskDAODatastore) Exist(id string) (bool, error) {
	ctx := context.Background()
	key := dao.datastoreKey(id)
	task := &model.Task{}
	if err := dao.client.Get(ctx, key, task); err != nil {
		return false, errors.New("task not found with id " + id)
	}
	return true, nil
}

func (dao *TaskDAODatastore) datastoreKey(idStr string) *datastore.Key {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil
	}
	return datastore.IDKey(entityName, id, nil)
}
