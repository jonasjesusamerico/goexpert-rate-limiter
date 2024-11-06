package limiter

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// Redis is a wrapper around the Redis client for rate limiting operations.
type Redis struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedis initializes a new Redis with the specified address and password.
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

// Get retrieves a value from Redis for the given key.
func (rs *Redis) Get(key string) (string, error) {
	val, err := rs.client.Get(rs.ctx, key).Result()
	if err == redis.Nil {
		return "", nil // Key does not exist
	}
	return val, err
}

// Set sets a value in Redis with an expiration time.
func (rs *Redis) Set(key string, value interface{}, expiration time.Duration) error {
	return rs.client.Set(rs.ctx, key, value, expiration).Err()
}

// Incr increments the value of a key in Redis.
func (rs *Redis) Incr(key string) error {
	return rs.client.Incr(rs.ctx, key).Err()
}

// Expire sets an expiration time on a key in Redis.
func (rs *Redis) Expire(key string, expiration time.Duration) error {
	return rs.client.Expire(rs.ctx, key, expiration).Err()
}
