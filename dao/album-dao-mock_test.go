package dao_test

import (
	"testing"
	"go-redis-sample/dao"
	"go-redis-sample/model"
)

func TestAlbumDAOMock(t *testing.T) {

	_, albumDaoMock, err := dao.GetDAO(dao.MockDAO, dao.DBConfigFileName)
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
		Id: 		 "1",
		Title:       "Test",
		Description: "Description Test",
		AuthorId:    "1",
		Songs:       songsToSave,
	}

	albumSaved, err := albumDaoMock.Upsert(&albumToSave)
	if err != nil {
		t.Error(err)
	}

	t.Log("album saved", albumSaved)

	albums, err := albumDaoMock.GetAll()
	if err != nil {
		t.Error(err)
	}
	if len(albums) != 1 {
		t.Errorf("Expected 1 albums, got %d", len(albums))
	}

	oneAlbum, err := albumDaoMock.Get(albumToSave.Id)
	if err != nil {
		t.Error(err)
	}
	// TODO bizarre bizarre
	if &albumToSave != oneAlbum {
		t.Error("Got wrong album by id")
	}

	err = albumDaoMock.Delete(oneAlbum.Id)
	if err != nil {
		t.Error(err)
	}

	oneAlbum, err = albumDaoMock.Get(oneAlbum.Id)
	if err == nil {
		t.Error("Album should have been deleted", oneAlbum)
	}
}
