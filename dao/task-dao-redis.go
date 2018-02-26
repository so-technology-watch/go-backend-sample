package dao

import (
	"encoding/json"
	"errors"

	"github.com/so-technology-watch/go-backend-sample/model"

	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"gopkg.in/redis.v5"
)

var _ TaskDAO = (*TaskDAORedis)(nil)

// TaskDAORedis is the redis implementation of the TaskDAO
type TaskDAORedis struct {
	redisCli *redis.Client
}

// NewTaskDAORedis creates a new TaskDAO redis implementation
func NewTaskDAORedis(redisCli *redis.Client) TaskDAO {
	return &TaskDAORedis{
		redisCli: redisCli,
	}
}

// Get return a task by its id
func (dao *TaskDAORedis) Get(id string) (*model.Task, error) {
	task, err := dao.get(id)
	if err != nil {
		return nil, err
	}

	return task, nil
}

// GetAll return all tasks
func (dao *TaskDAORedis) GetAll() ([]model.Task, error) {
	var tasks []model.Task

	// Collect all tasks identifiers
	keys := dao.redisCli.Keys("*").Val()
	if len(keys) == 0 {
		return nil, errors.New("no tasks")
	}

	for i := 0; i < len(keys); i++ {
		// Collect task by identifier
		task, err := dao.get(keys[i])
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, *task)
	}

	return tasks, nil
}

// Upsert update or create a task, returns true if updated, false otherwise or on error
func (dao *TaskDAORedis) Upsert(task *model.Task) (*model.Task, error) {
	if len(task.Id) == 0 {
		task.Id = uuid.NewV4().String()
	}

	task, err := dao.save(task)
	if err != nil {
		return nil, err
	}

	return task, nil
}

// Delete delete a task by its id
func (dao *TaskDAORedis) Delete(id string) error {
	result, err := dao.redisCli.Del(id).Result()
	if err != nil {
		return err
	} else if result == 0 {
		return errors.New(id + " don't exist")
	}

	return nil
}

// DeleteAll delete all tasks
func (dao *TaskDAORedis) DeleteAll() error {
	// Collect all identifiers of tasks
	keys := dao.redisCli.Keys("*").Val()
	if len(keys) == 0 { // If no tasks in database
		logrus.Warn("no tasks")
	}

	for i := 0; i < len(keys); i++ {
		// Deletion of task by identifier
		_, err := dao.redisCli.Del(keys[i]).Result()
		if err != nil {
			return err
		}
	}

	return nil
}

// Exist check if the task exist
func (dao *TaskDAORedis) Exist(id string) (bool, error) {
	result, err := dao.redisCli.Exists(id).Result()
	if err != nil {
		return false, err
	} else if result == false {
		return false, errors.New(id + " don't exist")
	}
	return result, nil
}

// save the task
func (dao *TaskDAORedis) save(task *model.Task) (*model.Task, error) {
	result, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}

	status := dao.redisCli.Set(task.Id, string(result), 0)
	if status.Err() != nil {
		return nil, status.Err()
	}

	return task, nil
}

// get a task by its id
func (dao *TaskDAORedis) get(id string) (*model.Task, error) {
	result, err := dao.redisCli.Get(id).Result()
	if err != nil {
		return nil, err
	} else if len(result) == 0 {
		return nil, errors.New(id + " don't exist")
	}

	task := model.Task{}
	err = json.Unmarshal([]byte(result), &task)
	if err != nil {
		return nil, err
	}

	return &task, nil
}
