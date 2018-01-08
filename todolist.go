package main

import (
	"flag"
	"go-backend-sample/dao"
	"go-backend-sample/web"
	"go-backend-sample/logger"
	"strconv"
	"github.com/urfave/negroni"
	"github.com/sirupsen/logrus"
)

var (
	taskDAO        dao.TaskDAO
	taskController *web.TaskController

	port, db int
	dbFile, logLevel string
)

// Main
func main() {
	// Get arguments
	flag.IntVar(&port, "p", 8020, "webserver port")
	flag.IntVar(&db, "db", 1, "database (0 - Redis | 1 - Mock)")
	flag.StringVar(&dbFile, "dbFile", "", "config file path")
	flag.StringVar(&logLevel, "log", "debug", "log level")

	// Parse arguments
	flag.Parse()

	// Initialisation log
	err := logger.InitLog(logLevel)
	if err != nil {
		logrus.Warn("error setting log level, using debug as default")
	}

	// Get DAO Redis
	taskDAO, err := dao.GetDAO(dao.DBType(db), dbFile)
	if err != nil {
		logrus.WithField("db", db).WithField("dbFile", dbFile).Error("unable to build the required DAO")
	}

	// New webserver
	webServer := negroni.New()

	recovery := negroni.NewRecovery()
	recovery.PrintStack = false
	webServer.Use(recovery)

	// New controller
	taskCtrl := web.NewTaskController(taskDAO)

	// New Router
	router := web.NewRouter(taskCtrl)

	webServer.UseHandler(router)
	webServer.Run(":" + strconv.Itoa(port))
}