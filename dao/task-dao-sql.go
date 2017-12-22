package dao

import (
	"database/sql"
	"errors"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"go-backend-sample/model"
)

var _ TaskDAO = (*TaskDAOMySQL)(nil)

// TaskDAOMySQL is the mysql  implementation of the TaskDAO
type TaskDAOMySQL struct {
	mysqlSession *sql.DB
}

// NewTaskDAOMySQL creates a new TaskDAO mysql implementation
func NewTaskDAOMySQL(mysqlSession *sql.DB) TaskDAO {
	return &TaskDAOMySQL{
		mysqlSession: mysqlSession,
	}
}

// Get return a task by its id
func (dao *TaskDAOMySQL) Get(id string) (*model.Task, error) {
	task, err := dao.get(id)
	if err != nil {
		return nil, err
	}

	return task, nil
}

// GetAll return all tasks
func (dao *TaskDAOMySQL) GetAll() ([]model.Task, error) {
	var tasks []model.Task

	// Collect all tasks identifiers
	rows, err := dao.mysqlSession.Query("SELECT * FROM task")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var task model.Task

		err = rows.Scan(&task.Id, &task.Title, &task.Description, &task.Status, &task.CreationDate, &task.ModificationDate)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

// Upsert update or create a task, returns true if updated, false otherwise or on error
func (dao *TaskDAOMySQL) Upsert(task *model.Task) (*model.Task, error) {
	if len(task.Id) == 0 {
		task.Id = uuid.NewV4().String()

		// insert
		stmt, err := dao.mysqlSession.Prepare("INSERT task SET id=?,title=?,description=?,status=?,creationDate=?,modificationDate=?")
		if err != nil {
			return nil, err
		}

		_, err = stmt.Exec(task.Id, task.Title, task.Description, task.Status, task.CreationDate, task.ModificationDate)
		if err != nil {
			return nil, err
		}
	} else {
		// update
		stmt, err := dao.mysqlSession.Prepare("update task set title=?,description=?,status=?,creationDate=?,modificationDate=? where id=?")
		if err != nil {
			return nil, err
		}

		_, err = stmt.Exec(task.Title, task.Description, task.Status, task.CreationDate, task.ModificationDate, task.Id)
		if err != nil {
			return nil, err
		}
	}

	return task, nil
}

// Delete delete a task by its id
func (dao *TaskDAOMySQL) Delete(id string) error {
	stmt, err := dao.mysqlSession.Prepare("delete from task where id=?")
	if err != nil {
		return err
	}

	result, err := stmt.Exec(id)
	if err != nil {
		return err
	}

	affect, err := result.RowsAffected()
	if err != nil {
		return err
	} else if affect == 0 {
		return errors.New(id + " don't exist")
	}

	return nil
}

// DeleteAll delete all tasks
func (dao *TaskDAOMySQL) DeleteAll() error {
	stmt, err := dao.mysqlSession.Prepare("delete * from task")
	if err != nil {
		return err
	}

	result, err := stmt.Exec()
	if err != nil {
		return err
	}

	affect, err := result.RowsAffected()
	if err != nil {
		return err
	} else if affect == 0 {
		logrus.Warn("no tasks")
	}

	return nil
}

// Exist check if the task exist
func (dao *TaskDAOMySQL) Exist(id string) (bool, error) {
	task := model.Task{}

	err := dao.mysqlSession.QueryRow("SELECT * FROM task WHERE id=?", id).
		Scan(&task.Id, &task.Title, &task.Description, &task.CreationDate, &task.ModificationDate)
	if err != nil {
		return false, err
	}

	return true, nil
}

// get a task by its id
func (dao *TaskDAOMySQL) get(id string) (*model.Task, error) {
	task := model.Task{}

	err := dao.mysqlSession.QueryRow("SELECT * FROM task WHERE id=?", id).
		Scan(&task.Id, &task.Title, &task.Description, &task.Status, &task.CreationDate, &task.ModificationDate)
	if err != nil {
		return nil, err
	}

	return &task, nil
}
