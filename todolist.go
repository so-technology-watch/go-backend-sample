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

	// Get DAO Redis
	taskDAO, err := dao.GetDAO(dao.DBType(db))
	if err != nil {
		fmt.Println(err)
	}

	// New controller
	taskController = web.NewTaskController(taskDAO)

	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/tasks", tasksHandler)

	fmt.Println("Starting web server on port : " + strconv.Itoa(port))
	http.ListenAndServe(":"+strconv.Itoa(port), nil)
}

// Welcome handler
func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome !")
}

// Tasks handler
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		taskController.GetTask(w, r)
	case "PUT":
		taskController.UpdateTask(w, r)
	case "POST":
		taskController.CreateTask(w, r)
	case "DELETE":
		taskController.DeleteTask(w, r)
	default:
		fmt.Fprintf(w, "Sorry only GET, POST, PUT and DELETE methods are supported.")
	}
}
