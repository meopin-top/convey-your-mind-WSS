package redis

import (
	"context"
	"errors"
	"log"

	"github.com/redis/go-redis/v9"
)

// Database is the data structure for the redis database
type Database struct {
	Client *redis.Client
}

var (
	lockSuffix = ":lock"
	// ErrNil is returned when key does not exist.
	ErrNil = errors.New("No such key")
	// Ctx is the context
	Ctx = context.Background()
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
		log.Fatal(err)
		return nil, err
	}

	return &Database{Client: client}, nil
}

// Get returns the value of the key
func (db *Database) Get(key string) (string, error) {
	val, err := db.Client.Get(Ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", ErrNil
		}
		return "", err
	}
	return val, nil
}

// GetLock returns the value of the key
func (db *Database) GetLock(key string) (string, error) {
	val, err := db.Client.Get(Ctx, key+lockSuffix).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", err
	}
	return val, nil
}

// IncrLock increases the value of the key
func (db *Database) IncrLock(key string) error {
	if err := db.Client.Incr(Ctx, key+lockSuffix).Err(); err != nil {
		return err
	}
	return nil
}

// DecrLock decreases the value of the key
func (db *Database) DecrLock(key string) error {
	if err := db.Client.Decr(Ctx, key+lockSuffix).Err(); err != nil {
		return err
	}
	return nil
}

// Add returns the value of the key and adds the value to the key
func (db *Database) Add(key string, value string) error {
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

// Set sets the value of the key
func (db *Database) Set(key string, value string) error {
	err := db.Client.Set(Ctx, key, value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}
