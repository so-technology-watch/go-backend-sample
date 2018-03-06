package dao_test

import (
	"testing"

	"go-backend-sample/dao"
	"go-backend-sample/model"
)

func TestTaskDAOMock(t *testing.T) {

	taskDaoMock, err := dao.GetDAO(dao.MockDAO, "")
	if err != nil {
		t.Error(err)
	}

	taskToSave := model.Task{
		Id:          uuid.NewV4().String(),
		Title:       "Title Mock",
		Description: "Description Mock",
	}

	taskSaved, err := taskDaoMock.Upsert(&taskToSave)
	if err != nil {
		t.Error(err)
	}

	t.Log("task saved", taskSaved)

	tasks, err := taskDaoMock.GetAll()
	if err != nil {
		t.Error(err)
	}
	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(tasks))
	}

	oneTask, err := taskDaoMock.Get(taskToSave.Id)
	if err != nil {
		t.Error(err)
	}
	if taskSaved != oneTask {
		t.Error("Got wrong task by id")
	}

	err = taskDaoMock.Delete(oneTask.Id)
	if err != nil {
		t.Error(err)
	}

	oneTask, err = taskDaoMock.Get(oneTask.Id)
	if err == nil {
		t.Error("Task should have been deleted", oneTask)
	}
}
