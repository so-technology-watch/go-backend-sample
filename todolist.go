package main

import (
	"fmt"
	"net/http"
	"go-backend-sample/model"
	"encoding/json"
	"time"
)

// Main
func main() {
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
	task := model.Task{
		Id:          	"1",
		Title:       	"Title Test",
		Description: 	"Description Test",
		Status: 		0,
		CreationDate:	time.Now(),
	}
	
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
