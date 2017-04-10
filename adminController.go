package main

import (
	"encoding/json"
	"net/http"
)

// Resource to delete all authors and albums with songs
func DeleteAll(w http.ResponseWriter, r *http.Request) {
	LogInfo.Println("Suppression de tout les clés")

	result, err := DeleteAllAuthorDB()
	if err != nil {
		LogError.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err = DeleteAllAlbumDB()
	if err != nil {
		LogError.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		LogError.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
