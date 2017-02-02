package main

import (
	"fmt"
	"net/http"
	"go-redis-sample/config"
	"go-redis-sample/controllers"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", HomePage).Methods("GET")
	router.HandleFunc("/albums", controllers.GetAlbums).Methods("GET")
	router.HandleFunc("/albums/{authorId}", controllers.GetAlbumsByAuthor).Methods("GET")
	router.HandleFunc("/authors", controllers.GetAuthors).Methods("GET")
	router.HandleFunc("/album/{albumId}", controllers.GetAlbum).Methods("GET")
	router.HandleFunc("/author/{authorId}", controllers.GetAuthor).Methods("GET")
	router.HandleFunc("/album/{albumId}", controllers.UpdateAlbum).Methods("PUT")
	router.HandleFunc("/author/{authorId}", controllers.UpdateAuthor).Methods("PUT")
	router.HandleFunc("/album/{albumId}", controllers.DeleteAlbum).Methods("DELETE")
	router.HandleFunc("/author/{authorId}", controllers.DeleteAuthor).Methods("DELETE")
	router.HandleFunc("/album", controllers.AddAlbum).Methods("POST")
	router.HandleFunc("/author", controllers.AddAuthor).Methods("POST")
	router.HandleFunc("/admin/delete", controllers.DeleteAll).Methods("DELETE")

	config.Error.Println(http.ListenAndServe(":8080", router))
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
}
