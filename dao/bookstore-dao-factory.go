package dao

import (
	"errors"
	"github.com/BurntSushi/toml"
	"go-redis-sample/utils"
	"gopkg.in/redis.v5"
)

// DBType lists the type of implementation the factory can return
type DBType int

type RedisConfig struct {
	Url      string
	Port     string
	Password string
	Db       string
}

const (
	// DAORedis is used for Redis implementation of AlbumDAO or AuthorDAO
	DAORedis DBType = iota
	// DAOMock is used for mocked implementation of AlbumDAO or AuthorDAO
	DAOMock
)

var (
	ErrorDAONotFound = errors.New("unknown DAO type")
	redisCli         *redis.Client
)

// GetDAO returns an AlbumDAO & an AuthorDAO according to type and params
func GetDAO(daoType DBType) (AuthorDAO, AlbumDAO, error) {
	switch daoType {
	case DAORedis:
		initRedis()
		return NewAuthorDAORedis(redisCli), NewAlbumDAORedis(redisCli), nil
	//case DAOMock:
	//	return NewAuthorDAOMock(), NewAlbumDAOMock(), nil
	default:
		return nil, nil, ErrorDAONotFound
	}
}

// Initialize Redis database
func initRedis() {
	utils.LogInfo.Println("connexion Redis")

	redisConfig := RedisConfig{}
	if _, err := toml.DecodeFile("./config/config.toml", &redisConfig); err != nil {
		utils.LogError.Println("paramètres de connexion Redis error :", err)
		panic(err)
	}

	// Connection to the REDIS database
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// Verification of connection
	ok, err := client.Ping().Result()
	if err != nil {
		utils.LogError.Println("connexion Redis error :", err.Error())
		panic(err)
	} else {
		utils.LogInfo.Println("connexion Redis OK :", ok)
	}
	redisCli = client
}
