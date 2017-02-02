package controllerAuthor

import (
	"github.com/gorilla/mux"
)

func Routes(router *mux.Router) {
	router.HandleFunc("/authors", GetAuthors).Methods("GET")
	router.HandleFunc("/author/{authorId}", GetAuthor).Methods("GET")
	router.HandleFunc("/author/{authorId}", UpdateAuthor).Methods("PUT")
	router.HandleFunc("/author/{authorId}", DeleteAuthor).Methods("DELETE")
	router.HandleFunc("/author", AddAuthor).Methods("POST")

}