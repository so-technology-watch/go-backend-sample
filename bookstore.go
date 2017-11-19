package main

import (
	"go-backend-sample/dao"
	"go-backend-sample/utils"
	"go-backend-sample/web"
	"net/http"
	"time"
	"gopkg.in/urfave/cli.v1"
	"strconv"
	"os"
)

var (
	// Version is the version of the software
	Version string
	// BuildStmp is the build date
	BuildStmp string
	// GitHash is the git build hash
	GitHash string

	port               = 8020
	statisticsDuration = 20 * time.Second
	db                 = 0
	dbConfigFile       = ""
)

func main() {
	// new app
	app := cli.NewApp()
	app.Name = utils.AppName
	app.Usage = "bookstore service launcher"

	timeStmp, err := strconv.Atoi(BuildStmp)
	if err != nil {
		timeStmp = 0
	}
	app.Version = Version + ", build on " + time.Unix(int64(timeStmp), 0).String() + ", git hash " + GitHash
	app.Authors = []cli.Author{{Name: "xma"}}
	app.Copyright = "Copyright " + strconv.Itoa(time.Now().Year())

	// command line flags
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Value:       port,
			Name:        "port, p",
			Usage:       "Set the listening port of the webserver",
			Destination: &port,
		},
		cli.DurationFlag{
			Value:       statisticsDuration,
			Name:        "statd, s",
			Usage:       "Set the statistics accumulation duration (ex : 1h, 2h30m, 30s, 300ms)",
			Destination: &statisticsDuration,
		},
		cli.IntFlag{
			Value:       db,
			Name:        "db",
			Usage:       "Set the database connection parameters (0 - Redis | 1 - MongoDB | 2 - Mock)",
			Destination: &db,
		},
		cli.StringFlag{
			Value:       dbConfigFile,
			Name:        "dbConfig",
			Usage:       "Set the path of database connection parameters file",
			Destination: &dbConfigFile,
		},
	}

	// main action
	// sub action are also possible
	app.Action = func(c *cli.Context) error {
		authorDAO, albumDAO, err := dao.GetDAO(dao.DBType(db), dbConfigFile)
		if err != nil {
			utils.LogError.Println("unable to build the required DAO")
		}

		authorCtrl := web.NewAuthorController(authorDAO)
		albumCtrl := web.NewAlbumController(albumDAO, authorDAO)
		adminCtrl := web.NewAdminController(albumDAO, authorDAO)

		router := web.NewRouter(authorCtrl, albumCtrl, adminCtrl)

		utils.LogError.Println(http.ListenAndServe(":" + strconv.Itoa(port), router))

		return nil
	}

	// run the app
	err = app.Run(os.Args)
	if err != nil {
		utils.LogError.Println("run error %q\n", err)
	}
}

