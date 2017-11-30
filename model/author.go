package model

import (
	"errors"
)

// Structure of an author
type Author struct {
	Id        string `json:"id, omitempty"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

func NewAuthor(id, firstname, lastname string) *Author {
	return &Author{
		Id:         id,
		Firstname:	firstname,
		Lastname: 	lastname,
	}
}

// Validation of an author structure
func (author Author) Valid() error {
	if author.Firstname == "" {
		return errors.New("firstname is mandatory")
	}
	if author.Lastname == "" {
		return errors.New("lastname is mandatory")
	}
	return nil
}
