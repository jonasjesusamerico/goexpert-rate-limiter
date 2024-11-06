package limiter

import "time"

type ILimiter interface {
	Get(key string) (string, error)
	Set(key string, value interface{}, expiration time.Duration) error
	Incr(key string) error
	Expire(key string, expiration time.Duration) error
}
