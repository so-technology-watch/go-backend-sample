package dao

import (
	"errors"
	"go-redis-sample/model"
	"github.com/satori/go.uuid"
)

// compilation time interface check
var _ AlbumDAO = (*AlbumDAOMock)(nil)

var MockSong1 = model.Song{
	Title:  "Test1",
	Number: "1",
}
var MockSong2 = model.Song{
	Title:  "Test2",
	Number: "2",
}
// MockedAlbum is the album returned by this mocked interface
var MockedAlbum = model.Album{
	Id:         "1",
	Title:    	"TestMock",
	Description:"TestMock",
	AuthorId:	"1",
}

// AlbumDAOMock is the mocked implementation of the AlbumDAOMock
type AlbumDAOMock struct {
	storage map[string]*model.Album
}

// NewAlbumDAOMock creates a new AlbumDAO with a mocked implementation
func NewAlbumDAOMock() AlbumDAO {
	daoMock := &AlbumDAOMock{
		storage: make(map[string]*model.Album),
	}

	// Adds some fake data
	var mockSongs []model.Song
	mockSongs = append(mockSongs, MockSong1)
	mockSongs = append(mockSongs, MockSong2)
	MockedAlbum.Songs = mockSongs
	daoMock.Upsert(&MockedAlbum)

	return daoMock
}

// Get returns an album by its id
func (s *AlbumDAOMock) Get(id string) (*model.Album, error) {
	album, ok := s.storage[id]
	if !ok {
		return nil, errors.New("album not found with id " + id)
	}
	return album, nil
}

// GetByAuthor returns all albums by author id
func (s *AlbumDAOMock) GetByAuthor(authorId string) ([]model.Album, error) {
	var albums []model.Album
	for albumId := range s.storage {
		album := s.storage[albumId]
		if album.AuthorId == authorId {
			albums = append(albums, *album)
		}
	}
	return albums, nil
}

// GetAll returns all albums
func (s *AlbumDAOMock) GetAll() ([]model.Album, error) {
	var albums []model.Album
	for albumId := range s.storage {
		album := s.storage[albumId]
		albums = append(albums, *album)
	}
	return albums, nil
}

// Upsert updates or creates an album
func (s *AlbumDAOMock) Upsert(album *model.Album) (*model.Album, error) {
	if album.Id == "" {
		album.Id = uuid.NewV4().String()
	}
	s.save(album)
	return album, nil
}

// Delete deletes an album by its id
func (s *AlbumDAOMock) Delete(id string) error {
	delete(s.storage, id)
	return nil
}

// DeleteAll deletes all albums
func (s *AlbumDAOMock) DeleteAll() error {
	for albumId := range s.storage {
		delete(s.storage, albumId)
	}
	return nil
}

// Exist checks if the album exist
func (s *AlbumDAOMock) Exist(id string) (bool, error) {
	if s.storage[id] != nil {
		return true, nil
	}
	return false, errors.New("album not found with id " + id)
}

// save saves the album
func (s *AlbumDAOMock) save(album *model.Album) *model.Album {
	s.storage[album.Id] = album
	return album
}

// get return an album by its id
func (s *AlbumDAOMock) get(id string) *model.Album {
	album := s.storage[id]
	return album
}

