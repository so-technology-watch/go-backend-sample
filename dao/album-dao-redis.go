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
var _ AlbumDAO = (*AlbumDAORedis)(nil)

const (
	AlbumStr   = "album"
	AlbumIdStr = "album:"
)

// AlbumDAORedis is the redis implementation of the AlbumDAO
type AlbumDAORedis struct {
	redisCli *redis.Client
}

// NewAlbumDAORedis creates a new AlbumDAO redis implementation
func NewAlbumDAORedis(redisCli *redis.Client) AlbumDAO {
	return &AlbumDAORedis{
		redisCli: redisCli,
	}
}

// Get returns an album by its id
func (dao *AlbumDAORedis) Get(id string) (*model.Album, error) {
	album, err := dao.get(AlbumIdStr + id)
	if err != nil {
		return nil, err
	}

	return album, nil
}

// GetAll returns all albums
func (dao *AlbumDAORedis) GetAll() ([]model.Album, error) {
	var albums []model.Album

	// Collect all albums identifiers
	keys := dao.redisCli.Keys(AlbumIdStr + "*").Val()
	if len(keys) == 0 {
		return nil, errors.New("no albums")
	}

	for i := 0; i < len(keys); i++ {
		// Collect album by identifier
		album, err := dao.get(keys[i])
		if err != nil {
			return nil, err
		}

		albums = append(albums, *album)
	}

	return albums, nil
}

// GetByAuthor returns all albums by author id in parameter
func (dao *AlbumDAORedis) GetByAuthor(authorId string) ([]model.Album, error) {
	var albums []model.Album

	// Collect all albums identifiers
	keys := dao.redisCli.Keys(AlbumIdStr + "*").Val()
	if len(keys) == 0 {
		return nil, errors.New("no albums")
	}

	for i := 0; i < len(keys); i++ {
		// Collect album by identifier
		album, err := dao.get(keys[i])
		if err != nil {
			return nil, err
		}

		if album.AuthorId == authorId {
			albums = append(albums, *album)
		}
	}

	if len(albums) == 0 {
		return nil, errors.New("no albums")
	}

	return albums, nil
}

// Upsert updates or creates an album, returns true if updated, false otherwise or on error
func (dao *AlbumDAORedis) Upsert(album *model.Album) (*model.Album, error) {
	if album.Id == "" {
		// Increment number of albums
		resultNewId, err := dao.redisCli.Incr(AlbumStr).Result()
		if err != nil {
			return nil, err
		}
		album.Id = strconv.FormatInt(resultNewId, 10)
	}

	album, err := dao.save(album)
	if err != nil {
		return nil, err
	}

	return album, nil
}

// Delete deletes an album by its id
func (dao *AlbumDAORedis) Delete(id string) error {
	result, err := dao.redisCli.Del(AlbumIdStr + id).Result()
	if err != nil {
		return err
	} else if result == 0 {
		return errors.New(AlbumIdStr + id + " don't exist")
	}

	return nil
}

// DeleteAll deletes all albums
func (dao *AlbumDAORedis) DeleteAll() error {
	// Collect all identifiers of albums
	keys := dao.redisCli.Keys(AlbumIdStr + "*").Val()
	if len(keys) == 0 { // If no albums in database
		utils.LogWarning.Println("no albums")
	}

	for i := 0; i < len(keys); i++ {
		// Deletion of album by identifier
		_, err := dao.redisCli.Del(keys[i]).Result()
		if err != nil {
			return err
		}
	}

	// Delete number of albums in database
	resultDelNbAlbum := dao.redisCli.Del(AlbumStr)
	if resultDelNbAlbum.Err() != nil {
		return resultDelNbAlbum.Err()
	}

	return nil
}

// Exist checks if the album exist
func (dao *AlbumDAORedis) Exist(id string) (bool, error) {
	result, err := dao.redisCli.Exists(AlbumIdStr + id).Result()
	if err != nil {
		return false, err
	} else if result == false {
		return false, errors.New(AlbumIdStr + id + " don't exist")
	}
	return result, nil
}

// Save saves the album
func (dao *AlbumDAORedis) save(album *model.Album) (*model.Album, error) {
	result, err := json.Marshal(album)
	if err != nil {
		return nil, err
	}

	// Save album with songs in database
	status := dao.redisCli.Set(AlbumIdStr+album.Id, string(result), 0)
	if status.Err() != nil {
		return nil, status.Err()
	}

	return album, nil
}

func (dao *AlbumDAORedis) get(id string) (*model.Album, error) {
	result, err := dao.redisCli.Get(id).Result()
	if err != nil {
		return nil, err
	} else if len(result) == 0 {
		return nil, errors.New(id + " don't exist !!")
	}

	album := model.Album{}
	err = json.Unmarshal([]byte(result), &album)
	if err != nil {
		return nil, err
	}

	return &album, nil
}
