package web

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Router define the router
type Router struct {
	*mux.Router
}

// Route define a route
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// NewRouter build a router and add routes
func NewRouter(taskCtrl *TaskController) *Router {
	router := Router{mux.NewRouter()}
	router.NotFoundHandler = NotFoundHandler()
	router.StrictSlash(false)

	AddTaskRoutes(taskCtrl, router)
	router.Methods("OPTIONS").HandlerFunc(PreflightTasks)

	return &router
}

// AddTaskRoutes add the routes of tasks
func AddTaskRoutes(taskCtrl *TaskController, router Router) {
	for _, route := range taskCtrl.Routes {
		router.
			Methods(route.Method).
			Path(taskCtrl.Prefix + route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
}
