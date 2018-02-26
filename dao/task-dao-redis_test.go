package dao_test

import (
	"github.com/satori/go.uuid"
	"github.com/so-technology-watch/go-backend-sample/dao"
	"github.com/so-technology-watch/go-backend-sample/model"
	"testing"
)

func TestTaskDAORedis(t *testing.T) {
	taskDao, err := dao.GetDAO(dao.RedisDAO, "")
	if err != nil {
		t.Error(err)
	}

	taskToSave := model.Task{
		Id:          uuid.NewV4().String(),
		Title:       "Title Redis",
		Description: "Description Redis",
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

	oneTask.Title = "Title Redis 2"
	oneTask.Description = "Description Redis 2"
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
