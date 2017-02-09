package config

import (
	"encoding/json"
	"os"
	"fmt"
	"gopkg.in/redis.v5"
)

type DBConfiguration struct {
	Url    string
	Port   string
	Password string
	Db int
}

var DB *redis.Client

func init() {
	var dbConfiguration DBConfiguration
	file, _ := os.Open("dbConfig.json")
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&dbConfiguration)
	if err != nil {
		fmt.Println("File config error :", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     dbConfiguration.Url + ":" + dbConfiguration.Port,
		Password: dbConfiguration.Password,
		DB:       dbConfiguration.Db,
	})

	pong, err := client.Ping().Result()
	if err != nil {
		fmt.Println("Connexion DB error :", err.Error())
	} else {
		fmt.Println("Connexion DB OK :", pong)
	}
	DB = client
}
