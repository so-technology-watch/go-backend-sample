package main

import (
	"fmt"
	"net/http"
	"go-backend-sample/dao"
	"go-backend-sample/web"
)

var (
	taskDAO dao.TaskDAO
	taskController *web.TaskController
)

// Main
func main() {
	// Get DAO Mock
	taskDAO = dao.GetDAO()
	
	// New controller
	taskController = web.NewTaskController(taskDAO)

	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/tasks", tasksHandler)

	fmt.Println("Starting web server...")
	http.ListenAndServe(":8020", nil)
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
