package web

import (
	"go-backend-sample/dao"
	"go-backend-sample/model"
	"net/http"
	"github.com/sirupsen/logrus"
)

const (
	prefixAuthor = "/authors"
)

// AuthorController is a controller for authors resource
type AuthorController struct {
	authorDao dao.AuthorDAO
	Routes    []Route
	Prefix    string
}

// NewAuthorController creates a new author controller to manage authors
func NewAuthorController(authorDAO dao.AuthorDAO) *AuthorController {
	controller := AuthorController{
		authorDao: authorDAO,
		Prefix:    prefixAuthor,
	}

	var routes []Route
	// GetAll
	routes = append(routes, Route{
		Name:        "Get all authors",
		Method:      http.MethodGet,
		Pattern:     "",
		HandlerFunc: controller.GetAuthors,
	})
	// Get
	routes = append(routes, Route{
		Name:        "Get one author",
		Method:      http.MethodGet,
		Pattern:     "/{id}",
		HandlerFunc: controller.GetAuthor,
	})
	// Create
	routes = append(routes, Route{
		Name:        "Create an author",
		Method:      http.MethodPost,
		Pattern:     "",
		HandlerFunc: controller.CreateAuthor,
	})
	// Update
	routes = append(routes, Route{
		Name:        "Update an author",
		Method:      http.MethodPut,
		Pattern:     "/{id}",
		HandlerFunc: controller.UpdateAuthor,
	})
	// Delete
	routes = append(routes, Route{
		Name:        "Delete an author",
		Method:      http.MethodDelete,
		Pattern:     "/{id}",
		HandlerFunc: controller.DeleteAuthor,
	})

	controller.Routes = routes

	return &controller
}

// GetAll retrieve all authors
func (ctrl *AuthorController) GetAuthors(w http.ResponseWriter, r *http.Request) {
	logrus.Println("list authors")

	authors, err := ctrl.authorDao.GetAll()
	if err != nil {
		logrus.Error(err)
		SendJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	SendJSONOk(w, authors)
}

// Get retrieve an author by id
func (ctrl *AuthorController) GetAuthor(w http.ResponseWriter, r *http.Request) {
	authorId := ParamAsString("id", r)
	logrus.Println("author : ", authorId)

	author, err := ctrl.authorDao.Get(authorId)
	if err != nil {
		logrus.Error(err)
		SendJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logrus.Println("author : ", author)
	SendJSONOk(w, author)
}

// Create create an author
func (ctrl *AuthorController) CreateAuthor(w http.ResponseWriter, r *http.Request) {
	logrus.Println("create author")
	author := &model.Author{}
	logrus.Println(r.Body)
	err := GetJSONContent(author, r)
	if err != nil {
		logrus.Error(err)
		SendJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	author, err = ctrl.authorDao.Upsert(author)
	if err != nil {
		logrus.Error(err)
		SendJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logrus.Println("author : ", author)
	SendJSONWithHTTPCode(w, author, http.StatusCreated)
}

// Update update an author by id
func (ctrl *AuthorController) UpdateAuthor(w http.ResponseWriter, r *http.Request) {
	author := &model.Author{}
	err := GetJSONContent(author, r)
	if err != nil {
		logrus.Error(err)
		SendJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	logrus.Println("update author : ", author.Id)

	authorExist, err := ctrl.authorDao.Exist(author.Id)
	if err != nil {
		logrus.Error(err)
		SendJSONError(w, err.Error(), http.StatusNotFound)
		return
	} else if authorExist == false {
		SendJSONError(w, "author not found", http.StatusNotFound)
		return
	}

	author, err = ctrl.authorDao.Upsert(author)
	if err != nil {
		logrus.Error(err)
		SendJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logrus.Println("author : ", author)
	SendJSONOk(w, author)
}

// Delete delete an entity by id
func (ctrl *AuthorController) DeleteAuthor(w http.ResponseWriter, r *http.Request) {
	authorId := ParamAsString("id", r)
	logrus.Println("delete author : ", authorId)

	err := ctrl.authorDao.Delete(authorId)
	if err != nil {
		logrus.Error(err)
		SendJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logrus.Println("deleted author : ", authorId)
	SendJSONWithHTTPCode(w, nil, http.StatusNoContent)
}
