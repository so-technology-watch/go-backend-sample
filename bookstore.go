package main

import (
	"go-redis-sample/dao"
	"go-redis-sample/utils"
	"go-redis-sample/web"
	"net/http"
)

func main() {
	authorDAO, albumDAO, err := dao.GetDAO(dao.DAORedis)
	if err != nil {
		utils.LogError.Println("unable to build the required DAO")
	}

	authorCtrl := web.NewAuthorController(authorDAO)
	albumCtrl := web.NewAlbumController(albumDAO, authorDAO)
	adminCtrl := web.NewAdminController(albumDAO, authorDAO)

	router := web.NewRouter(authorCtrl, albumCtrl, adminCtrl)

	utils.LogError.Println(http.ListenAndServe(":8080", router))
}
