package web

import (
	"encoding/json"
	"net/http"
	"go-backend-sample/dao"
	"go-backend-sample/model"
	"time"
)

// TaskController is a controller for tasks resources
type TaskController struct {
	taskDao dao.TaskDAO
}

// NewTaskController creates a new task controller to manage tasks
func NewTaskController(taskDAO dao.TaskDAO) *TaskController {
	controller := TaskController{
		taskDao: taskDAO,
	}

	return &controller
}

// Get retrieve a task by its id
func (ctrl *TaskController) GetTask(w http.ResponseWriter, r *http.Request) {
	taskId := r.URL.Query().Get("id")

	task, err := ctrl.taskDao.Get(taskId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Update update a task by its id
func (ctrl *TaskController) UpdateTask(w http.ResponseWriter, r *http.Request) {
	task := &model.Task{}
	err := json.NewDecoder(r.Body).Decode(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	taskExist, err := ctrl.taskDao.Exist(task.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	} else if taskExist == false {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	task.ModificationDate = time.Now()

	task, err = ctrl.taskDao.Upsert(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Create create a task
func (ctrl *TaskController) CreateTask(w http.ResponseWriter, r *http.Request) {
	task := &model.Task{}
	err := json.NewDecoder(r.Body).Decode(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	task.CreationDate = time.Now()
	task.Status = 0

	task, err = ctrl.taskDao.Upsert(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Delete delete a task by its id
func (ctrl *TaskController) DeleteTask(w http.ResponseWriter, r *http.Request) {
	taskId := r.URL.Query().Get("id")

	err := ctrl.taskDao.Delete(taskId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNoContent)
	if err := json.NewEncoder(w).Encode(true); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
