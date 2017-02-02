package controllers

import (
	"encoding/json"
	"net/http"
	"go-redis-sample/config"
	"go-redis-sample/models"

	"github.com/gorilla/mux"
)

func GetAlbums(w http.ResponseWriter, r *http.Request) {
	config.Info.Println("List albums")

	albums, err := models.GetAlbumsDB()
	if err != nil {
		config.Error.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(albums); err != nil {
		config.Error.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func GetAlbumsByAuthor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	authorId := vars["authorId"]

	config.Info.Println("List albums of author : " + authorId)

	albums, err := models.GetAlbumsByAuthorDB(authorId)
	if err != nil {
		config.Error.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(albums); err != nil {
		config.Error.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func GetAlbum(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	albumId := vars["albumId"]

	album, err := models.GetAlbumDB(albumId)
	if err != nil {
		config.Error.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	config.Info.Println("Album : ", album)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(album); err != nil {
		config.Error.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func UpdateAlbum(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	albumId := vars["albumId"]

	config.Info.Println("Update album : ", albumId)

	var songs []models.Song
	json.Unmarshal([]byte(r.FormValue("songs")), &songs)
	album := &models.Album{Id: "album:" + vars["albumId"], Title: r.FormValue("title"), Description: r.FormValue("description"), IdAuthor: r.FormValue("authorId"), Songs: songs}

	album, err := models.UpdateAlbumDB(album)
	if err != nil {
		config.Error.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(album); err != nil {
		config.Error.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func AddAlbum(w http.ResponseWriter, r *http.Request) {
	config.Info.Println("Add album : " + r.FormValue("title") + " " + r.FormValue("description") + " "+ r.FormValue("authorId"))

	var songs []models.Song
	json.Unmarshal([]byte(r.FormValue("songs")), &songs)
	album := &models.Album{Id: "album:" + r.FormValue("authorId"), Title: r.FormValue("title"), Description: r.FormValue("description"), IdAuthor: r.FormValue("authorId"), Songs: songs}

	id, err := models.CreateAlbumDB(album)
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

func DeleteAlbum(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	albumId := vars["albumId"]

	config.Info.Println("Delete album : ", albumId)

	result, err := models.DeleteAlbumDB(albumId)
	if err != nil {
		config.Error.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		config.Error.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
