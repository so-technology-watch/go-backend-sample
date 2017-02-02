package controllers

import (
	"encoding/json"
	"net/http"
	"go-redis-sample/config"
	"go-redis-sample/models"

	"github.com/gorilla/mux"
)

func GetAuthors(w http.ResponseWriter, r *http.Request) {
	config.Info.Println("List authors")

	authors, err := models.GetAuthorsDB()
	if err != nil {
		config.Error.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(authors); err != nil {
		config.Error.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}



func GetAuthor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	authorId := vars["authorId"]

	author, err := models.GetAuthorDB(authorId)
	if err != nil {
		config.Error.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	config.Info.Println("Author : ", author)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(author); err != nil {
		config.Error.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func UpdateAuthor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	config.Info.Println("Update author : ", vars["authorId"])

	author := &models.Author{Id: "author:" + vars["authorId"], Firstname: r.FormValue("firstname"), Lastname: r.FormValue("lastname")}

	author, err := models.UpdateAuthorDB(author)
	if err != nil {
		config.Error.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(author); err != nil {
		config.Error.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func AddAuthor(w http.ResponseWriter, r *http.Request) {
	config.Info.Println("Add author : " + r.FormValue("firstname") + " " + r.FormValue("lastname"))

	author := &models.Author{Firstname: r.FormValue("firstname"), Lastname: r.FormValue("lastname")}

	id, err := models.CreateAuthorDB(author)
	if err != nil {
		config.Error.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(id); err != nil {
		config.Error.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func DeleteAuthor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	authorId := vars["authorId"]

	config.Info.Println("Delete author : ", authorId)

	result, err := models.DeleteAuthorDB(authorId)
	if err != nil {
		config.Error.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		config.Error.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
