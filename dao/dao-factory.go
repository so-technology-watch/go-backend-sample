package dao

import (
	"errors"

	"cloud.google.com/go/datastore"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

// DBType is the database type
type DBType int

const (
	DatastoreDAO DBType = iota
	MockDAO

	AppName   = "Todolist"
	projectId = "go-back-sample"
)

var (
	ErrorDAONotFound = errors.New("unknown DAO type")
)

// GetDAO get TaskDAO
func GetDAO(daoType DBType, dbConfigFile string) (TaskDAO, error) {
	switch daoType {
	case DatastoreDAO:
		clientDatastore := initDatastore()
		return NewTaskDAODatastore(clientDatastore), nil
	case MockDAO:
		return NewTaskDAOMock(), nil
	default:
		return nil, ErrorDAONotFound
	}
}

func initDatastore() *datastore.Client {
	logrus.Println("datastore connexion ")
	ctx := context.Background()

	// connection to the Datastore
	client, err := datastore.NewClient(ctx, projectId)
	if err != nil {
		logrus.Error("datastore connexion error :", err)
		panic(err)
	}

	// verification of connection
	t, err := client.NewTransaction(ctx)
	if err != nil {
		logrus.Error("datastore connexion error :", err)
		panic(err)
	}
	if err := t.Rollback(); err != nil {
		logrus.Error("datastore connexion error :", err)
		panic(err)
	}
	return client
}
