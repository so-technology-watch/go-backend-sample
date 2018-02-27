package web

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	requestHeaderAccessControlRequestMethod = "Access-Control-Request-Method"
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

	router.Methods(http.MethodOptions).HandlerFunc(preflightTasks)

	addTaskRoutes(taskCtrl, router)

	return &router
}

// addTaskRoutes add the routes of tasks
func addTaskRoutes(taskCtrl *TaskController, router Router) {
	for _, route := range taskCtrl.Routes {
		router.
			Methods(route.Method).
			Path(taskCtrl.Prefix + route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
}

// preflightTasks handles the preflight requests
func preflightTasks(w http.ResponseWriter, r *http.Request) {
	logrus.Println("preflight request handled")
	if r.Header.Get(requestHeaderAccessControlRequestMethod) == http.MethodDelete || r.Header.Get(requestHeaderAccessControlRequestMethod) == http.MethodPut || r.Header.Get(requestHeaderAccessControlRequestMethod) == http.MethodPost {
		SendJSONOk(w, nil)
	} else {
		SendJSONError(w, "Unsupported method", http.StatusUnauthorized)
	}
}
