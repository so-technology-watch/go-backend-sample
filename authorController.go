package main

import (
	"encoding/json"
	"net/http"
	

	"github.com/gorilla/mux"
)

func GetAuthors(w http.ResponseWriter, r *http.Request) {
	LogInfo.Println("List authors")

	authors, err := GetAuthorsDB()
	if err != nil {
		LogError.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(authors); err != nil {
		LogError.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func GetAuthor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	authorId := vars["authorId"]

	author, err := GetAuthorDB(authorId)
	if err != nil {
		LogError.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	LogInfo.Println("Author :", author)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(author); err != nil {
		LogError.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func UpdateAuthor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	author := &Author{Id: AuthorIdStr + vars["authorId"], Firstname: r.FormValue("firstname"), Lastname: r.FormValue("lastname")}
	err := author.Valid()
	if err != nil {
		LogError.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	author, err = UpdateAuthorDB(author)
	if err != nil {
		LogError.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	LogInfo.Println("Update author :", author)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(author); err != nil {
		LogError.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func AddAuthor(w http.ResponseWriter, r *http.Request) {
	author := &Author{Firstname: r.FormValue("firstname"), Lastname: r.FormValue("lastname")}
	err := author.Valid()
	if err != nil {
		LogError.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	LogInfo.Println("Add author :", author)

	id, err := CreateAuthorDB(author)
	if err != nil {
		LogError.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(id); err != nil {
		LogError.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func DeleteAuthor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	authorId := vars["authorId"]

	LogInfo.Println("Delete author : Id=" + authorId)

	result, err := DeleteAuthorDB(authorId)
	if err != nil {
		LogError.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		LogError.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}