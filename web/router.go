package web

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Router struct {
	*mux.Router
}

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

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
