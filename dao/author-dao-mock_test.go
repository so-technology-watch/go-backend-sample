package dao_test

import (
	"testing"
	"go-redis-sample/dao"
	"go-redis-sample/model"
)

func TestAuthorDAOMock(t *testing.T) {

	authorDaoMock, _, err := dao.GetDAO(dao.MockDAO, dao.DBConfigFileName)
	if err != nil {
		t.Error(err)
	}

	authorToSave := model.Author{
		Id:           "1",
		Firstname:    "TestMock",
		Lastname:  	  "TestMock",
	}

	authorSaved, err := authorDaoMock.Upsert(&authorToSave)
	if err != nil {
		t.Error(err)
	}

	t.Log("author saved", authorSaved)

	authors, err := authorDaoMock.GetAll()
	if err != nil {
		t.Error(err)
	}
	if len(authors) != 1 {
		t.Errorf("Expected 1 authors, got %d", len(authors))
	}

	oneAuthor, err := authorDaoMock.Get(authorToSave.Id)
	if err != nil {
		t.Error(err)
	}
	if authorToSave != *oneAuthor {
		t.Error("Got wrong author by id")
	}

	err = authorDaoMock.Delete(oneAuthor.Id)
	if err != nil {
		t.Error(err)
	}

	oneAuthor, err = authorDaoMock.Get(oneAuthor.Id)
	if err == nil {
		t.Error("Author should have been deleted", oneAuthor)
	}
}
