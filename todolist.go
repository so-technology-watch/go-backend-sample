package main

import (
	"flag"
	"fmt"
	"go-backend-sample/dao"
	"go-backend-sample/web"
	"strconv"
	"github.com/urfave/negroni"
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

	// Get DAO Redis
	taskDAO, err := dao.GetDAO(dao.DBType(db), dbFile)
	if err != nil {
		fmt.Println(err)
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