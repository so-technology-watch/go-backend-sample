package main

import (
	"errors"
	"strconv"
)

const AuthorStr = "author"
const AuthorIdStr = "author:"

// Structure of an author
type Author struct {
	Id            	string `json:"id"`
	Firstname     	string `json:"firstname"`
	Lastname	string `json:"lastname"`
}

// Validation of an author structure
func (author Author) valid() (error) {
	if author.Firstname == "" {
		return errors.New("Firstname is mandatory")
	}
	if author.Lastname == "" {
		return errors.New("Lastname is mandatory")
	}
	return nil
}

// Verification if author exist
func (author Author) exist() (error) {
	resultAuthorExist := DB.Exists(author.Id)
	if resultAuthorExist.Err() != nil {
		return resultAuthorExist.Err()
	} else if resultAuthorExist.Val() == false {
		return errors.New(author.Id + " don't exist !!")
	}
	return nil
}

// Save author
func (author Author) save() (error) {
	if err := author.valid(); err != nil {
		return err
	}

	mapAuthor := map[string]string{
		"firstname":	author.Firstname,
		"lastname": 	author.Lastname,
	}

	// Save author in database
	result := DB.HMSet(author.Id, mapAuthor)
	return result.Err()
}

// Create an author in database
func CreateAuthorDB(firstname, lastname string) (int64, error) {
	// Increment number of authors
	newId := DB.Incr(AuthorStr)
	if newId.Err() != nil {
		return -1, newId.Err()
	}

	author := Author{Id: AuthorIdStr + strconv.FormatInt(newId.Val(), 10), Firstname: firstname, Lastname: lastname}

	if err := author.save(); err != nil {
		return -1, err
	}

	return newId.Val(), nil
}

// Collect all authors from database
func GetAuthorsDB() ([]*Author, error) {
	var authors []*Author

	// Collect all authors identifiers
	keys := DB.Keys(AuthorIdStr + "*")
	if len(keys.Val()) == 0 {
		return nil, errors.New("No authors !!")
	}

	for i := 0; i < len(keys.Val()); i++ {
		// Collect author by identifier
		author, err := GetAuthorDB(keys.Val()[i])
		if err != nil {
			return nil, err
		}

		authors = append(authors, author)
	}

	return authors, nil
}

// Collect an author from database
func GetAuthorDB(id string) (*Author, error) {
	result := DB.HGetAll(id)
	if result.Err() != nil {
		return nil, result.Err()
	} else if len(result.Val()) == 0 {
		return nil, errors.New(id + " don't exist !!")
	}

	author := Author{Id: id, Firstname: result.Val()["firstname"], Lastname: result.Val()["lastname"]}

	return &author, nil
}

// Update an author in database
func UpdateAuthorDB(id, firstname, lastname string) (*Author, error) {
	author := Author{Id: AuthorIdStr + id, Firstname: firstname, Lastname: lastname}

	if err := author.exist(); err != nil {
		return nil, err
	}

	if err := author.save(); err != nil {
		return nil, err
	}

	return &author, nil
}

// Delete an author in database
func DeleteAuthorDB(id string) (bool, error) {
	albums, err := GetAlbumsByAuthorDB(id)
	if err != nil {
		return false, err
	}
	for i := 0; i < len(albums); i++ {
		resultDelAlbum := DB.Del(albums[i].Id)
		if resultDelAlbum.Err() != nil {
			return false, resultDelAlbum.Err()
		} else if resultDelAlbum.Val() == 0 {
			return false, errors.New(albums[i].Id + " don't exist !!")
		}
	}

	// Delete author in database
	resultDelAuthor := DB.Del(AuthorIdStr + id)
	if resultDelAuthor.Err() != nil {
		return false, resultDelAuthor.Err()
	} else if resultDelAuthor.Val() == 0 {
		return false, errors.New(AuthorIdStr + id + " don't exist !!")
	}

	return true, nil
}

// Delete all authors in database
func DeleteAllAuthorDB() (bool, error) {
	// Collect all identifiers of authors
	keys := DB.Keys(AuthorIdStr + "*")
	if len(keys.Val()) == 0 {
		LogWarning.Println("No authors !!")
	}

	for i := 0; i < len(keys.Val()); i++ {
		// Deletion of author by identifier
		resultDelAuthors := DB.Del(keys.Val()[i])
		if resultDelAuthors.Err() != nil {
			return false, resultDelAuthors.Err()
		}
	}

	// Delete number of authors in database
	resultDelNbAuthor := DB.Del(AuthorStr)
	if resultDelNbAuthor.Err() != nil {
		return false, resultDelNbAuthor.Err()
	}

	return true, nil
}
