package db

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
)

type Database struct {
	Client *redis.Client
}

var (
	ErrNil = errors.New("No such key")
	Ctx    = context.Background()
)

// db instance
var instance *Database

// GetInstance returns the singleton instance
func GetInstance() *Database {
	if instance == nil {
		instance = &Database{
			Client: redis.NewClient(&redis.Options{
				Addr: "localhost:6379",
			}),
		}
	}
	return instance
}

// NewDatabase returns a new instance of Database
func NewDatabase(address string) (*Database, error) {
	client := redis.NewClient(&redis.Options{
		Addr: address,
	})

	if err := client.Ping(Ctx).Err(); err != nil {
		return nil, err
	}

	return &Database{Client: client}, nil
}

// Get returns the value of the key
func (db *Database) Get(key string) (string, error) {
	val, err := db.Client.Get(Ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

// GetAndAdd returns the value of the key and adds the value to the key
func (db *Database) GetAndAdd(key string, value string) error {
	val, err := db.Client.Get(Ctx, key).Result()
	if err != nil {
		return err
	}
	val = val + value
	err = db.Client.Set(Ctx, key, val, 0).Err()
	if err != nil {
		return err
	}
	return nil
}
