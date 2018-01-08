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
	taskController = web.NewTaskController(taskDAO)

	http.HandleFunc("/", welcomeHandler)

	// Tasks Handler
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
