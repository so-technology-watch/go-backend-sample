package web

import (
	"encoding/json"
	"go-backend-sample/model"
	"net/http"
)

// Get retrieve a task by its id
func GetTask(w http.ResponseWriter, r *http.Request) {
	taskId := r.URL.Query().Get("id")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(taskId); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Update update a task by its id
func UpdateTask(w http.ResponseWriter, r *http.Request) {
	task := &model.Task{}
	err := json.NewDecoder(r.Body).Decode(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Create create a task
func CreateTask(w http.ResponseWriter, r *http.Request) {
	task := &model.Task{}
	err := json.NewDecoder(r.Body).Decode(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Delete delete a task by its id
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	taskId := r.URL.Query().Get("id")

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNoContent)
	if err := json.NewEncoder(w).Encode(taskId); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
