package model

import (
	"errors"
)

// Structure of a song
type Song struct {
	Number string `json:"number"`
	Title  string `json:"title"`
}

// Validation of a song structure
func (song Song) Valid() error {
	if song.Number == "" {
		return errors.New("Number of song is mandatory")
	}
	if song.Title == "" {
		return errors.New("Title of song is mandatory")
	}
	return nil
}
