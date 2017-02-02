package controllerAuthor

import (
	"encoding/json"
	"net/http"
	"go-redis-sample/config"
	"go-redis-sample/models"

	"github.com/gorilla/mux"
)

func GetAuthors(w http.ResponseWriter, r *http.Request) {
	config.LogInfo.Println("List authors")

	authors, err := models.GetAuthors()
	if err != nil {
		config.LogError.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(authors); err != nil {
		config.LogError.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}



func GetAuthor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	authorId := vars["authorId"]

	author, err := models.GetAuthor(authorId)
	if err != nil {
		config.LogError.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	config.LogInfo.Println("Author :", author)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(author); err != nil {
		config.LogError.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func UpdateAuthor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	author := &models.Author{Id: "author:" + vars["authorId"], Firstname: r.FormValue("firstname"), Lastname: r.FormValue("lastname")}

	author, err := models.UpdateAuthor(author)
	if err != nil {
		config.LogError.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	config.LogInfo.Println("Update author :", author)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(author); err != nil {
		config.LogError.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func AddAuthor(w http.ResponseWriter, r *http.Request) {
	author := &models.Author{Firstname: r.FormValue("firstname"), Lastname: r.FormValue("lastname")}

	config.LogInfo.Println("Add author :", author)

	id, err := models.CreateAuthor(author)
	if err != nil {
		config.LogError.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(id); err != nil {
		config.LogError.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func DeleteAuthor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	authorId := vars["authorId"]

	config.LogInfo.Println("Delete author : Id=" + authorId)

	result, err := models.DeleteAuthor(authorId)
	if err != nil {
		config.LogError.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		config.LogError.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
