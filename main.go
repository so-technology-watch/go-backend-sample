package main

import (
	"fmt"
	"net/http"
	"go-redis-sample/config"
	"go-redis-sample/controllers/admin"
	"go-redis-sample/controllers/album"
	"go-redis-sample/controllers/author"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)

	controllerAdmin.Routes(router)
	controllerAlbum.Routes(router)
	controllerAuthor.Routes(router)

	config.LogError.Println(http.ListenAndServe(":8080", router))
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
}
