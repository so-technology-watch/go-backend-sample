package dao_test

import (
	"go-backend-sample/dao"
	"go-backend-sample/model"
	"testing"
	"time"
)

func TestTaskDAOMySQL(t *testing.T) {
	taskDao, err := dao.GetDAO(dao.MySQLDAO, "")
	if err != nil {
		t.Error(err)
	}

	taskToSave := model.Task{
		Title:            "Title MySQL",
		Description:      "Description MySQL",
		Status:           0,
		CreationDate:     time.Now(),
		ModificationDate: time.Now(),
	}

	taskSaved, err := taskDao.Upsert(&taskToSave)
	if err != nil {
		t.Error(err)
	}

	t.Log("task saved", taskSaved)

	tasks, err := taskDao.GetAll()
	if err != nil {
		t.Error(err)
	}

	t.Log("task found all", tasks[0])

	oneTask, err := taskDao.Get(taskSaved.Id)
	if err != nil {
		t.Error(err)
	}

	t.Log("task found one", oneTask)

	oneTask.Title = "Title MySQL 2"
	oneTask.Description = "Description MySQL 2"
	oneTaskChanged, err := taskDao.Upsert(oneTask)
	if err != nil {
		t.Error(err)
	}

	t.Log("task modified", oneTaskChanged, oneTask)

	oneTask, err = taskDao.Get(oneTaskChanged.Id)
	if err != nil {
		t.Error(err)
	}

	t.Log("task found one modified", oneTask)

	err = taskDao.Delete(oneTask.Id)
	if err != nil {
		t.Error(err)
	}

	oneTask, err = taskDao.Get(oneTask.Id)
	if err != nil {
		t.Log("task deleted", err, oneTask)
	}
}

func TestTaskDAOSQLite(t *testing.T) {
	taskDao, err := dao.GetDAO(dao.SQLiteDAO, "")
	if err != nil {
		t.Error(err)
	}

	taskToSave := model.Task{
		Title:            "Title SQLite",
		Description:      "Description SQLite",
		Status:           0,
		CreationDate:     time.Now(),
		ModificationDate: time.Now(),
	}

	taskSaved, err := taskDao.Upsert(&taskToSave)
	if err != nil {
		t.Error(err)
	}

	t.Log("task saved", taskSaved)

	tasks, err := taskDao.GetAll()
	if err != nil {
		t.Error(err)
	}

	t.Log("task found all", tasks[0])

	oneTask, err := taskDao.Get(taskSaved.Id)
	if err != nil {
		t.Error(err)
	}

	t.Log("task found one", oneTask)

	oneTask.Title = "Title SQLite 2"
	oneTask.Description = "Description SQLite 2"
	oneTaskChanged, err := taskDao.Upsert(oneTask)
	if err != nil {
		t.Error(err)
	}

	t.Log("task modified", oneTaskChanged, oneTask)

	oneTask, err = taskDao.Get(oneTaskChanged.Id)
	if err != nil {
		t.Error(err)
	}

	t.Log("task found one modified", oneTask)

	err = taskDao.Delete(oneTask.Id)
	if err != nil {
		t.Error(err)
	}

	oneTask, err = taskDao.Get(oneTask.Id)
	if err != nil {
		t.Log("task deleted", err, oneTask)
	}
}
