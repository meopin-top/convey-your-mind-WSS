package db

import "github.com/redis/go-redis/v9"

func (db *Database) Publish(channel string, message interface{}) error {
	return db.Client.Publish(Ctx, channel, message).Err()
}

func (db *Database) Subscribe(channel string) *redis.PubSub {
	return db.Client.Subscribe(Ctx, channel)
}

func (db *Database) Broadcast(channel string, message interface{}) error {
	return db.Client.Publish(Ctx, channel, message).Err()
}

func (db *Database) Get(key string) (string, error) {
	return db.Client.Get(Ctx, key).Result()
}

func (db *Database) Set(key string, value interface{}) error {

	return db.Client.Set(Ctx, key, value, 0).Err()
}
