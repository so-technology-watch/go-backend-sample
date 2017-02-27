package models

import (
	"strconv"
	"encoding/json"
	"go-redis-sample/config"
	"errors"
)

const AlbumStr = "album"
const AlbumIdStr = "album:"

// Structure of an album
type Album struct {
	Id            	string 	`json:"id"`
	Title          	string 	`json:"title"`
	Description	string 	`json:"description"`
	IdAuthor	string 	`json:"idAuthor"`
	Songs		[]Song	`json:"songs"`
}

// Structure of a song
type Song struct {
	Number		string	`json:"number"`
	Title          	string 	`json:"title"`
}

// Validation of an album structure
func (album Album) Valid() (error) {
	if album.Title == "" {
		return errors.New("Title is mandatory")
	}
	if album.Description == "" {
		return errors.New("Title is mandatory")
	}
	if album.IdAuthor == "" {
		return errors.New("Title is mandatory")
	}
	for i:=0; i<len(album.Songs); i++ {
		err := album.Songs[i].Valid()
		if err != nil {
			return err
		}

	}
	return nil
}

// Validation of a song structure
func (song Song) Valid() (error) {
	if song.Number == "" {
		return errors.New("Number of song is mandatory")
	}
	if song.Title == "" {
		return errors.New("Title of song is mandatory")
	}
	return nil
}

// Create an album in database
func CreateAlbum(album *Album) (int64, error) {
	// Verification if author exist
	resultAuthorExist := config.DB.Exists(AuthorIdStr + album.IdAuthor)
	if resultAuthorExist.Err() != nil {
		return -1, resultAuthorExist.Err()
	} else if resultAuthorExist.Val() == false {
		return -1, errors.New(album.IdAuthor + " don't exist !!")
	}

	songs, err := json.Marshal(album.Songs)
	if err != nil {
		return -1, err
	}

	mapAlbum := map[string]string{
		"title":    	album.Title,
		"description":	album.Description,
		"idAuthor": 	album.IdAuthor,
		"songs": 	string(songs),
	}

	// Increment number of albums
	newId := config.DB.Incr(AlbumStr)
	if newId.Err() != nil {
		return -1, newId.Err()
	}

	// Insert album in database
	result := config.DB.HMSet(AlbumIdStr + strconv.FormatInt(newId.Val(), 10), mapAlbum)
	if result.Err() != nil {
		return -1, result.Err()
	}

	return newId.Val(), nil
}

// Collect all albums from database
func GetAlbums() ([]*Album, error) {
	var albums []*Album

	// Collect all albums identifiers
	keys := config.DB.Keys(AlbumIdStr + "*")
	if len(keys.Val()) == 0 {
		return nil, errors.New("No albums !!")
	}

	for i := 0; i < len(keys.Val()); i++ {
		// Collect album by identifier
		result := config.DB.HGetAll(keys.Val()[i])
		if result.Err() != nil {
			return nil, result.Err()
		}

		var songs []Song
		json.Unmarshal([]byte(result.Val()["songs"]), &songs)
		album := &Album{Id: keys.Val()[i], Title: result.Val()["title"], Description: result.Val()["description"], IdAuthor: result.Val()["idAuthor"], Songs: songs}

		albums = append(albums, album)
	}

	return albums, nil
}

// Collect albums from database by author identifier
func GetAlbumsByAuthor(idAuthor string) ([]*Album, error) {
	var albums []*Album

	// Verification if author exist
	resultAuthorExist := config.DB.Exists(AuthorIdStr + idAuthor)
	if resultAuthorExist.Err() != nil {
		return nil, resultAuthorExist.Err()
	} else if resultAuthorExist.Val() == false {
		return nil, errors.New(AuthorIdStr + idAuthor + " don't exist !!")
	}

	// Collect all albums identifiers
	keys := config.DB.Keys(AlbumIdStr + "*")
	if len(keys.Val()) == 0 {
		return nil, errors.New("No albums !!")
	}

	for i := 0; i < len(keys.Val()); i++ {
		// Collect album by identifier
		result := config.DB.HGetAll(keys.Val()[i])
		if result.Err() != nil {
			return nil, result.Err()
		}

		if result.Val()["idAuthor"] == idAuthor {
			var songs []Song
			json.Unmarshal([]byte(result.Val()["songs"]), &songs)
			album := &Album{Id: keys.Val()[i], Title: result.Val()["title"], Description: result.Val()["description"], IdAuthor: result.Val()["idAuthor"], Songs: songs}
			albums = append(albums, album)
		}
	}

	if len(albums) == 0 {
		return nil, errors.New("No albums !!")
	}

	return albums, nil
}

// Collect an album from database
func GetAlbum(id string) (*Album, error) {
	result := config.DB.HGetAll(AlbumIdStr + id)
	if result.Err() != nil {
		return nil, result.Err()
	} else if len(result.Val()) == 0 {
		return nil, errors.New(AlbumIdStr + id + " don't exist !!")
	}

	var songs []Song
	json.Unmarshal([]byte(result.Val()["songs"]), &songs)
	album := &Album{Id: AlbumIdStr + id, Title: result.Val()["title"], Description: result.Val()["description"], IdAuthor: result.Val()["idAuthor"], Songs: songs}

	return album, nil
}

// Update an album in database
func UpdateAlbum(album *Album) (*Album, error) {
	// Verification if album exist
	resultAlbumExist := config.DB.Exists(album.Id)
	if resultAlbumExist.Err() != nil {
		return album, resultAlbumExist.Err()
	} else if resultAlbumExist.Val() == false {
		return album, errors.New(album.Id + " don't exist !!")
	}

	// Verification if author exist
	resultAuthorExist := config.DB.Exists("author:" + album.IdAuthor)
	if resultAuthorExist.Err() != nil {
		return album, resultAuthorExist.Err()
	} else if resultAuthorExist.Val() == false {
		return album, errors.New(AuthorIdStr + album.IdAuthor + " don't exist !!")
	}

	songs, err := json.Marshal(album.Songs)
	if err != nil {
		return album, err
	}

	mapAlbum := map[string]string{
		"title":     	album.Title,
		"description":  album.Description,
		"idAuthor": 	album.IdAuthor,
		"songs": 	string(songs),
	}

	// Update album with songs in database
	result := config.DB.HMSet(album.Id, mapAlbum)
	if result.Err() != nil {
		return album, result.Err()
	}

	return album, nil
}

// Delete an album in database
func DeleteAlbum(id string) (bool, error) {
	result := config.DB.Del(AlbumIdStr + id)
	if result.Err() != nil {
		return false, result.Err()
	} else if result.Val() == 0 {
		return false, errors.New(AlbumIdStr + id + " don't exist !!")
	}

	return true, nil
}

// Delete all albums in database
func DeleteAllAlbum() (bool, error) {
	// Collect all identifiers of albums
	keys := config.DB.Keys(AlbumIdStr + "*")
	if len(keys.Val()) == 0 { // If no albums in database
		config.LogWarning.Println("No albums !!")
	}

	for i := 0; i < len(keys.Val()); i++ {
		// Deletion of album by identifier
		resultDelAlbums := config.DB.Del(keys.Val()[i])
		if resultDelAlbums.Err() != nil {
			return false, resultDelAlbums.Err()
		}
	}

	// Delete number of albums in database
	resultDelNbAlbum := config.DB.Del(AlbumStr)
	if resultDelNbAlbum.Err() != nil {
		return false, resultDelNbAlbum.Err()
	}

	return true, nil
}