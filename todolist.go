package main

import (
	"fmt"
	"net/http"
	"go-backend-sample/web"
)

// Main
func main() {

	http.HandleFunc("/", welcomeHandler)

	// Tasks Handler
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
		web.GetTask(w, r)
	case "PUT":
		web.UpdateTask(w, r)
	case "POST":
		web.CreateTask(w, r)
	case "DELETE":
		web.DeleteTask(w, r)
	default:
		fmt.Fprintf(w, "Sorry only GET, POST, PUT and DELETE methods are supported.")
	}
}
