package controllerAlbum

import (
	"github.com/gorilla/mux"
)

// Routes for album management
func Routes(router *mux.Router) {
	router.HandleFunc("/albums", GetAlbums).Methods("GET")
	router.HandleFunc("/albums/{authorId}", GetAlbumsByAuthor).Methods("GET")
	router.HandleFunc("/album/{albumId}", GetAlbum).Methods("GET")
	router.HandleFunc("/album/{albumId}", UpdateAlbum).Methods("PUT")
	router.HandleFunc("/album/{albumId}", DeleteAlbum).Methods("DELETE")
	router.HandleFunc("/album", AddAlbum).Methods("POST")
}