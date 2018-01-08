package dao

import (
	"gopkg.in/redis.v5"
	"fmt"
)

type DBType int

const (
	// RedisDAO is used for Redis implementation of TaskDAO
	RedisDAO DBType = iota
	// MockDAO is used for mocked implementation of TaskDAO
	MockDAO
)

// GetDAO returns a TaskDAO according to type and params
func GetDAO(daoType DBType) TaskDAO {
	switch daoType {
	case RedisDAO:
		redisCli := initRedis()
		return NewTaskDAORedis(redisCli)
	case MockDAO:
		return NewTaskDAOMock()
	default:
		return nil
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
