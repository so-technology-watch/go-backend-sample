package config

import (
	"encoding/json"
	"os"
	"fmt"
	"gopkg.in/redis.v5"
)

// Structure of database configuration
type DBConfiguration struct {
	Url    string
	Port   string
	Password string
	Db int
}

var DB *redis.Client

// Initialize REDIS database
func init() {
	// Collect configuration
	var dbConfiguration DBConfiguration
	file, _ := os.Open("dbConfig.json")
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&dbConfiguration)
	if err != nil {
		fmt.Println("File config error :", err)
	}

	// Connection to the REDIS database
	client := redis.NewClient(&redis.Options{
		Addr:     dbConfiguration.Url + ":" + dbConfiguration.Port,
		Password: dbConfiguration.Password,
		DB:       dbConfiguration.Db,
	})

	// Verification of connection
	ok, err := client.Ping().Result()
	if err != nil {
		fmt.Println("Connexion DB error :", err.Error())
	} else {
		fmt.Println("Connexion DB OK :", ok)
	}
	DB = client
}
