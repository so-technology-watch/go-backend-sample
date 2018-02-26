package dao

import (
	"database/sql"
	"errors"

	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"github.com/so-technology-watch/go-backend-sample/model"
)

var _ TaskDAO = (*TaskDAOSQL)(nil)

// TaskDAOSQL is the sql implementation of the TaskDAO
type TaskDAOSQL struct {
	sqlSession *sql.DB
}

// NewTaskDAOSQL creates a new TaskDAO sql implementation
func NewTaskDAOSQL(sqlSession *sql.DB) TaskDAO {
	return &TaskDAOSQL{
		sqlSession: sqlSession,
	}
}

// Get return a task by its id
func (dao *TaskDAOSQL) Get(id string) (*model.Task, error) {
	task, err := dao.get(id)
	if err != nil {
		return nil, err
	}

	return task, nil
}

// GetAll return all tasks
func (dao *TaskDAOSQL) GetAll() ([]model.Task, error) {
	var tasks []model.Task

	// Collect all tasks identifiers
	rows, err := dao.sqlSession.Query("SELECT * FROM task")
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
func (dao *TaskDAOSQL) Upsert(task *model.Task) (*model.Task, error) {
	if len(task.Id) == 0 {
		task.Id = uuid.NewV4().String()

		// insert
		stmt, err := dao.sqlSession.Prepare("INSERT INTO task (id, title, description, status, creationDate, modificationDate) VALUES (?,?,?,?,?,?)")
		if err != nil {
			return nil, err
		}

		_, err = stmt.Exec(task.Id, task.Title, task.Description, task.Status, task.CreationDate, task.ModificationDate)
		if err != nil {
			return nil, err
		}
	} else {
		// update
		stmt, err := dao.sqlSession.Prepare("UPDATE task SET title=?,description=?,status=?,creationDate=?,modificationDate=? WHERE id=?")
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
func (dao *TaskDAOSQL) Delete(id string) error {
	stmt, err := dao.sqlSession.Prepare("DELETE FROM task WHERE id=?")
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
func (dao *TaskDAOSQL) DeleteAll() error {
	stmt, err := dao.sqlSession.Prepare("DELETE * FROM task")
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
func (dao *TaskDAOSQL) Exist(id string) (bool, error) {
	task := model.Task{}

	err := dao.sqlSession.QueryRow("SELECT * FROM task WHERE id=?", id).
		Scan(&task.Id, &task.Title, &task.Description, &task.CreationDate, &task.ModificationDate)
	if err != nil {
		return false, err
	}

	return true, nil
}

// get a task by its id
func (dao *TaskDAOSQL) get(id string) (*model.Task, error) {
	task := model.Task{}

	err := dao.sqlSession.QueryRow("SELECT * FROM task WHERE id=?", id).
		Scan(&task.Id, &task.Title, &task.Description, &task.Status, &task.CreationDate, &task.ModificationDate)
	if err != nil {
		return nil, err
	}

	return &task, nil
}
