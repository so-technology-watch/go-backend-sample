package controllerAdmin

import (
	"github.com/gorilla/mux"
)

// Routes for admin
func Routes(router *mux.Router) {
	router.HandleFunc("/admin/delete", DeleteAll).Methods("DELETE")
}