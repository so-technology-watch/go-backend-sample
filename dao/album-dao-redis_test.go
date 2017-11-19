package dao_test

import (
	"go-redis-sample/dao"
	"go-redis-sample/model"
	"testing"
)

func TestAlbumDAORedis(t *testing.T) {
	authorDao, albumDao, err := dao.GetDAO(dao.RedisDAO, dao.DBConfigFileName)
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

	song1ToSave := model.Song{
		Title:  "Test1",
		Number: "1",
	}
	song2ToSave := model.Song{
		Title:  "Test2",
		Number: "2",
	}

	var songsToSave []model.Song
	songsToSave = append(songsToSave, song1ToSave)
	songsToSave = append(songsToSave, song2ToSave)

	albumToSave := model.Album{
		Title:       "Test",
		Description: "Description Test",
		AuthorId:    authorSaved.Id,
		Songs:       songsToSave,
	}

	albumSaved, err := albumDao.Upsert(&albumToSave)
	if err != nil {
		t.Error(err)
	}

	t.Log("album saved", albumSaved)

	albums, err := albumDao.GetAll()
	if err != nil {
		t.Error(err)
	}

	t.Log("album found all", len(albums))

	oneAlbum, err := albumDao.Get(albumSaved.Id)
	if err != nil {
		t.Error(err)
	}

	t.Log("album found one", oneAlbum)

	oneAlbum.Title = "Use Go(lang)"
	oneAlbum.Description = "Let's build a REST service in Go !"
	oneAlbumChanged, err := albumDao.Upsert(oneAlbum)
	if err != nil {
		t.Error(err)
	}

	t.Log("album modified", oneAlbumChanged, oneAlbum)

	oneAlbum, err = albumDao.Get(oneAlbumChanged.Id)
	if err != nil {
		t.Error(err)
	}

	t.Log("album found one modified", oneAlbum)

	err = albumDao.Delete(oneAlbum.Id)
	if err != nil {
		t.Error(err)
	}

	err = authorDao.Delete(authorSaved.Id)
	if err != nil {
		t.Error(err)
	}

	oneAlbum, err = albumDao.Get(oneAlbum.Id)
	if err != nil {
		t.Log("album deleted", err, oneAlbum)
	}
}
