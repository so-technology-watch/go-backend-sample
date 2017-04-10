package main

import (
	"github.com/gorilla/mux"
)

// Routes
func Routes(router *mux.Router) {
	router.HandleFunc("/albums", GetAlbums).Methods("GET")
	router.HandleFunc("/albums/{authorId}", GetAlbumsByAuthor).Methods("GET")
	router.HandleFunc("/album/{albumId}", GetAlbum).Methods("GET")
	router.HandleFunc("/album/{albumId}", UpdateAlbum).Methods("PUT")
	router.HandleFunc("/album/{albumId}", DeleteAlbum).Methods("DELETE")
	router.HandleFunc("/album", AddAlbum).Methods("POST")

	router.HandleFunc("/authors", GetAuthors).Methods("GET")
	router.HandleFunc("/author/{authorId}", GetAuthor).Methods("GET")
	router.HandleFunc("/author/{authorId}", UpdateAuthor).Methods("PUT")
	router.HandleFunc("/author/{authorId}", DeleteAuthor).Methods("DELETE")
	router.HandleFunc("/author", AddAuthor).Methods("POST")

	router.HandleFunc("/admin/delete", DeleteAll).Methods("DELETE")
}