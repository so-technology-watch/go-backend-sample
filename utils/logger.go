package utils

import (
	"log"
	"os"
)

const (
	AppName = "bookstore"
)

var (
	LogInfo    *log.Logger
	LogWarning *log.Logger
	LogError   *log.Logger
)

// Initialize the logger
func init() {
	LogInfo = log.New(os.Stdout,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	LogWarning = log.New(os.Stdout,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	LogError = log.New(os.Stderr,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}
