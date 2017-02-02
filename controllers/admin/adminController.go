package controllerAdmin

import (
	"encoding/json"
	"net/http"
	"go-redis-sample/config"
	"go-redis-sample/models"
)

func DeleteAll(w http.ResponseWriter, r *http.Request) {
	config.LogInfo.Println("Suppression de tout les clés")

	result, err := models.DeleteAllAuthor()
	if err != nil {
		config.LogError.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err = models.DeleteAllAlbum()
	if err != nil {
		config.LogError.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		config.LogError.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
