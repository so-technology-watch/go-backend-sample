package controllerAlbum

import (
	"encoding/json"
	"net/http"
	"go-redis-sample/config"
	"go-redis-sample/models"

	"github.com/gorilla/mux"
)

func GetAlbums(w http.ResponseWriter, r *http.Request) {
	config.LogInfo.Println("List albums")

	albums, err := models.GetAlbums()
	if err != nil {
		config.LogError.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(albums); err != nil {
		config.LogError.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func GetAlbumsByAuthor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	authorId := vars["authorId"]

	config.LogInfo.Println("List albums of author : AuthorId=" + authorId)

	albums, err := models.GetAlbumsByAuthor(authorId)
	if err != nil {
		config.LogError.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(albums); err != nil {
		config.LogError.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func GetAlbum(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	albumId := vars["albumId"]

	album, err := models.GetAlbum(albumId)
	if err != nil {
		config.LogError.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	config.LogInfo.Println("Album :", album)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(album); err != nil {
		config.LogError.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func UpdateAlbum(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var songs []models.Song
	json.Unmarshal([]byte(r.FormValue("songs")), &songs)
	album := &models.Album{Id: "album:" + vars["albumId"], Title: r.FormValue("title"), Description: r.FormValue("description"), IdAuthor: r.FormValue("authorId"), Songs: songs}

	album, err := models.UpdateAlbum(album)
	if err != nil {
		config.LogError.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	config.LogInfo.Println("Update album :", album)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(album); err != nil {
		config.LogError.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func AddAlbum(w http.ResponseWriter, r *http.Request) {
	var songs []models.Song
	json.Unmarshal([]byte(r.FormValue("songs")), &songs)
	album := &models.Album{Title: r.FormValue("title"), Description: r.FormValue("description"), IdAuthor: r.FormValue("authorId"), Songs: songs}

	config.LogInfo.Println("Add album :", album)

	id, err := models.CreateAlbum(album)
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

func DeleteAlbum(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	albumId := vars["albumId"]

	config.LogError.Println("Delete album : Id=" + albumId)

	result, err := models.DeleteAlbum(albumId)
	if err != nil {
		config.LogError.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		config.LogError.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
