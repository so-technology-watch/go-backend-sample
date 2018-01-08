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
	port               = 8020
	logLevel           = "warning"
	db                 = 1
	dbConfigFile       = ""
)

// Main
func main() {
	// Get arguments
	flag.IntVar(&port, "p", port, "webserver port")
	flag.IntVar(&db, "db", db, "database (0 - Redis | 1 - Mock)")
	flag.StringVar(&dbConfigFile, "dbConf", dbConfigFile, "config file path")
	flag.StringVar(&logLevel, "log", logLevel, "log level")

	// Parse arguments
	flag.Parse()

	// Initialisation log
	err := logger.InitLog(logLevel)
	if err != nil {
		logrus.Warn("error setting log level, using debug as default")
	}

	// Get DAO Redis
	taskDAO, err := dao.GetDAO(dao.DBType(db), dbConfigFile)
	if err != nil {
		logrus.WithField("db", db).WithField("dbConf", dbConfigFile).Error("unable to build the required DAO")
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