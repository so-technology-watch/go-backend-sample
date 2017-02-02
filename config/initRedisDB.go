package config

import (
	"fmt"
	"gopkg.in/redis.v5"
)

var DB *redis.Client

func init() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Connexion BDD OK : " + pong)
	}
	DB = client
}
