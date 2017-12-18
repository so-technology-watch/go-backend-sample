package dao

import (
	"errors"
	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/redis.v5"
	"os"
	"time"
)

type DBType int

type DBConfig struct {
	Url      string
	Port     string
	Password string
	Database string
}

const (
	// RedisDAO is used for Redis implementation of TaskDAO
	RedisDAO DBType = iota
	// MongoDAO is used for Mongo implementation of TaskDAO
	MongoDAO
	// MockDAO is used for mocked implementation of TaskDAO
	MockDAO

	// mongo timeout
	timeout = 5 * time.Second
	// poolSize of mongo connection pool
	poolSize = 35
)

var (
	ErrorDAONotFound = errors.New("unknown DAO type")

	redisLocalConfig = DBConfig{
		Url:      os.Getenv("URL_DB"),
		Password: "",
		Database: "",
		Port:     "6379",
	}

	mongoLocalConfig = DBConfig{
		Url:      os.Getenv("URL_DB"),
		Password: "",
		Database: "todolist",
		Port:     "27017",
	}
)

// GetDAO returns a TaskDAO according to type and params
func GetDAO(daoType DBType, dbConfigFile string) (TaskDAO, error) {
	switch daoType {
	case RedisDAO:
		config := getConfig(RedisDAO, dbConfigFile)
		redisCli := initRedis(config)
		return NewTaskDAORedis(redisCli), nil
	case MongoDAO:
		config := getConfig(MongoDAO, dbConfigFile)
		mongoSession := initMongo(config)
		return NewTaskDAOMongo(mongoSession), nil
	case MockDAO:
		return NewTaskDAOMock(), nil
	default:
		return nil, ErrorDAONotFound
	}
}

// Initialize Redis database
func initRedis(dbConfig DBConfig) *redis.Client {
	logrus.Println("redis connexion " + dbConfig.Url)

	// Connection to the Redis database
	redisCli := redis.NewClient(&redis.Options{
		Addr:     dbConfig.Url + ":" + dbConfig.Port,
		Password: dbConfig.Password,
		DB:       int(RedisDAO),
	})

	// Verification of connection
	ok, err := redisCli.Ping().Result()
	if err != nil {
		logrus.Error("redis connexion error :", err.Error())
		panic(err)
	} else {
		logrus.Println("redis connexion OK :", ok)
	}

	return redisCli
}

// Initialize Mongo database
func initMongo(dbConfig DBConfig) *mgo.Session {
	logrus.Info("mongodb connexion " + dbConfig.Url)

	// Connection to the Mongo database
	mongoSession, err := mgo.DialWithTimeout("mongodb://"+dbConfig.Url+":"+dbConfig.Port+"/"+dbConfig.Database, timeout)
	if err != nil {
		logrus.Error("mongodb connexion error :", err.Error())
		panic(err)
	} else {
		logrus.Info("mongodb connexion OK")
	}

	mongoSession.SetSyncTimeout(timeout)
	mongoSession.SetSocketTimeout(timeout)
	mongoSession.SetMode(mgo.Monotonic, true)
	mongoSession.SetPoolLimit(poolSize)

	return mongoSession
}

func getConfig(daoType DBType, dbConfigFile string) DBConfig {
	var config DBConfig
	if dbConfigFile == "" {
		switch daoType {
		case RedisDAO:
			config = redisLocalConfig
		case MongoDAO:
			config = mongoLocalConfig
		}
	} else {
		if _, err := toml.DecodeFile(dbConfigFile, &config); err != nil {
			logrus.Error("connexion parameters error :", err)
			panic(err)
		}
	}
	return config
}
