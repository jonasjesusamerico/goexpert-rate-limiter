package limiter

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

// Limiter handles rate limiting for both IP and Token-based requests.
type Limiter struct {
	IPMaxRequestsPerSecond    int
	IPBlockDurationSeconds    int
	TokenMaxRequestsPerSecond int
	TokenBlockDurationSeconds int
	IPILimiter                ILimiter
	TokenILimiter             ILimiter
}

// NewLimiter creates a new Limiter instance with provided configurations and limiters.
func NewLimiter(ipMaxRequestsPerSecond, ipBlockDurationSeconds, tokenMaxRequestsPerSecond, tokenBlockDurationSeconds int, ipILimiter, tokenILimiter ILimiter) *Limiter {
	return &Limiter{
		IPMaxRequestsPerSecond:    ipMaxRequestsPerSecond,
		IPBlockDurationSeconds:    ipBlockDurationSeconds,
		TokenMaxRequestsPerSecond: tokenMaxRequestsPerSecond,
		TokenBlockDurationSeconds: tokenBlockDurationSeconds,
		IPILimiter:                ipILimiter,
		TokenILimiter:             tokenILimiter,
	}
}

// AllowRequest determines if a request should be allowed based on either IP or Token.
func (l *Limiter) AllowRequest(ip, token string) bool {
	if token != "" {
		return l.checkRateLimit(fmt.Sprintf("ratelimit:token:%s", token), l.TokenMaxRequestsPerSecond, l.TokenBlockDurationSeconds, l.TokenILimiter)
	}
	return l.checkRateLimit(fmt.Sprintf("ratelimit:ip:%s", ip), l.IPMaxRequestsPerSecond, l.IPBlockDurationSeconds, l.IPILimiter)
}

// checkRateLimit checks and updates the request count for a given key in the limiter.
func (l *Limiter) checkRateLimit(key string, maxRequestsPerSecond, blockDurationSeconds int, limiter ILimiter) bool {
	count, err := l.getRequestCount(key, limiter)
	if err != nil {
		log.Printf("Error retrieving request count: %v", err)
		return false
	}

	if count >= maxRequestsPerSecond {
		log.Printf("Rate limit exceeded for key: %s", key)
		return false
	}

	if err := l.updateRequestCount(key, blockDurationSeconds, limiter); err != nil {
		log.Printf("Error updating request count: %v", err)
		return false
	}

	log.Printf("Request allowed for key: %s", key)
	return true
}

// getRequestCount retrieves the current request count from the limiter.
func (l *Limiter) getRequestCount(key string, limiter ILimiter) (int, error) {
	countStr, err := limiter.Get(key)
	if err != nil {
		return 0, err
	}

	count, _ := strconv.Atoi(countStr) // Ignoring error since default 0 is okay
	log.Printf("Key: %s, Current Count: %d", key, count)
	return count, nil
}

// updateRequestCount increments the request count and sets the expiration in the limiter.
func (l *Limiter) updateRequestCount(key string, blockDurationSeconds int, limiter ILimiter) error {
	if err := limiter.Incr(key); err != nil {
		return fmt.Errorf("error incrementing key: %s, %v", key, err)
	}

	expiration := time.Duration(blockDurationSeconds) * time.Second
	if err := limiter.Expire(key, expiration); err != nil {
		return fmt.Errorf("error setting expiration for key: %s, %v", key, err)
	}

	return nil
}
