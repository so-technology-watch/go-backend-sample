package dao

import (
	"errors"
	"github.com/BurntSushi/toml"
	"go-redis-sample/utils"
	"gopkg.in/redis.v5"
	"gopkg.in/mgo.v2"
	"time"
)

// DBType lists the type of implementation the factory can return
type DBType int

type DefaultDBsConfig struct {
	Redis DBConfig
	Mongo DBConfig
}

type DBConfig struct {
	Url      string
	Port     string
	Password string
	Database string
}

const (
	// RedisDAO is used for Redis implementation of AlbumDAO or AuthorDAO
	RedisDAO DBType = iota
	// MongoDAO is used for Mongo implementation of AlbumDAO or AuthorDAO
	MongoDAO
	// MockDAO is used for mocked implementation of AlbumDAO or AuthorDAO
	MockDAO

	// mongo timeout
	timeout = 5 * time.Second
	// poolSize of mongo connection pool
	poolSize = 35

	dbConfigFileName = "default-dbs-config.toml"
)

var (
	ErrorDAONotFound = errors.New("unknown DAO type")
)

// GetDAO returns an AlbumDAO & an AuthorDAO according to type and params
func GetDAO(daoType DBType, dbConfigFile string) (AuthorDAO, AlbumDAO, error) {
	switch daoType {
	case RedisDAO:
		config := getConfig(RedisDAO, dbConfigFile)
		redisCli := initRedis(config)
		return NewAuthorDAORedis(redisCli), NewAlbumDAORedis(redisCli), nil
	case MongoDAO:
		config := getConfig(MongoDAO, dbConfigFile)
		mongoSession := initMongo(config)
		return NewAuthorDAOMongo(mongoSession), NewAlbumDAOMongo(mongoSession), nil
	case MockDAO:
		return NewAuthorDAOMock(), NewAlbumDAOMock(), nil
	default:
		return nil, nil, ErrorDAONotFound
	}
}

// Initialize Redis database
func initRedis(dbConfig DBConfig) *redis.Client {
	utils.LogInfo.Println("connexion Redis")

	utils.LogInfo.Println("redis " + dbConfig.Url)

	// Connection to the Redis database
	redisCli := redis.NewClient(&redis.Options{
		Addr:     dbConfig.Url + ":" + dbConfig.Port,
		Password: dbConfig.Password,
		DB:       int(RedisDAO),
	})

	// Verification of connection
	ok, err := redisCli.Ping().Result()
	if err != nil {
		utils.LogError.Println("redis connexion error :", err.Error())
		panic(err)
	} else {
		utils.LogInfo.Println("redis connexion OK :", ok)
	}

	return redisCli
}

func initMongo(dbConfig DBConfig) *mgo.Session {
	utils.LogInfo.Println("mongodb connexion")

	utils.LogInfo.Println("mongodb " + dbConfig.Url)

	// Connection to the Mongo database
	mongoSession, err := mgo.DialWithTimeout("mongodb://" + dbConfig.Url + ":" + dbConfig.Port + "/" + dbConfig.Database, timeout)
	if err != nil {
		utils.LogError.Println("mongodb connexion error :", err.Error())
		panic(err)
	} else {
		utils.LogInfo.Println("mongodb connexion OK")
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
		var defaultConfig DefaultDBsConfig
		if _, err := toml.DecodeFile(dbConfigFileName, &defaultConfig); err != nil {
			utils.LogError.Println("connexion parameters error :", err)
			panic(err)
		}
		switch daoType {
		case RedisDAO:
			config = defaultConfig.Redis
		case MongoDAO:
			config = defaultConfig.Mongo
		}
	} else {
		if _, err := toml.DecodeFile(dbConfigFile, &config); err != nil {
			utils.LogError.Println("connexion parameters error :", err)
			panic(err)
		}
	}
	return config
}