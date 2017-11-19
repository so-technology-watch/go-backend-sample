package dao_test

import (
	"github.com/satori/go.uuid"
	"testing"
	"go-redis-sample/dao"
	"go-redis-sample/model"
)

func TestAuthorDAOMongo(t *testing.T) {
	authorDao, _, err := dao.GetDAO(dao.MongoDAO, dao.DBConfigFileName)
	if err != nil {
		t.Error(err)
	}

	toSave := model.Author{
		Id:           uuid.NewV4().String(),
		Firstname:    "Use Go",
		Lastname:     "Let's use the Go programming language in a real project.",
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

	oneAuthor.Firstname = "Use Go(lang)"
	oneAuthor.Lastname = "Let's build a REST service in Go !"
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
