package dao

import (
	"go-redis-sample/model"
)

// AuthorDAO is the DAO interface to work with authors
type AuthorDAO interface {

	// Get returns an author by its id
	Get(id string) (*model.Author, error)

	// GetAll returns all authors
	GetAll() ([]model.Author, error)

	// Upsert updates or creates an author, returns true if updated, false otherwise or on error
	Upsert(author *model.Author) (*model.Author, error)

	// Delete deletes an author by its ID
	Delete(id string) error

	// DeleteAll deletes all authors
	DeleteAll() error

	// Exist checks if the author exist
	Exist(id string) (bool, error)
}
