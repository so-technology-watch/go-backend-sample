package dao

import (
	"errors"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"go-backend-sample/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// compilation time interface check
var _ AlbumDAO = (*AlbumDAOMongo)(nil)

var (
	// ErrInvalidAlbumUUID is used on invalid UUID number for an album
	ErrInvalidAlbumUUID = errors.New("invalid input to UUID")
)

const (
	collectionAlbums = "albums"
	indexAlbum       = "id"
)

// AlbumDAOMongo is the mongo implementation of the AlbumDAO
type AlbumDAOMongo struct {
	session *mgo.Session
}

// NewAlbumDAOMongo creates a new AlbumDAO mongo implementation
func NewAlbumDAOMongo(session *mgo.Session) AlbumDAO {
	// create index
	err := session.DB("").C(collectionAlbums).EnsureIndex(mgo.Index{
		Key:        []string{indexAlbum},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	})

	if err != nil {
		logrus.Error("mongodb connexion error :", err)
		panic(err)
	}

	return &AlbumDAOMongo{
		session: session,
	}
}

// Get returns a album by its id
func (s *AlbumDAOMongo) Get(id string) (*model.Album, error) {
	if _, err := uuid.FromString(id); err != nil {
		return nil, ErrInvalidAlbumUUID
	}

	album := model.Album{}

	session := s.session.Copy()
	defer session.Close()
	c := session.DB("").C(collectionAlbums)
	err := c.Find(bson.M{"id": id}).One(&album)
	if err != nil {
		return nil, err
	}
	return &album, err
}

// GetByAuthor returns all albums by author id in parameter
func (s *AlbumDAOMongo) GetByAuthor(authorId string) ([]model.Album, error) {
	var err error
	var albums []model.Album

	session := s.session.Copy()
	defer session.Close()
	c := session.DB("").C(collectionAlbums)
	err = c.Find(bson.M{"authorId": authorId}).All(&albums)
	if err != nil {
		return nil, err
	}
	return albums, err
}

// GetAll returns all albums
func (s *AlbumDAOMongo) GetAll() ([]model.Album, error) {
	var err error
	var albums []model.Album

	session := s.session.Copy()
	defer session.Close()
	c := session.DB("").C(collectionAlbums)
	err = c.Find(nil).All(&albums)
	if err != nil {
		return nil, err
	}
	return albums, err
}

// Upsert updates or creates an album
func (s *AlbumDAOMongo) Upsert(album *model.Album) (*model.Album, error) {
	if len(album.Id) == 0 {
		album.Id = uuid.NewV4().String()
	}

	session := s.session.Copy()
	defer session.Close()
	c := session.DB("").C(collectionAlbums)
	_, err := c.Upsert(bson.M{"id": album.Id}, album)
	if err != nil {
		return nil, err
	}
	return album, nil
}

// Delete deletes an album by its id
func (s *AlbumDAOMongo) Delete(id string) error {
	if _, err := uuid.FromString(id); err != nil {
		return ErrInvalidAlbumUUID
	}

	session := s.session.Copy()
	defer session.Close()
	c := session.DB("").C(collectionAlbums)
	err := c.Remove(bson.M{"id": id})
	return err
}

// DeleteAll deletes all albums
func (s *AlbumDAOMongo) DeleteAll() error {
	session := s.session.Copy()
	defer session.Close()
	c := session.DB("").C(collectionAlbums)
	_, err := c.RemoveAll(nil)
	return err
}

// Exist check if the album exist
func (s *AlbumDAOMongo) Exist(id string) (bool, error) {
	session := s.session.Copy()
	defer session.Close()
	c := session.DB("").C(collectionAlbums)
	count, err := c.Find(bson.M{"id": id}).Count()
	if err != nil {
		return false, err
	}
	return count == 1, err
}
