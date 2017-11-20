package dao_test

import (
	"github.com/satori/go.uuid"
	"go-backend-sample/dao"
	"go-backend-sample/model"
	"testing"
)

func TestAuthorDAOMongo(t *testing.T) {
	authorDao, _, err := dao.GetDAO(dao.MongoDAO, "")
	if err != nil {
		t.Error(err)
	}

	toSave := model.Author{
		Id:        uuid.NewV4().String(),
		Firstname: "Test1",
		Lastname:  "Test2",
	}

	authorSaved, err := authorDao.Upsert(&toSave)
	if err != nil {
		t.Error(err)
	}

	t.Log("author saved", authorSaved)

	authors, err := authorDao.GetAll()
	if err != nil {
		t.Error(err)
	}

	t.Log("author found all", authors[0])

	oneAuthor, err := authorDao.Get(authors[0].Id)
	if err != nil {
		t.Error(err)
	}

	t.Log("author found one", oneAuthor)

	oneAuthor.Firstname = "Test3"
	oneAuthor.Lastname = "Test4"
	chg, err := authorDao.Upsert(oneAuthor)
	if err != nil {
		t.Error(err)
	}

	t.Log("author modified", chg, oneAuthor)

	oneAuthor, err = authorDao.Get(oneAuthor.Id)
	if err != nil {
		t.Error(err)
	}

	t.Log("author found one modified", oneAuthor)

	err = authorDao.Delete(oneAuthor.Id)
	if err != nil {
		t.Error(err)
	}

	oneAuthor, err = authorDao.Get(oneAuthor.Id)
	if err != nil {
		t.Log("author deleted", err)
	}
}
