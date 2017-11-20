package dao

import (
	"errors"
	"github.com/satori/go.uuid"
	"go-backend-sample/model"
	"go-backend-sample/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// compilation time interface check
var _ AuthorDAO = (*AuthorDAOMongo)(nil)

var (
	// ErrInvalidUUIDAuthor is used on invalid UUID number for an author
	ErrInvalidUUIDAuthor = errors.New("invalid input to UUID")
)

const (
	collectionAuthors = "authors"
	indexAuthor       = "id"
)

// AuthorDAOMongo is the mongo implementation of the AuthorDAO
type AuthorDAOMongo struct {
	session *mgo.Session
}

// NewAuthorDAOMongo creates a new AuthorDAO mongo implementation
func NewAuthorDAOMongo(session *mgo.Session) AuthorDAO {
	// create index
	err := session.DB("").C(collectionAuthors).EnsureIndex(mgo.Index{
		Key:        []string{indexAuthor},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	})

	if err != nil {
		utils.LogError.Println("mongodb connexion error :", err)
		panic(err)
	}

	return &AuthorDAOMongo{
		session: session,
	}
}

// Get returns a author by its id
func (s *AuthorDAOMongo) Get(id string) (*model.Author, error) {
	if _, err := uuid.FromString(id); err != nil {
		return nil, ErrInvalidUUIDAuthor
	}

	author := model.Author{}

	session := s.session.Copy()
	defer session.Close()
	c := session.DB("").C(collectionAuthors)
	err := c.Find(bson.M{"id": id}).One(&author)
	if err != nil {
		return nil, err
	}
	return &author, err
}

// GetAll returns all authors
func (s *AuthorDAOMongo) GetAll() ([]model.Author, error) {
	var err error
	var authors []model.Author

	session := s.session.Copy()
	defer session.Close()
	c := session.DB("").C(collectionAuthors)
	err = c.Find(nil).All(&authors)
	if err != nil {
		return nil, err
	}
	return authors, err
}

// Upsert updates or creates an author
func (s *AuthorDAOMongo) Upsert(author *model.Author) (*model.Author, error) {
	if len(author.Id) == 0 {
		author.Id = uuid.NewV4().String()
	}

	session := s.session.Copy()
	defer session.Close()
	c := session.DB("").C(collectionAuthors)
	_, err := c.Upsert(bson.M{"id": author.Id}, author)
	if err != nil {
		return nil, err
	}
	return author, nil
}

// Delete deletes an author by its id
func (s *AuthorDAOMongo) Delete(id string) error {
	if _, err := uuid.FromString(id); err != nil {
		return ErrInvalidUUIDAuthor
	}

	session := s.session.Copy()
	defer session.Close()
	c := session.DB("").C(collectionAuthors)
	err := c.Remove(bson.M{"id": id})
	return err
}

// DeleteAll deletes all authors
func (s *AuthorDAOMongo) DeleteAll() error {
	session := s.session.Copy()
	defer session.Close()
	c := session.DB("").C(collectionAuthors)
	_, err := c.RemoveAll(nil)
	return err
}

// Exist check if the author exist
func (s *AuthorDAOMongo) Exist(id string) (bool, error) {
	session := s.session.Copy()
	defer session.Close()
	c := session.DB("").C(collectionAuthors)
	count, err := c.Find(bson.M{"id": id}).Count()
	if err != nil {
		return false, err
	}
	return count == 1, err
}
