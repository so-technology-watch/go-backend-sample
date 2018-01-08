package dao

import (
	"errors"
	"fmt"
	"gopkg.in/redis.v5"
)

type DBType int

const (
	// RedisDAO is used for Redis implementation of TaskDAO
	RedisDAO DBType = iota
	// MockDAO is used for mocked implementation of TaskDAO
	MockDAO
)

// GetDAO returns a TaskDAO according to type and params
func GetDAO(daoType DBType) (TaskDAO, error) {
	switch daoType {
	case RedisDAO:
		redisCli := initRedis()
		return NewTaskDAORedis(redisCli), nil
	case MockDAO:
		return NewTaskDAOMock(), nil
	default:
		return nil, errors.New("unknown DAO type")
	}
}

// Initialize Redis database
func initRedis() *redis.Client {
	// Connection to the Redis database
	redisCli := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// Verification of connection
	ok, err := redisCli.Ping().Result()
	if err != nil {
		panic(err)
	} else {
		fmt.Println("redis connexion OK :", ok)
	}

	return redisCli
}
