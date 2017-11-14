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
func NewRouter(authorCtrl *AuthorController, albumCtrl *AlbumController, adminCtrl *AdminController) *Router {
	router := Router{mux.NewRouter()}
	router.NotFoundHandler = NotFoundHandler()
	router.StrictSlash(false)

	AddAuthorRoutes(authorCtrl, router)
	AddAlbumRoutes(albumCtrl, router)
	AddAdminRoutes(adminCtrl, router)

	return &router
}

func AddAuthorRoutes(authorCtrl *AuthorController, router Router) {
	for _, route := range authorCtrl.Routes {
		router.
			Methods(route.Method).
			Path(authorCtrl.Prefix + route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
}

func AddAlbumRoutes(albumCtrl *AlbumController, router Router) {
	for _, route := range albumCtrl.Routes {
		router.
			Methods(route.Method).
			Path(albumCtrl.Prefix + route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
}

func AddAdminRoutes(adminCtrl *AdminController, router Router) {
	for _, route := range adminCtrl.Routes {
		router.
			Methods(route.Method).
			Path(adminCtrl.Prefix + route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
}
