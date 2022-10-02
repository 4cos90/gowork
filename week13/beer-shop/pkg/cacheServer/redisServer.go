package cacheServer

import (
	"github.com/go-redis/redis"
)

type redisManager struct {
	Name   string
	client *redis.Client
}

type RedisServer interface {
	Get(key string) (string, error)
	Set(key string, value string) error
}

func NewClient(name string, addr string, password string, db int) RedisServer {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &redisManager{
		Name:   name,
		client: client,
	}
}

func (s *redisManager) Get(key string) (string, error) {
	return s.client.Get(key).Result()
}

func (s *redisManager) Set(key string, value string) error {
	return s.client.Set(key, value, 0).Err()
}
