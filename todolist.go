package main

import (
	"flag"
	"fmt"
	"net/http"
	"go-backend-sample/dao"
	"go-backend-sample/web"
	"strconv"
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
	taskDAO, err := dao.GetDAO(dao.DBType(db))
	if err != nil {
		fmt.Println(err)
	}

	// New Controller
	taskCtrl := web.NewTaskController(taskDAO)

	// New Router
	router := web.NewRouter(taskCtrl)
	
	fmt.Println("Starting web server on port : " + strconv.Itoa(port))
	http.ListenAndServe(":"+strconv.Itoa(port), router)
}