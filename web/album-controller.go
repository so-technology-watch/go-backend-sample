package web

import (
	"go-backend-sample/dao"
	"go-backend-sample/model"
	"go-backend-sample/utils"
	"net/http"
)

const (
	prefixAlbum = "/albums"
)

// AlbumController is a controller for albums resource
type AlbumController struct {
	albumDao  dao.AlbumDAO
	authorDao dao.AuthorDAO
	Routes    []Route
	Prefix    string
}

// NewAlbumController creates a new album controller to manage albums
func NewAlbumController(albumDAO dao.AlbumDAO, authorDAO dao.AuthorDAO) *AlbumController {
	controller := AlbumController{
		albumDao:  albumDAO,
		authorDao: authorDAO,
		Prefix:    prefixAlbum,
	}

	var routes []Route
	// GetAll
	routes = append(routes, Route{
		Name:        "Get all albums",
		Method:      http.MethodGet,
		Pattern:     "",
		HandlerFunc: controller.GetAlbums,
	})
	// GetByAuthor
	routes = append(routes, Route{
		Name:        "Get albums by author",
		Method:      http.MethodGet,
		Pattern:     "/author/{authorId}",
		HandlerFunc: controller.GetAlbumsByAuthor,
	})
	// Get
	routes = append(routes, Route{
		Name:        "Get one album",
		Method:      http.MethodGet,
		Pattern:     "/{id}",
		HandlerFunc: controller.GetAlbum,
	})
	// Create
	routes = append(routes, Route{
		Name:        "Create an album",
		Method:      http.MethodPost,
		Pattern:     "",
		HandlerFunc: controller.CreateAlbum,
	})
	// Update
	routes = append(routes, Route{
		Name:        "Update an album",
		Method:      http.MethodPut,
		Pattern:     "/{id}",
		HandlerFunc: controller.UpdateAlbum,
	})
	// Delete
	routes = append(routes, Route{
		Name:        "Delete an album",
		Method:      http.MethodDelete,
		Pattern:     "/{id}",
		HandlerFunc: controller.DeleteAlbum,
	})

	controller.Routes = routes

	return &controller
}

// GetAlbums retrieve all albums
func (ctrl *AlbumController) GetAlbums(w http.ResponseWriter, r *http.Request) {
	utils.LogInfo.Println("list albums")

	albums, err := ctrl.albumDao.GetAll()
	if err != nil {
		utils.LogError.Println(err)
		SendJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	SendJSONOk(w, albums)
}

// GetAlbumsByAuthor retrieve albums by author id
func (ctrl *AlbumController) GetAlbumsByAuthor(w http.ResponseWriter, r *http.Request) {
	authorId := ParamAsString("authorId", r)
	utils.LogInfo.Println("list albums of author : ", authorId)

	albums, err := ctrl.albumDao.GetByAuthor(authorId)
	if err != nil {
		utils.LogError.Println(err)
		SendJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	SendJSONOk(w, albums)
}

// GetAlbum retrieve an album by id
func (ctrl *AlbumController) GetAlbum(w http.ResponseWriter, r *http.Request) {
	albumId := ParamAsString("id", r)
	utils.LogInfo.Println("album : ", albumId)

	album, err := ctrl.albumDao.Get(albumId)
	if err != nil {
		utils.LogError.Println(err)
		SendJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.LogInfo.Println("album : ", album)
	SendJSONOk(w, album)
}

// CreateAlbum create an album
func (ctrl *AlbumController) CreateAlbum(w http.ResponseWriter, r *http.Request) {
	utils.LogInfo.Println("create album")
	album := &model.Album{}
	err := GetJSONContent(album, r)
	if err != nil {
		utils.LogError.Println(err)
		SendJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	authorExist, err := ctrl.authorDao.Exist(album.AuthorId)
	if err != nil {
		utils.LogError.Println(err)
		SendJSONError(w, err.Error(), http.StatusNotFound)
		return
	} else if authorExist == false {
		SendJSONError(w, "author not found", http.StatusNotFound)
		return
	}

	album, err = ctrl.albumDao.Upsert(album)
	if err != nil {
		utils.LogError.Println(err)
		SendJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.LogInfo.Println("album : ", album)
	SendJSONWithHTTPCode(w, album, http.StatusCreated)
}

// UpdateAlbum update an album by id
func (ctrl *AlbumController) UpdateAlbum(w http.ResponseWriter, r *http.Request) {
	album := &model.Album{}
	err := GetJSONContent(album, r)
	if err != nil {
		utils.LogError.Println(err)
		SendJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	utils.LogInfo.Println("update album : ", album.Id)

	authorExist, err := ctrl.authorDao.Exist(album.AuthorId)
	if err != nil {
		utils.LogError.Println(err)
		SendJSONError(w, err.Error(), http.StatusNotFound)
		return
	} else if authorExist == false {
		SendJSONError(w, "author not found", http.StatusNotFound)
		return
	}

	albumExist, err := ctrl.albumDao.Exist(album.Id)
	if err != nil {
		utils.LogError.Println(err)
		SendJSONError(w, err.Error(), http.StatusNotFound)
		return
	} else if albumExist == false {
		SendJSONError(w, "album not found", http.StatusNotFound)
		return
	}

	album, err = ctrl.albumDao.Upsert(album)
	if err != nil {
		utils.LogError.Println(err)
		SendJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.LogInfo.Println("album : ", album)
	SendJSONOk(w, album)
}

// DeleteAlbum delete an album by id
func (ctrl *AlbumController) DeleteAlbum(w http.ResponseWriter, r *http.Request) {
	albumId := ParamAsString("id", r)
	utils.LogInfo.Println("delete album : ", albumId)

	err := ctrl.albumDao.Delete(albumId)
	if err != nil {
		utils.LogError.Println(err)
		SendJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.LogInfo.Println("deleted album : ", albumId)
	SendJSONWithHTTPCode(w, nil, http.StatusNoContent)
}
