package config

import (
	"io"
	"log"
	"os"
)

var (
	LogInfo    *log.Logger
	LogWarning *log.Logger
	LogError   *log.Logger
)

// Initialize the logger
func init() {
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file", os.Stdout, ":", err)
	}

	multiOut := io.MultiWriter(logFile, os.Stdout)
	multiErr := io.MultiWriter(logFile, os.Stderr)

	LogInfo = log.New(multiOut,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	LogWarning = log.New(multiOut,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	LogError = log.New(multiErr,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}
