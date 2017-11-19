package dao

import (
	"encoding/json"
	"errors"
	"go-backend-sample/model"
	"go-backend-sample/utils"
	"gopkg.in/redis.v5"
	"strconv"
)

// compilation time interface check
var _ AuthorDAO = (*AuthorDAORedis)(nil)

const (
	AuthorStr   = "author"
	AuthorIdStr = "author:"
)

// AuthorDAORedis is the redis implementation of the AuthorDAO
type AuthorDAORedis struct {
	redisCli *redis.Client
}

// NewAuthorDAORedis creates a new AuthorDAO redis implementation
func NewAuthorDAORedis(redisCli *redis.Client) AuthorDAO {
	return &AuthorDAORedis{
		redisCli: redisCli,
	}
}

// Get returns an author by its id
func (dao *AuthorDAORedis) Get(id string) (*model.Author, error) {
	author, err := dao.get(AuthorIdStr + id)
	if err != nil {
		return nil, err
	}

	return author, nil
}

// GetAll returns all authors
func (dao *AuthorDAORedis) GetAll() ([]model.Author, error) {
	var authors []model.Author

	// Collect all authors identifiers
	keys := dao.redisCli.Keys(AuthorIdStr + "*").Val()
	if len(keys) == 0 {
		return nil, errors.New("no authors")
	}

	for i := 0; i < len(keys); i++ {
		// Collect author by identifier
		author, err := dao.get(keys[i])
		if err != nil {
			return nil, err
		}

		authors = append(authors, *author)
	}

	return authors, nil
}

// Upsert updates or creates an author, returns true if updated, false otherwise or on error
func (dao *AuthorDAORedis) Upsert(author *model.Author) (*model.Author, error) {
	if author.Id == "" {
		// Increment number of authors
		resultNewId, err := dao.redisCli.Incr(AuthorStr).Result()
		if err != nil {
			return nil, err
		}
		author.Id = strconv.FormatInt(resultNewId, 10)
	}

	author, err := dao.save(author)
	if err != nil {
		return nil, err
	}

	return author, nil
}

// Delete deletes an author by its id
func (dao *AuthorDAORedis) Delete(id string) error {
	result, err := dao.redisCli.Del(AuthorIdStr + id).Result()
	if err != nil {
		return err
	} else if result == 0 {
		return errors.New(AuthorIdStr + id + " don't exist")
	}

	return nil
}

// DeleteAll deletes all authors
func (dao *AuthorDAORedis) DeleteAll() error {
	// Collect all identifiers of authors
	keys := dao.redisCli.Keys(AuthorIdStr + "*").Val()
	if len(keys) == 0 { // If no authors in database
		utils.LogWarning.Println("no authors")
	}

	for i := 0; i < len(keys); i++ {
		// Deletion of author by identifier
		_, err := dao.redisCli.Del(keys[i]).Result()
		if err != nil {
			return err
		}
	}

	// Delete number of authors in database
	resultDelNbAuthor := dao.redisCli.Del(AuthorStr)
	if resultDelNbAuthor.Err() != nil {
		return resultDelNbAuthor.Err()
	}

	return nil
}

// Exist checks if the author exist
func (dao *AuthorDAORedis) Exist(id string) (bool, error) {
	result, err := dao.redisCli.Exists(AuthorIdStr + id).Result()
	if err != nil {
		return false, err
	} else if result == false {
		return false, errors.New(AuthorIdStr + id + " don't exist")
	}
	return result, nil
}

// Save saves the author
func (dao *AuthorDAORedis) save(author *model.Author) (*model.Author, error) {
	result, err := json.Marshal(author)
	if err != nil {
		return nil, err
	}

	// Save author with songs in database
	status := dao.redisCli.Set(AuthorIdStr+author.Id, string(result), 0)
	if status.Err() != nil {
		return nil, status.Err()
	}

	return author, nil
}

func (dao *AuthorDAORedis) get(id string) (*model.Author, error) {
	result, err := dao.redisCli.Get(id).Result()
	if err != nil {
		return nil, err
	} else if len(result) == 0 {
		return nil, errors.New(id + " don't exist")
	}

	author := model.Author{}
	err = json.Unmarshal([]byte(result), &author)
	if err != nil {
		return nil, err
	}

	return &author, nil
}
