package dao_test

import (
	"go-backend-sample/dao"
	"go-backend-sample/model"
	"testing"
)

func TestTaskDAOMongo(t *testing.T) {
	taskDao, err := dao.GetDAO(dao.MongoDAO, "")
	if err != nil {
		t.Error(err)
	}

	toSave := model.Task{
		Id:          uuid.NewV4().String(),
		Title:       "Title Mongo",
		Description: "Description Mongo",
	}

	taskSaved, err := taskDao.Upsert(&toSave)
	if err != nil {
		t.Error(err)
	}

	t.Log("task saved", taskSaved)

	tasks, err := taskDao.GetAll()
	if err != nil {
		t.Error(err)
	}

	t.Log("task found all", tasks[0])

	oneTask, err := taskDao.Get(tasks[0].Id)
	if err != nil {
		t.Error(err)
	}

	t.Log("task found one", oneTask)

	oneTask.Title = "Title Mongo 2"
	oneTask.Description = "Description Mongo 2"
	chg, err := taskDao.Upsert(oneTask)
	if err != nil {
		t.Error(err)
	}

	t.Log("task modified", chg, oneTask)

	oneTask, err = taskDao.Get(oneTask.Id)
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
		t.Log("task deleted", err)
	}
}
