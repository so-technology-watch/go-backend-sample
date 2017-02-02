package models

import (
	"errors"
	"strconv"
	"go-redis-sample/config"
)

type Author struct {
	Id            	string `json:"id"`
	Firstname     	string `json:"firstname"`
	Lastname	string `json:"lastname"`
}

func CreateAuthorDB(author *Author) (int64, error) {
	mapAuthor := map[string]string{
		"firstname":	author.Firstname,
		"lastname": 	author.Lastname,
	}

	incr := config.DB.Incr("author")
	if incr.Err() != nil {
		config.Error.Println(incr.Err())
		return -1, incr.Err()
	}

	id, err := config.DB.Get("author").Int64()
	if err != nil {
		config.Error.Println(err)
		return -1, err
	}

	result := config.DB.HMSet("author:" + strconv.FormatInt(id, 10), mapAuthor)
	if result.Err() != nil {
		config.Error.Println(result.Err())
		return -1, result.Err()
	}

	return id, nil
}

func GetAuthorsDB() ([]*Author, error) {
	var authors []*Author

	keys := config.DB.Keys("author:*")
	if len(keys.Val()) == 0 {
		config.Info.Println("No authors !!")
		return nil, errors.New("No authors !!")
	}

	for i := 0; i < len(keys.Val()); i++ {
		result := config.DB.HGetAll(keys.Val()[i])
		if result.Err() != nil {
			config.Error.Println(result.Err())
			return nil, result.Err()
		}

		author := &Author{keys.Val()[i], result.Val()["firstname"], result.Val()["lastname"]}

		authors = append(authors, author)
	}

	return authors, nil
}

func GetAuthorDB(id string) (*Author, error) {
	result := config.DB.HGetAll("author:" + id)
	if result.Err() != nil {
		config.Error.Println(result.Err())
		return nil, result.Err()
	} else if len(result.Val()) == 0 {
		config.Error.Println("Author : " + id + " don't exist !!")
		return nil, errors.New("Author : " + id + " don't exist !!")
	}

	author := &Author{Id: "author:" + id, Firstname: result.Val()["firstname"], Lastname: result.Val()["lastname"]}

	return author, nil
}

func UpdateAuthorDB(author *Author) (*Author, error) {
	resultAuthorExist := config.DB.Exists(author.Id)
	if resultAuthorExist.Err() != nil {
		config.Error.Println(resultAuthorExist.Err())
		return nil, resultAuthorExist.Err()
	} else if resultAuthorExist.Val() == false {
		config.Error.Println("Author : " + author.Id + " don't exist !!")
		return nil, errors.New("Author : " + author.Id + " don't exist !!")
	}
	mapAuthor := map[string]string{
		"firstname":	author.Firstname,
		"lastname": 	author.Lastname,
	}

	result := config.DB.HMSet(author.Id, mapAuthor)
	if result.Err() != nil {
		config.Error.Println(result.Err())
		return nil, result.Err()
	}

	return author, nil
}

func DeleteAuthorDB(id string) (bool, error) {
	keys := config.DB.Keys("album:*")
	if len(keys.Val()) == 0 {
		config.Info.Println("No albums !!")
	}

	for i := 0; i < len(keys.Val()); i++ {
		result := config.DB.HGetAll(keys.Val()[i])
		if result.Err() != nil {
			config.Error.Println(result.Err())
			return false, result.Err()
		} else if len(result.Val()) == 0 {
			config.Error.Println("Author : " + id + " don't exist !!")
			return false, errors.New("Author : " + id + " don't exist !!")
		}

		if id == result.Val()["idAuthor"] {
			resultDelAlbum := config.DB.Del(keys.Val()[i])
			if resultDelAlbum.Err() != nil {
				config.Error.Println(resultDelAlbum.Err())
				return false, resultDelAlbum.Err()
			} else if resultDelAlbum.Val() == 0 {
				config.Error.Println("Album : " + keys.Val()[i] + " don't exist !!")
				return false, errors.New("Album : " + keys.Val()[i] + " don't exist !!")
			}
		}
	}

	resultDelAuthor := config.DB.Del("author:" + id)
	if resultDelAuthor.Err() != nil {
		config.Error.Println(resultDelAuthor.Err())
		return false, resultDelAuthor.Err()
	} else if resultDelAuthor.Val() == 0 {
		config.Error.Println("Author : " + id + " don't exist !!")
		return false, errors.New("Author : " + id + " don't exist !!")
	}

	return true, nil
}

func DeleteAllAuthorDB() (bool, error) {
	keys := config.DB.Keys("author:*")
	if len(keys.Val()) == 0 {
		config.Info.Println("No authors !!")
	}

	for i := 0; i < len(keys.Val()); i++ {
		resultDelAuthors := config.DB.Del(keys.Val()[i])
		if resultDelAuthors.Err() != nil {
			config.Error.Println(resultDelAuthors.Err())
			return false, resultDelAuthors.Err()
		}
	}

	resultDelNbAuthor := config.DB.Del("author")
	if resultDelNbAuthor.Err() != nil {
		config.Error.Println(resultDelNbAuthor.Err())
		return false, resultDelNbAuthor.Err()
	}

	return true, nil
}
