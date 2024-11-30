package limiter

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedis(addr, password string) *Redis {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	return &Redis{
		client: client,
		ctx:    context.Background(),
	}
}

func (rs *Redis) Get(key string) (string, error) {
	val, err := rs.client.Get(rs.ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

func (rs *Redis) Set(key string, value interface{}, expiration time.Duration) error {
	return rs.client.Set(rs.ctx, key, value, expiration).Err()
}

func (rs *Redis) Incr(key string) error {
	return rs.client.Incr(rs.ctx, key).Err()
}

func (rs *Redis) Expire(key string, expiration time.Duration) error {
	return rs.client.Expire(rs.ctx, key, expiration).Err()
}
