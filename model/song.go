package model

import (
	"errors"
)

// Structure of a song
type Song struct {
	Number string `json:"number"`
	Title  string `json:"title"`
}

func NewSong(number, title string) *Song {
	return &Song{
		Number: number,
		Title:	title,
	}
}

// Validation of a song structure
func (song Song) Valid() error {
	if song.Number == "" {
		return errors.New("number of song is mandatory")
	}
	if song.Title == "" {
		return errors.New("title of song is mandatory")
	}
	return nil
}
