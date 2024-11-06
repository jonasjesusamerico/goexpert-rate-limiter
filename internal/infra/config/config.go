package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	BlockDurationSeconds      int
	MaxRequestsPerSecond      int
	TokenMaxRequestsPerSecond int
	IPBlockDurationSeconds    int
	IPMaxRequestsPerSecond    int
	RedisPassword             string
	RedisHost                 string
	RedisPort                 string
	TokenBlockDurationSeconds int
}

func LoadConfig(envFile string) (*Config, error) {
	if err := godotenv.Load(envFile); err != nil {
		log.Printf("Warning: %v\n", err)
	}

	redisHost := getEnv("REDIS_HOST", "localhost")
	redisPort := getEnv("REDIS_PORT", "6379")
	redisPassword := getEnv("REDIS_PASSWORD", "")

	blockDurationSeconds, _ := strconv.Atoi(getEnv("BLOCK_DURATION_SECONDS", "60"))
	maxRequestsPerSecond, _ := strconv.Atoi(getEnv("MAX_REQUESTS_PER_SECOND", "100"))
	tokenMaxRequestsPerSecond, _ := strconv.Atoi(getEnv("TOKEN_MAX_REQUESTS_PER_SECOND", "10"))
	ipBlockDurationSeconds, _ := strconv.Atoi(getEnv("IP_BLOCK_DURATION_SECONDS", "300"))
	tokenBlockDurationSeconds, _ := strconv.Atoi(getEnv("TOKEN_BLOCK_DURATION_SECONDS", "300"))
	ipMaxRequestsPerSecond, _ := strconv.Atoi(getEnv("IP_MAX_REQUESTS_PER_SECOND", "5"))

	cfg := &Config{
		MaxRequestsPerSecond:      maxRequestsPerSecond,
		TokenMaxRequestsPerSecond: tokenMaxRequestsPerSecond,
		IPMaxRequestsPerSecond:    ipMaxRequestsPerSecond,
		BlockDurationSeconds:      blockDurationSeconds,
		RedisPassword:             redisPassword,
		RedisPort:                 redisPort,
		RedisHost:                 redisHost,
		TokenBlockDurationSeconds: tokenBlockDurationSeconds,
		IPBlockDurationSeconds:    ipBlockDurationSeconds,
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
