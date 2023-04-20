package userstore

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(host *string, password *string, db *int) *RedisStore {
	client := redis.NewClient(&redis.Options{
		Addr:     *host,
		Password: *password,
		DB:       *db,
	})
	return &RedisStore{
		client: client,
	}
}

func (r *RedisStore) Put(user *User) error {
	err := r.client.Set(context.Background(), user.Username, time.Time(user.DateOfBirth).Format("2006-01-02"), 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisStore) Get(username string) (*User, error) {
	val, err := r.client.Get(context.Background(), username).Result()
	// If it does not exist, return no user and no error
	// If there is any other error, return that error
	if errors.Is(err, redis.Nil) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	dob, err := time.Parse("2006-01-02", val)
	if err != nil {
		return nil, err
	}

	return &User{
		Username:    username,
		DateOfBirth: DateOfBirth(dob),
	}, nil
}
