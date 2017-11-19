package dao_test

import (
	"go-redis-sample/dao"
	"go-redis-sample/model"
	"testing"
)

func TestAuthorDAORedis(t *testing.T) {
	authorDao, _, err := dao.GetDAO(dao.RedisDAO, dao.DBConfigFileName)
	if err != nil {
		t.Error(err)
	}

	authorToSave := model.Author{
		Lastname:  "Lastname",
		Firstname: "Firstname",
	}

	authorSaved, err := authorDao.Upsert(&authorToSave)
	if err != nil {
		t.Error(err)
	}

	t.Log("author saved", authorSaved)

	authors, err := authorDao.GetAll()
	if err != nil {
		t.Error(err)
	}

	t.Log("author found all", authors[0])

	oneAuthor, err := authorDao.Get(authorSaved.Id)
	if err != nil {
		t.Error(err)
	}

	t.Log("author found one", oneAuthor)

	oneAuthor.Lastname = "Lastname2"
	oneAuthor.Firstname = "Firstname2"
	oneAuthorChanged, err := authorDao.Upsert(oneAuthor)
	if err != nil {
		t.Error(err)
	}

	t.Log("author modified", oneAuthorChanged, oneAuthor)

	oneAuthor, err = authorDao.Get(oneAuthorChanged.Id)
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
		t.Log("author deleted", err, oneAuthor)
	}
}
