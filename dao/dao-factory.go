package dao

import (
	"database/sql"
	"errors"
	"os"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/BurntSushi/toml"
	// driver mysql
	_ "github.com/go-sql-driver/mysql"
	// driver sqlite
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"gopkg.in/mgo.v2"
	"gopkg.in/redis.v5"
)

// DBType is the database type
type DBType int

// DBConfig is the database configuration
type DBConfig struct {
	Url      string
	Port     string
	User     string
	Password string
	Database string
	File     string
}

const (
	RedisDAO DBType = iota
	MongoDAO
	DatastoreDAO
	MySQLDAO
	SQLiteDAO
	MockDAO

	AppName = "Todolist"

	timeout  = 5 * time.Second
	poolSize = 35
)

var (
	ErrorDAONotFound = errors.New("unknown DAO type")

	redisLocalConfig = DBConfig{
		Url:      os.Getenv("URL_DB"),
		Password: "",
		Database: "",
		Port:     "6379",
	}

	mongoLocalConfig = DBConfig{
		Url:      os.Getenv("URL_DB"),
		Password: "",
		Database: "todolist",
		Port:     "27017",
	}

	mySQLLocalConfig = DBConfig{
		Url:      os.Getenv("URL_DB"),
		User:     "root",
		Password: "password",
		Database: "todolist",
		Port:     "3306",
	}

	sqliteLocalConfig = DBConfig{
		File: ":memory:",
	}

	createDatabaseSQLStatements = []string{
		`CREATE DATABASE IF NOT EXISTS todolist DEFAULT CHARACTER SET = 'utf8' DEFAULT COLLATE 'utf8_general_ci';`,
		`USE todolist;`,
	}
	createTableSQLStatement = `CREATE TABLE IF NOT EXISTS task (
			id VARCHAR(50) NOT NULL,
			title VARCHAR(100) NULL,
			description VARCHAR(500) NULL,
			status INT NULL,
			creationDate DATETIME NULL,
			modificationDate DATETIME NULL,
			PRIMARY KEY (id)
		)`
)

// GetDAO get TaskDAO
func GetDAO(daoType DBType, dbConfigFile string) (TaskDAO, error) {
	switch daoType {
	case RedisDAO:
		config := getConfig(RedisDAO, dbConfigFile)
		redisCli := initRedis(config)
		return NewTaskDAORedis(redisCli), nil
	case MongoDAO:
		config := getConfig(MongoDAO, dbConfigFile)
		mongoSession := initMongo(config)
		return NewTaskDAOMongo(mongoSession), nil
	case DatastoreDAO:
		clientDatastore := initDatastore()
		return NewTaskDAODatastore(clientDatastore), nil
	case MySQLDAO:
		config := getConfig(MySQLDAO, dbConfigFile)
		mySQLSession := initMySQL(config)
		createDatabaseSQL(mySQLSession)
		createTableSQL(mySQLSession)
		return NewTaskDAOSQL(mySQLSession), nil
	case SQLiteDAO:
		config := getConfig(SQLiteDAO, dbConfigFile)
		sqliteSession := initSQLite(config)
		createTableSQL(sqliteSession)
		return NewTaskDAOSQL(sqliteSession), nil
	case MockDAO:
		return NewTaskDAOMock(), nil
	default:
		return nil, ErrorDAONotFound
	}
}

func initDatastore() *datastore.Client {
	logrus.Println("datastore connexion ")
	ctx := context.Background()

	// connection to the Datastore
	client, err := datastore.NewClient(ctx, AppName)
	if err != nil {
		logrus.Error("datastore connexion error :", err)
		panic(err)
	}

	// verification of connection
	t, err := client.NewTransaction(ctx)
	if err != nil {
		logrus.Error("datastore connexion error :", err)
		panic(err)
	}
	if err := t.Rollback(); err != nil {
		logrus.Error("datastore connexion error :", err)
		panic(err)
	}
	return client
}

func initRedis(dbConfig DBConfig) *redis.Client {
	logrus.Println("redis connexion " + dbConfig.Url)

	// connection to the Redis database
	redisCli := redis.NewClient(&redis.Options{
		Addr:     dbConfig.Url + ":" + dbConfig.Port,
		Password: dbConfig.Password,
		DB:       int(RedisDAO),
	})

	// verification of connection
	ok, err := redisCli.Ping().Result()
	if err != nil {
		logrus.Error("redis connexion error :", err.Error())
		panic(err)
	} else {
		logrus.Println("redis connexion OK :", ok)
	}

	return redisCli
}

func initMongo(dbConfig DBConfig) *mgo.Session {
	logrus.Info("mongodb connexion " + dbConfig.Url)

	// connection to the Mongo database
	mongoSession, err := mgo.DialWithTimeout("mongodb://"+dbConfig.Url+":"+dbConfig.Port+"/"+dbConfig.Database, timeout)
	if err != nil {
		logrus.Error("mongodb connexion error :", err.Error())
		panic(err)
	} else {
		logrus.Info("mongodb connexion OK")
	}

	mongoSession.SetSyncTimeout(timeout)
	mongoSession.SetSocketTimeout(timeout)
	mongoSession.SetMode(mgo.Monotonic, true)
	mongoSession.SetPoolLimit(poolSize)

	return mongoSession
}

func initMySQL(dbConfig DBConfig) *sql.DB {
	logrus.Println("mysql connexion " + dbConfig.Url)

	mySQLSession, err := sql.Open("mysql", dbConfig.User+":"+dbConfig.Password+"@tcp("+dbConfig.Url+":"+dbConfig.Port+")/"+dbConfig.Database+"?parseTime=true")
	if err != nil {
		logrus.Error("mysql connexion error :", err.Error())
		panic(err.Error())
	}
	err = mySQLSession.Ping()
	if err != nil {
		logrus.Error("mysql connexion error :", err.Error())
		panic(err.Error())
	}

	return mySQLSession
}

func initSQLite(dbConfig DBConfig) *sql.DB {
	logrus.Println("sqlite connexion " + dbConfig.File)

	sqliteSession, err := sql.Open("sqlite3", dbConfig.File)
	if err != nil {
		logrus.Error("sqlite connexion error :", err.Error())
		panic(err.Error())
	}
	err = sqliteSession.Ping()
	if err != nil {
		logrus.Error("sqlite connexion error :", err.Error())
		panic(err.Error())
	}

	return sqliteSession
}

func createDatabaseSQL(session *sql.DB) error {
	for _, stmt := range createDatabaseSQLStatements {
		_, err := session.Exec(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

func createTableSQL(session *sql.DB) error {
	_, err := session.Exec(createTableSQLStatement)
	if err != nil {
		return err
	}
	return nil
}

func getConfig(daoType DBType, dbConfigFile string) DBConfig {
	var config DBConfig
	if dbConfigFile == "" {
		switch daoType {
		case RedisDAO:
			config = redisLocalConfig
		case MongoDAO:
			config = mongoLocalConfig
		case MySQLDAO:
			config = mySQLLocalConfig
		case SQLiteDAO:
			config = sqliteLocalConfig
		}
	} else {
		if _, err := toml.DecodeFile(dbConfigFile, &config); err != nil {
			logrus.Error("connexion parameters error :", err)
			panic(err)
		}
	}
	return config
}
