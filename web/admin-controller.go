package web

import (
	"go-backend-sample/dao"
	"go-backend-sample/utils"
	"net/http"
)

const (
	prefixAdmin = "/admin"
)

// AdminController is a controller for admin resource
type AdminController struct {
	albumDao  dao.AlbumDAO
	authorDao dao.AuthorDAO
	Routes    []Route
	Prefix    string
}

// NewAdminController creates a new album controller to manage albums & authors
func NewAdminController(albumDAO dao.AlbumDAO, authorDAO dao.AuthorDAO) *AdminController {
	controller := AdminController{
		albumDao:  albumDAO,
		authorDao: authorDAO,
		Prefix:    prefixAdmin,
	}

	var routes []Route
	// DeleteAll
	routes = append(routes, Route{
		Name:        "Delete all albums & authors",
		Method:      http.MethodDelete,
		Pattern:     "/delete",
		HandlerFunc: controller.DeleteAll,
	})

	controller.Routes = routes

	return &controller
}

// DeleteAll deletes all authors and albums with songs
func (ctrl *AdminController) DeleteAll(w http.ResponseWriter, r *http.Request) {
	utils.LogInfo.Println("delete all keys")

	err := ctrl.authorDao.DeleteAll()
	if err != nil {
		utils.LogError.Println(err)
		SendJSONError(w, err.Error(), http.StatusNotFound)
		return
	}

	err = ctrl.albumDao.DeleteAll()
	if err != nil {
		utils.LogError.Println(err)
		SendJSONError(w, err.Error(), http.StatusNotFound)
		return
	}

	SendJSONWithHTTPCode(w, nil, http.StatusNoContent)
}
