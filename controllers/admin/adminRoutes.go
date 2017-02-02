package controllerAdmin

import (
	"github.com/gorilla/mux"
)

func Routes(router *mux.Router) {
	router.HandleFunc("/admin/delete", DeleteAll).Methods("DELETE")
}