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
func (author Author) Valid() (error) {
	if author.Firstname == "" {
		return errors.New("Firstname is mandatory")
	}
	if author.Lastname == "" {
		return errors.New("Lastname is mandatory")
	}
	return nil
}

// Create an author in database
func CreateAuthorDB(author *Author) (int64, error) {
	mapAuthor := map[string]string{
		"firstname":	author.Firstname,
		"lastname": 	author.Lastname,
	}

	// Increment number of authors
	newId := DB.Incr(AuthorStr)
	if newId.Err() != nil {
		return -1, newId.Err()
	}

	// Insert author in database
	result := DB.HMSet(AuthorIdStr + strconv.FormatInt(newId.Val(), 10), mapAuthor)
	if result.Err() != nil {
		return -1, result.Err()
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
		result := DB.HGetAll(keys.Val()[i])
		if result.Err() != nil {
			return nil, result.Err()
		}

		author := &Author{keys.Val()[i], result.Val()["firstname"], result.Val()["lastname"]}
		authors = append(authors, author)
	}

	return authors, nil
}

// Collect an author from database
func GetAuthorDB(id string) (*Author, error) {
	result := DB.HGetAll(AuthorIdStr + id)
	if result.Err() != nil {
		return nil, result.Err()
	} else if len(result.Val()) == 0 {
		return nil, errors.New(AuthorIdStr + id + " don't exist !!")
	}

	author := &Author{Id: AuthorIdStr + id, Firstname: result.Val()["firstname"], Lastname: result.Val()["lastname"]}

	return author, nil
}

// Update an author in database
func UpdateAuthorDB(author *Author) (*Author, error) {
	// Verification if author exist
	resultAuthorExist := DB.Exists(author.Id)
	if resultAuthorExist.Err() != nil {
		return nil, resultAuthorExist.Err()
	} else if resultAuthorExist.Val() == false {
		return nil, errors.New(author.Id + " don't exist !!")
	}
	mapAuthor := map[string]string{
		"firstname":	author.Firstname,
		"lastname": 	author.Lastname,
	}

	// Update author in database
	result := DB.HMSet(author.Id, mapAuthor)
	if result.Err() != nil {
		return nil, result.Err()
	}

	return author, nil
}

// Delete an author in database
func DeleteAuthorDB(id string) (bool, error) {
	// Collect all albums identifiers
	keys := DB.Keys(AlbumIdStr + "*")
	if len(keys.Val()) == 0 {
		LogWarning.Println("No albums !!")
	}

	for i := 0; i < len(keys.Val()); i++ {
		// Collect album by identifier
		result := DB.HGetAll(keys.Val()[i])
		if result.Err() != nil {
			return false, result.Err()
		} else if len(result.Val()) == 0 {
			return false, errors.New(AuthorIdStr + id + " don't exist !!")
		}

		// If author of the album is the same as the author in parameter,
		// album is deleted in database
		if id == result.Val()["idAuthor"] {
			resultDelAlbum := DB.Del(keys.Val()[i])
			if resultDelAlbum.Err() != nil {
				return false, resultDelAlbum.Err()
			} else if resultDelAlbum.Val() == 0 {
				return false, errors.New(keys.Val()[i] + " don't exist !!")
			}
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
