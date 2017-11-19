package dao

import (
	"errors"
	"go-backend-sample/model"
	"github.com/satori/go.uuid"
)

// compilation time interface check
var _ AuthorDAO = (*AuthorDAOMock)(nil)

// MockedAuthor is the author returned by this mocked interface
var MockedAuthor = model.Author{
	Id:           "1",
	Firstname:    "TestMock",
	Lastname:  	  "TestMock",
}

// AuthorDAOMock is the mocked implementation of the AuthorDAOMock
type AuthorDAOMock struct {
	storage map[string]*model.Author
}

// NewAuthorDAOMock creates a new AuthorDAO with a mocked implementation
func NewAuthorDAOMock() AuthorDAO {
	daoMock := &AuthorDAOMock{
		storage: make(map[string]*model.Author),
	}

	// Adds some fake data
	daoMock.Upsert(&MockedAuthor)

	return daoMock
}

// Get returns an author by its id
func (s *AuthorDAOMock) Get(id string) (*model.Author, error) {
	author, ok := s.storage[id]
	if !ok {
		return nil, errors.New("author not found with id " + id)
	}
	return author, nil
}

// GetAll returns all authors
func (s *AuthorDAOMock) GetAll() ([]model.Author, error) {
	var authors []model.Author
	for authorId := range s.storage {
		author := s.storage[authorId]
		authors = append(authors, *author)
	}
	return authors, nil
}

// Upsert updates or creates an author
func (s *AuthorDAOMock) Upsert(author *model.Author) (*model.Author, error) {
	if author.Id == "" {
		author.Id = uuid.NewV4().String()
	}
	s.save(author)
	return author, nil
}

// Delete deletes an author by its id
func (s *AuthorDAOMock) Delete(id string) error {
	delete(s.storage, id)
	return nil
}

// DeleteAll deletes all authors
func (s *AuthorDAOMock) DeleteAll() error {
	for authorId := range s.storage {
		delete(s.storage, authorId)
	}
	return nil
}

// Exist checks if the author exist
func (s *AuthorDAOMock) Exist(id string) (bool, error) {
	if s.storage[id] != nil {
		return true, nil
	}
	return false, errors.New("author not found with id " + id)
}

// save saves the author
func (s *AuthorDAOMock) save(author *model.Author) *model.Author {
	s.storage[author.Id] = author
	return author
}

// get return an author by its id
func (s *AuthorDAOMock) get(id string) *model.Author {
	author := s.storage[id]
	return author
}

