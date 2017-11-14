package dao

import (
	"go-redis-sample/model"
)

// AlbumDAO is the DAO interface to work with albums
type AlbumDAO interface {

	// Get returns an album by its id
	Get(id string) (*model.Album, error)

	// GetByAuthor returns all albums by author id in parameter
	GetByAuthor(authorId string) ([]model.Album, error)

	// GetAll returns all albums
	GetAll() ([]model.Album, error)

	// Upsert updates or creates an album, returns true if updated, false otherwise or on error
	Upsert(album *model.Album) (*model.Album, error)

	// Delete deletes an album by its id
	Delete(id string) error

	// DeleteAll deletes all albums
	DeleteAll() error

	// Check if the album exist
	Exist(id string) (bool, error)
}
