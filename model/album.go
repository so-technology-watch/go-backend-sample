package model

import (
	"encoding/json"
	"errors"
)

// Structure of an album
type Album struct {
	Id          string `json:"id, omitempty"`
	Title       string `json:"title"`
	Description string `json:"description"`
	AuthorId    string `json:"authorId"`
	Songs       []Song `json:"songs"`
}

func NewAlbum(id, title, description, authorId string, tabSongs []byte) *Album {
	var songs []Song
	json.Unmarshal(tabSongs, &songs)
	return &Album{
		Id:          id,
		Title:       title,
		Description: description,
		AuthorId:    authorId,
		Songs:       songs,
	}
}

// Validation of an album structure
func (album Album) Valid() error {
	if album.Title == "" {
		return errors.New("Title is mandatory")
	}
	if album.Description == "" {
		return errors.New("Description is mandatory")
	}
	if album.AuthorId == "" {
		return errors.New("Author ID is mandatory")
	}
	for i := 0; i < len(album.Songs); i++ {
		err := album.Songs[i].Valid()
		if err != nil {
			return err
		}

	}
	return nil
}
