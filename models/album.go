package models

import (
	"strconv"
	"encoding/json"
	"go-redis-sample/config"
	"errors"
)

type Album struct {
	Id            	string 	`json:"id"`
	Title          	string 	`json:"title"`
	Description	string 	`json:"description"`
	IdAuthor	string 	`json:"idAuthor"`
	Songs		[]Song	`json:"songs"`
}

type Song struct {
	Number		string	`json:"number"`
	Title          	string 	`json:"title"`
}

func CreateAlbumDB(album *Album) (int64, error) {
	resultAuthorExist := config.DB.Exists("author:" + album.IdAuthor)
	if resultAuthorExist.Err() != nil {
		config.Error.Println(resultAuthorExist.Err())
		return -1, resultAuthorExist.Err()
	} else if resultAuthorExist.Val() == false {
		config.Error.Println("Author : " + album.IdAuthor + " don't exist !!")
		return -1, resultAuthorExist.Err()
	}

	songs, err := json.Marshal(album.Songs)
	if err != nil {
		config.Error.Println(err)
		return -1, err
	}

	mapAlbum := map[string]string{
		"title":    	album.Title,
		"description":	album.Description,
		"idAuthor": 	album.IdAuthor,
		"songs": 	string(songs),
	}

	incr := config.DB.Incr("album")
	if incr.Err() != nil {
		config.Error.Println(incr.Err())
		return -1, incr.Err()
	}

	id, err := config.DB.Get("album").Int64()
	if err != nil {
		config.Error.Println(err)
		return -1, err
	}

	result := config.DB.HMSet("album:" + strconv.FormatInt(id, 10), mapAlbum)
	if result.Err() != nil {
		config.Error.Println(result.Err())
		return -1, result.Err()
	}

	return id, nil
}

func GetAlbumsDB() ([]*Album, error) {
	var albums []*Album

	keys := config.DB.Keys("album:*")
	if len(keys.Val()) == 0 {
		config.Info.Println("No albums !!")
		return nil, errors.New("No albums !!")
	}

	for i := 0; i < len(keys.Val()); i++ {
		result := config.DB.HGetAll(keys.Val()[i])
		if result.Err() != nil {
			config.Error.Println(result.Err())
			return nil, result.Err()
		}

		var songs []Song
		json.Unmarshal([]byte(result.Val()["songs"]), &songs)
		album := &Album{Id: keys.Val()[i], Title: result.Val()["title"], Description: result.Val()["description"], IdAuthor: result.Val()["idAuthor"], Songs: songs}

		albums = append(albums, album)
	}

	return albums, nil
}

func GetAlbumsByAuthorDB(idAuthor string) ([]*Album, error) {
	var albums []*Album
	resultAuthorExist := config.DB.Exists("author:" + idAuthor)
	if resultAuthorExist.Err() != nil {
		config.Error.Println(resultAuthorExist.Err())
		return nil, resultAuthorExist.Err()
	} else if resultAuthorExist.Val() == false {
		config.Error.Println("Author : " + idAuthor + " don't exist !!")
		return nil, errors.New("Author : " + idAuthor + " don't exist !!")
	}

	keys := config.DB.Keys("album:*")
	if len(keys.Val()) == 0 {
		config.Info.Println("No albums !!")
		return nil, errors.New("No albums !!")
	}

	for i := 0; i < len(keys.Val()); i++ {
		result := config.DB.HGetAll(keys.Val()[i])
		if result.Err() != nil {
			config.Error.Println(result.Err())
			return nil, result.Err()
		}

		if result.Val()["idAuthor"] == idAuthor {
			var songs []Song
			json.Unmarshal([]byte(result.Val()["songs"]), &songs)
			album := &Album{Id: keys.Val()[i], Title: result.Val()["title"], Description: result.Val()["description"], IdAuthor: result.Val()["idAuthor"], Songs: songs}
			albums = append(albums, album)
		}
	}

	return albums, nil
}

func GetAlbumDB(id string) (*Album, error) {
	result := config.DB.HGetAll("album:" + id)
	if result.Err() != nil {
		config.Error.Println(result.Err())
		return nil, result.Err()
	} else if len(result.Val()) == 0 {
		config.Error.Println("Album : " + id + " don't exist !!")
		return nil, errors.New("Album : " + id + " don't exist !!")
	}

	var songs []Song
	json.Unmarshal([]byte(result.Val()["songs"]), &songs)
	album := &Album{Id: "album:" + id, Title: result.Val()["title"], Description: result.Val()["description"], IdAuthor: result.Val()["idAuthor"], Songs: songs}

	return album, nil
}

func UpdateAlbumDB(album *Album) (*Album, error) {
	resultAlbumExist := config.DB.Exists("album:" + album.Id)
	if resultAlbumExist.Err() != nil {
		config.Error.Println(resultAlbumExist.Err())
		return album, resultAlbumExist.Err()
	} else if resultAlbumExist.Val() == false {
		config.Error.Println("Album : " + album.Id + " don't exist !!")
		return album, errors.New("Album : " + album.Id + " don't exist !!")
	}

	resultAuthorExist := config.DB.Exists("author:" + album.IdAuthor)
	if resultAuthorExist.Err() != nil {
		config.Error.Println(resultAuthorExist.Err())
		return album, resultAuthorExist.Err()
	} else if resultAuthorExist.Val() == false {
		config.Error.Println("Author : " + album.IdAuthor + " don't exist !!")
		return album, errors.New("Author : " + album.IdAuthor + " don't exist !!")
	}

	songs, err := json.Marshal(album.Songs)
	if err != nil {
		config.Error.Println(err)
		return album, err
	}

	mapAlbum := map[string]string{
		"title":     	album.Title,
		"description":  album.Description,
		"idAuthor": 	album.IdAuthor,
		"songs": 	string(songs),
	}

	result := config.DB.HMSet("album:" + album.Id, mapAlbum)
	if result.Err() != nil {
		config.Error.Println(result.Err())
		return album, result.Err()
	}

	return album, nil
}

func DeleteAlbumDB(id string) (bool, error) {
	result := config.DB.Del("album:" + id)
	if result.Err() != nil {
		config.Error.Println(result.Err())
		return false, result.Err()
	} else if result.Val() == 0 {
		config.Error.Println("Album : " + id + " don't exist !!")
		return false, errors.New("Album : " + id + " don't exist !!")
	}

	return true, nil
}

func DeleteAllAlbumDB() (bool, error) {
	keys := config.DB.Keys("album:*")
	if len(keys.Val()) == 0 {
		config.Info.Println("No albums !!")
	}

	for i := 0; i < len(keys.Val()); i++ {
		resultDelAlbums := config.DB.Del(keys.Val()[i])
		if resultDelAlbums.Err() != nil {
			config.Error.Println(resultDelAlbums.Err())
			return false, resultDelAlbums.Err()
		}
	}

	resultDelNbAlbum := config.DB.Del("album")
	if resultDelNbAlbum.Err() != nil {
		config.Error.Println(resultDelNbAlbum.Err())
		return false, resultDelNbAlbum.Err()
	}

	return true, nil
}