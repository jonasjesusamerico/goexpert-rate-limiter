package model

type RateLimiter struct {
	IPBlockDurationSeconds    int
	IPMaxRequestsPerSecond    int
	MaxRequestsPerSecond      int
	TokenBlockDurationSeconds int
	TokenMaxRequestsPerSecond int
	BlockDurationSeconds      int
}
