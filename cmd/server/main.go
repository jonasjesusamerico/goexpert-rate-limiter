package main

import (
	"log"
	"net/http"

	"github.com/jonasjesusamerico/goexpert-rate-limiter/internal/application"
	"github.com/jonasjesusamerico/goexpert-rate-limiter/internal/domain/service"
	"github.com/jonasjesusamerico/goexpert-rate-limiter/internal/infra/config"
	"github.com/jonasjesusamerico/goexpert-rate-limiter/internal/infra/limiter"
	"github.com/jonasjesusamerico/goexpert-rate-limiter/internal/infra/middleware"
)

func main() {
	// Load configuration from the .env file
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}
	log.Printf("Redis configuration: host=%s, port=%s, password=%s", cfg.RedisHost, cfg.RedisPort, cfg.RedisPassword)

	// Setup Redis limiters for IP and token rate limiting
	redisAddr := cfg.RedisHost + ":" + cfg.RedisPort
	ipILimiter := limiter.NewRedis(redisAddr, cfg.RedisPassword)
	tokenILimiter := limiter.NewRedis(redisAddr, cfg.RedisPassword)

	// Initialize rate limiter with configured settings
	rateLimiter := limiter.NewLimiter(
		cfg.IPMaxRequestsPerSecond,
		cfg.IPBlockDurationSeconds,
		cfg.TokenMaxRequestsPerSecond,
		cfg.TokenBlockDurationSeconds,
		ipILimiter,
		tokenILimiter,
	)

	// Setup service, application, and middleware layers
	rateLimiterService := service.NewRateLimiterService(rateLimiter)
	rateLimiterApp := application.NewRateLimiterApp(rateLimiterService)
	rateLimiterMiddleware := middleware.NewRateLimiterMiddleware(rateLimiterApp)

	// Configure the HTTP server with the rate limiter middleware
	mux := http.NewServeMux()
	mux.Handle("/", rateLimiterMiddleware.Handler(http.HandlerFunc(handler)))

	// Start the HTTP server
	log.Println("Server running on port 8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

// Simple HTTP handler function
func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}
