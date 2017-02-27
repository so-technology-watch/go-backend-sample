package config

import (
	"io"
	"log"
	"os"
	"time"
)

var (
	LogInfo    *log.Logger
	LogWarning *log.Logger
	LogError   *log.Logger
)

// Initialize the logger
func init() {
	// Verification if folder "logs" exist
	logsFolder := "logs"
	if _, err := os.Stat(logsFolder); os.IsNotExist(err) {
		err := os.Mkdir(logsFolder, os.ModeDir)
		if err != nil {
			log.Fatalln("Failed to create logs folder", os.Stdout, ":", err)
		}
	}

	// Verification if log file exist
	current := time.Now()
	logFileName := logsFolder + "/go-redis-sample-" + current.Format("02-01-2006") + ".log"
	var logFile *os.File
	if _, err := os.Stat(logFileName); os.IsNotExist(err) {
		logFile, err = os.Create(logFileName)
		if err != nil {
			log.Fatalln("Failed to create log file", os.Stdout, ":", err)
		}
	} else {
		logFile, err = os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalln("Failed to open log file", os.Stdout, ":", err)
		}
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
