package main

import (
	"strconv"
	"encoding/json"
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
func (album Album) valid() (error) {
	if album.Title == "" {
		return errors.New("Title is mandatory")
	}
	if album.Description == "" {
		return errors.New("Description is mandatory")
	}
	if album.IdAuthor == "" {
		return errors.New("Author ID is mandatory")
	}
	for i:=0; i<len(album.Songs); i++ {
		err := album.Songs[i].valid()
		if err != nil {
			return err
		}

	}
	return nil
}

// Validation of a song structure
func (song Song) valid() (error) {
	if song.Number == "" {
		return errors.New("Number of song is mandatory")
	}
	if song.Title == "" {
		return errors.New("Title of song is mandatory")
	}
	return nil
}

// Verification if album exist
func (album Album) exist() (error) {
	resultAuthorExist := DB.Exists(album.Id)
	if resultAuthorExist.Err() != nil {
		return resultAuthorExist.Err()
	} else if resultAuthorExist.Val() == false {
		return errors.New(album.Id + " don't exist !!")
	}
	return nil
}

// Save album
func (album Album) save() (error) {
	if err := album.valid(); err != nil {
		return err
	}

	songs, err := json.Marshal(album.Songs)
	if err != nil {
		return err
	}

	mapAlbum := map[string]string{
		"title":     	album.Title,
		"description":  album.Description,
		"idAuthor": 	album.IdAuthor,
		"songs": 	string(songs),
	}

	// Save album with songs in database
	result := DB.HMSet(album.Id, mapAlbum)
	return result.Err()
}

func constructAlbum(id, title, description, authorId string, tabSongs []byte) (album Album) {
	var songs []Song
	json.Unmarshal(tabSongs, &songs)
	return Album{Id: id, Title: title, Description: description, IdAuthor: authorId, Songs: songs}
}

// Create an album in database
func CreateAlbumDB(title, description, authorId string, songs []byte) (int64, error) {
	// Increment number of albums
	newId := DB.Incr(AlbumStr)
	if newId.Err() != nil {
		return -1, newId.Err()
	}

	album := constructAlbum(AlbumIdStr + strconv.FormatInt(newId.Val(), 10), title, description, authorId, songs)

	if err := album.save(); err != nil {
		return -1, err
	}

	return newId.Val(), nil
}

// Collect all albums from database
func GetAlbumsDB() ([]*Album, error) {
	var albums []*Album

	// Collect all albums identifiers
	keys := DB.Keys(AlbumIdStr + "*")
	if len(keys.Val()) == 0 {
		return nil, errors.New("No albums !!")
	}

	for i := 0; i < len(keys.Val()); i++ {
		// Collect album by identifier
		album, err := GetAlbumDB(keys.Val()[i])
		if err != nil {
			return nil, err
		}

		albums = append(albums, album)
	}

	return albums, nil
}

// Collect albums from database by author identifier
func GetAlbumsByAuthorDB(idAuthor string) ([]*Album, error) {
	var albums []*Album

	// Verification if author exist
	resultAuthorExist := DB.Exists(AuthorIdStr + idAuthor)
	if resultAuthorExist.Err() != nil {
		return nil, resultAuthorExist.Err()
	} else if resultAuthorExist.Val() == false {
		return nil, errors.New(AuthorIdStr + idAuthor + " don't exist !!")
	}

	// Collect all albums identifiers
	keys := DB.Keys(AlbumIdStr + "*")
	if len(keys.Val()) == 0 {
		return nil, errors.New("No albums !!")
	}

	for i := 0; i < len(keys.Val()); i++ {
		// Collect album by identifier
		result := DB.HGetAll(keys.Val()[i])
		if result.Err() != nil {
			return nil, result.Err()
		}

		if result.Val()["idAuthor"] == idAuthor {
			album := constructAlbum(keys.Val()[i], result.Val()["title"], result.Val()["description"], result.Val()["idAuthor"], []byte(result.Val()["songs"]))
			albums = append(albums, &album)
		}
	}

	if len(albums) == 0 {
		return nil, errors.New("No albums !!")
	}

	return albums, nil
}

// Collect an album from database
func GetAlbumDB(id string) (*Album, error) {
	result := DB.HGetAll(id)
	if result.Err() != nil {
		return nil, result.Err()
	} else if len(result.Val()) == 0 {
		return nil, errors.New(AlbumIdStr + id + " don't exist !!")
	}

	album := constructAlbum(AlbumIdStr + id, result.Val()["title"], result.Val()["description"], result.Val()["idAuthor"], []byte(result.Val()["songs"]))

	return &album, nil
}

// Update an album in database
func UpdateAlbumDB(id, title, description, authorId string, songs []byte) (*Album, error) {
	album := constructAlbum(AlbumIdStr + id, title, description, authorId, songs)

	// Verification if author exist
	resultAuthorExist := DB.Exists(AuthorIdStr + album.IdAuthor)
	if resultAuthorExist.Err() != nil {
		return nil, resultAuthorExist.Err()
	} else if resultAuthorExist.Val() == false {
		return nil, errors.New(AuthorIdStr + album.IdAuthor + " don't exist !!")
	}

	if err := album.exist(); err != nil {
		return nil, err
	}

	if err := album.save(); err != nil {
		return nil, err
	}

	return &album, nil
}

// Delete an album in database
func DeleteAlbumDB(id string) (bool, error) {
	result := DB.Del(AlbumIdStr + id)
	if result.Err() != nil {
		return false, result.Err()
	} else if result.Val() == 0 {
		return false, errors.New(AlbumIdStr + id + " don't exist !!")
	}

	return true, nil
}

// Delete all albums in database
func DeleteAllAlbumDB() (bool, error) {
	// Collect all identifiers of albums
	keys := DB.Keys(AlbumIdStr + "*")
	if len(keys.Val()) == 0 { // If no albums in database
		LogWarning.Println("No albums !!")
	}

	for i := 0; i < len(keys.Val()); i++ {
		// Deletion of album by identifier
		resultDelAlbums := DB.Del(keys.Val()[i])
		if resultDelAlbums.Err() != nil {
			return false, resultDelAlbums.Err()
		}
	}

	// Delete number of albums in database
	resultDelNbAlbum := DB.Del(AlbumStr)
	if resultDelNbAlbum.Err() != nil {
		return false, resultDelNbAlbum.Err()
	}

	return true, nil
}