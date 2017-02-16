package models

import (
	"errors"
	"strconv"
	"go-redis-sample/config"
)

const AuthorStr = "author"
const AuthorIdStr = "author:"

type Author struct {
	Id            	string `json:"id"`
	Firstname     	string `json:"firstname"`
	Lastname	string `json:"lastname"`
}

func (author Author) Valid() (error) {
	if author.Firstname == "" {
		return errors.New("Firstname is mandatory")
	}
	if author.Lastname == "" {
		return errors.New("Lastname is mandatory")
	}
	return nil
}

func CreateAuthor(author *Author) (int64, error) {
	mapAuthor := map[string]string{
		"firstname":	author.Firstname,
		"lastname": 	author.Lastname,
	}

	newId := config.DB.Incr(AuthorStr)
	if newId.Err() != nil {
		return -1, newId.Err()
	}

	result := config.DB.HMSet(AuthorIdStr + strconv.FormatInt(newId.Val(), 10), mapAuthor)
	if result.Err() != nil {
		return -1, result.Err()
	}

	return newId.Val(), nil
}

func GetAuthors() ([]*Author, error) {
	var authors []*Author

	keys := config.DB.Keys(AuthorIdStr + "*")
	if len(keys.Val()) == 0 {
		return nil, errors.New("No authors !!")
	}

	for i := 0; i < len(keys.Val()); i++ {
		result := config.DB.HGetAll(keys.Val()[i])
		if result.Err() != nil {
			return nil, result.Err()
		}

		author := &Author{keys.Val()[i], result.Val()["firstname"], result.Val()["lastname"]}
		authors = append(authors, author)
	}

	return authors, nil
}

func GetAuthor(id string) (*Author, error) {
	result := config.DB.HGetAll(AuthorIdStr + id)
	if result.Err() != nil {
		return nil, result.Err()
	} else if len(result.Val()) == 0 {
		return nil, errors.New(AuthorIdStr + id + " don't exist !!")
	}

	author := &Author{Id: AuthorIdStr + id, Firstname: result.Val()["firstname"], Lastname: result.Val()["lastname"]}

	return author, nil
}

func UpdateAuthor(author *Author) (*Author, error) {
	resultAuthorExist := config.DB.Exists(author.Id)
	if resultAuthorExist.Err() != nil {
		return nil, resultAuthorExist.Err()
	} else if resultAuthorExist.Val() == false {
		return nil, errors.New(author.Id + " don't exist !!")
	}
	mapAuthor := map[string]string{
		"firstname":	author.Firstname,
		"lastname": 	author.Lastname,
	}

	result := config.DB.HMSet(author.Id, mapAuthor)
	if result.Err() != nil {
		return nil, result.Err()
	}

	return author, nil
}

func DeleteAuthor(id string) (bool, error) {
	keys := config.DB.Keys(AlbumIdStr + "*")
	if len(keys.Val()) == 0 {
		config.LogWarning.Println("No albums !!")
	}

	for i := 0; i < len(keys.Val()); i++ {
		result := config.DB.HGetAll(keys.Val()[i])
		if result.Err() != nil {
			return false, result.Err()
		} else if len(result.Val()) == 0 {
			return false, errors.New(AuthorIdStr + id + " don't exist !!")
		}

		if id == result.Val()["idAuthor"] {
			resultDelAlbum := config.DB.Del(keys.Val()[i])
			if resultDelAlbum.Err() != nil {
				return false, resultDelAlbum.Err()
			} else if resultDelAlbum.Val() == 0 {
				return false, errors.New(keys.Val()[i] + " don't exist !!")
			}
		}
	}

	resultDelAuthor := config.DB.Del(AuthorIdStr + id)
	if resultDelAuthor.Err() != nil {
		return false, resultDelAuthor.Err()
	} else if resultDelAuthor.Val() == 0 {
		return false, errors.New(AuthorIdStr + id + " don't exist !!")
	}

	return true, nil
}

func DeleteAllAuthor() (bool, error) {
	keys := config.DB.Keys(AuthorIdStr + "*")
	if len(keys.Val()) == 0 {
		config.LogWarning.Println("No authors !!")
	}

	for i := 0; i < len(keys.Val()); i++ {
		resultDelAuthors := config.DB.Del(keys.Val()[i])
		if resultDelAuthors.Err() != nil {
			return false, resultDelAuthors.Err()
		}
	}

	resultDelNbAuthor := config.DB.Del(AuthorStr)
	if resultDelNbAuthor.Err() != nil {
		return false, resultDelNbAuthor.Err()
	}

	return true, nil
}
