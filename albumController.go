package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func GetAlbums(w http.ResponseWriter, r *http.Request) {
	LogInfo.Println("List albums")

	albums, err := GetAlbumsDB()
	if err != nil {
		LogError.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(albums); err != nil {
		LogError.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func GetAlbumsByAuthor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	authorId := vars["authorId"]

	LogInfo.Println("List albums of author : AuthorId=" + authorId)

	albums, err := GetAlbumsByAuthorDB(authorId)
	if err != nil {
		LogError.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(albums); err != nil {
		LogError.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func GetAlbum(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	albumId := vars["albumId"]

	album, err := GetAlbumDB(albumId)
	if err != nil {
		LogError.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	LogInfo.Println("Album :", album)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(album); err != nil {
		LogError.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func UpdateAlbum(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var songs []Song
	json.Unmarshal([]byte(r.FormValue("songs")), &songs)
	album := &Album{Id: AlbumIdStr + vars["albumId"], Title: r.FormValue("title"), Description: r.FormValue("description"), IdAuthor: r.FormValue("authorId"), Songs: songs}
	err := album.Valid()
	if err != nil {
		LogError.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	album, err = UpdateAlbumDB(album)
	if err != nil {
		LogError.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	LogInfo.Println("Update album :", album)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(album); err != nil {
		LogError.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func AddAlbum(w http.ResponseWriter, r *http.Request) {
	var songs []Song
	json.Unmarshal([]byte(r.FormValue("songs")), &songs)
	album := &Album{Title: r.FormValue("title"), Description: r.FormValue("description"), IdAuthor: r.FormValue("authorId"), Songs: songs}
	err := album.Valid()
	if err != nil {
		LogError.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	LogInfo.Println("Add album :", album)

	id, err := CreateAlbumDB(album)
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

func DeleteAlbum(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	albumId := vars["albumId"]

	LogError.Println("Delete album : Id=" + albumId)

	result, err := DeleteAlbumDB(albumId)
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
