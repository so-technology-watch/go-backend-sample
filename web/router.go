package web

import (
	"github.com/gorilla/mux"
	"net/http"
)

// Router is the struct use for routing
type Router struct {
	*mux.Router
}

// Route is a structure of Route
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// NewRouter creates a new router instance
func NewRouter(taskCtrl *TaskController) *Router {
	router := Router{mux.NewRouter()}
	router.NotFoundHandler = NotFoundHandler()
	router.StrictSlash(false)

	AddTaskRoutes(taskCtrl, router)
	return &router
}

func AddTaskRoutes(taskCtrl *TaskController, router Router) {
	for _, route := range taskCtrl.Routes {
		router.
			Methods(route.Method).
			Path(taskCtrl.Prefix + route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
}
