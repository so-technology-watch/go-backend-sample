package web

import (
	"github.com/sirupsen/logrus"
	"go-backend-sample/dao"
	"go-backend-sample/model"
	"net/http"
	"time"
)

const (
	prefixTask = "/tasks"
)

// TaskController define the controller for tasks
type TaskController struct {
	taskDao dao.TaskDAO
	Routes  []Route
	Prefix  string
}

// NewTaskController build the controller for tasks
func NewTaskController(taskDAO dao.TaskDAO) *TaskController {
	controller := TaskController{
		taskDao: taskDAO,
		Prefix:  prefixTask,
	}

	var routes []Route
	// GetAll
	routes = append(routes, Route{
		Name:        "Get all tasks",
		Method:      http.MethodGet,
		Pattern:     "",
		HandlerFunc: controller.GetTasks,
	})
	// Get
	routes = append(routes, Route{
		Name:        "Get one task",
		Method:      http.MethodGet,
		Pattern:     "/{id}",
		HandlerFunc: controller.GetTask,
	})
	// Create
	routes = append(routes, Route{
		Name:        "Create an task",
		Method:      http.MethodPost,
		Pattern:     "",
		HandlerFunc: controller.CreateTask,
	})
	// Update
	routes = append(routes, Route{
		Name:        "Update an task",
		Method:      http.MethodPut,
		Pattern:     "/{id}",
		HandlerFunc: controller.UpdateTask,
	})
	// Delete
	routes = append(routes, Route{
		Name:        "Delete an task",
		Method:      http.MethodDelete,
		Pattern:     "/{id}",
		HandlerFunc: controller.DeleteTask,
	})
	// Delete All
	routes = append(routes, Route{
		Name:        "Delete all tasks",
		Method:      http.MethodDelete,
		Pattern:     "",
		HandlerFunc: controller.DeleteTasks,
	})

	controller.Routes = routes

	return &controller
}

// GetTasks retrieve all tasks
func (ctrl *TaskController) GetTasks(w http.ResponseWriter, r *http.Request) {
	logrus.Println("list tasks")

	tasks, err := ctrl.taskDao.GetAll()
	if err != nil {
		logrus.Error(err)
		SendJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	SendJSONOk(w, tasks)
}

// GetTask retrieve a task by its id
func (ctrl *TaskController) GetTask(w http.ResponseWriter, r *http.Request) {
	taskId := ParamAsString("id", r)
	logrus.Println("task : ", taskId)

	task, err := ctrl.taskDao.Get(taskId)
	if err != nil {
		logrus.Error(err)
		SendJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logrus.Println("task : ", task)
	SendJSONOk(w, task)
}

// CreateTask create a task
func (ctrl *TaskController) CreateTask(w http.ResponseWriter, r *http.Request) {
	task := &model.Task{}
	logrus.Println(r.Body)
	err := GetJSONContent(task, r)
	if err != nil {
		logrus.Error(err)
		SendJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	logrus.Println("create task")

	task.CreationDate = time.Now()
	task.ModificationDate = time.Now()
	task.Status = 0

	task, err = ctrl.taskDao.Upsert(task)
	if err != nil {
		logrus.Error(err)
		SendJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logrus.Println("task : ", task)
	SendJSONWithHTTPCode(w, task, http.StatusCreated)
}

// UpdateTask update a task by its id
func (ctrl *TaskController) UpdateTask(w http.ResponseWriter, r *http.Request) {
	task := &model.Task{}
	err := GetJSONContent(task, r)
	if err != nil {
		logrus.Error(err)
		SendJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	logrus.Println("update task : ", task.Id)

	task.ModificationDate = time.Now()

	taskExist, err := ctrl.taskDao.Exist(task.Id)
	if err != nil {
		logrus.Error(err)
		SendJSONError(w, err.Error(), http.StatusNotFound)
		return
	} else if taskExist == false {
		SendJSONError(w, "task not found", http.StatusNotFound)
		return
	}

	task, err = ctrl.taskDao.Upsert(task)
	if err != nil {
		logrus.Error(err)
		SendJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logrus.Println("task : ", task)
	SendJSONOk(w, task)
}

// DeleteTask delete a task by its id
func (ctrl *TaskController) DeleteTask(w http.ResponseWriter, r *http.Request) {
	taskId := ParamAsString("id", r)
	logrus.Println("delete task : ", taskId)

	err := ctrl.taskDao.Delete(taskId)
	if err != nil {
		logrus.Error(err)
		SendJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logrus.Println("deleted task : ", taskId)
	SendJSONWithHTTPCode(w, true, http.StatusNoContent)
}

// DeleteTasks delete all tasks
func (ctrl *TaskController) DeleteTasks(w http.ResponseWriter, r *http.Request) {
	logrus.Println("delete all tasks")

	err := ctrl.taskDao.DeleteAll()
	if err != nil {
		logrus.Error(err)
		SendJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logrus.Println("all tasks deleted")
	SendJSONWithHTTPCode(w, true, http.StatusNoContent)
}
