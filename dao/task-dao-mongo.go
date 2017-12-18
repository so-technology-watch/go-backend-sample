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
var _ TaskDAO = (*TaskDAOMongo)(nil)

var (
	// ErrInvalidUUIDTask is used on invalid UUID number for a task
	ErrInvalidUUIDTask = errors.New("invalid input to UUID")
)

const (
	collectionTasks = "tasks"
	indexTask       = "id"
)

// TaskDAOMongo is the mongo implementation of the TaskDAO
type TaskDAOMongo struct {
	session *mgo.Session
}

// NewTaskDAOMongo creates a new TaskDAO mongo implementation
func NewTaskDAOMongo(session *mgo.Session) TaskDAO {
	// create index
	err := session.DB("").C(collectionTasks).EnsureIndex(mgo.Index{
		Key:        []string{indexTask},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	})

	if err != nil {
		logrus.Error("mongodb connexion error :", err)
		panic(err)
	}

	return &TaskDAOMongo{
		session: session,
	}
}

// Get returns a task by its id
func (s *TaskDAOMongo) Get(id string) (*model.Task, error) {
	if _, err := uuid.FromString(id); err != nil {
		return nil, ErrInvalidUUIDTask
	}

	task := model.Task{}

	session := s.session.Copy()
	defer session.Close()
	c := session.DB("").C(collectionTasks)
	err := c.Find(bson.M{"id": id}).One(&task)
	if err != nil {
		return nil, err
	}
	return &task, err
}

// GetAll returns all tasks
func (s *TaskDAOMongo) GetAll() ([]model.Task, error) {
	var err error
	var tasks []model.Task

	session := s.session.Copy()
	defer session.Close()
	c := session.DB("").C(collectionTasks)
	err = c.Find(nil).All(&tasks)
	if err != nil {
		return nil, err
	}
	return tasks, err
}

// Upsert updates or creates a task
func (s *TaskDAOMongo) Upsert(task *model.Task) (*model.Task, error) {
	if len(task.Id) == 0 {
		task.Id = uuid.NewV4().String()
	}

	session := s.session.Copy()
	defer session.Close()
	c := session.DB("").C(collectionTasks)
	_, err := c.Upsert(bson.M{"id": task.Id}, task)
	if err != nil {
		return nil, err
	}
	return task, nil
}

// Delete deletes a task by its id
func (s *TaskDAOMongo) Delete(id string) error {
	if _, err := uuid.FromString(id); err != nil {
		return ErrInvalidUUIDTask
	}

	session := s.session.Copy()
	defer session.Close()
	c := session.DB("").C(collectionTasks)
	err := c.Remove(bson.M{"id": id})
	return err
}

// DeleteAll deletes all tasks
func (s *TaskDAOMongo) DeleteAll() error {
	session := s.session.Copy()
	defer session.Close()
	c := session.DB("").C(collectionTasks)
	_, err := c.RemoveAll(nil)
	return err
}

// Exist check if the task exist
func (s *TaskDAOMongo) Exist(id string) (bool, error) {
	session := s.session.Copy()
	defer session.Close()
	c := session.DB("").C(collectionTasks)
	count, err := c.Find(bson.M{"id": id}).Count()
	if err != nil {
		return false, err
	}
	return count == 1, err
}
