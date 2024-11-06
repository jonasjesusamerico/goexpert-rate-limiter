package application

import (
	"github.com/jonasjesusamerico/goexpert-rate-limiter/internal/domain/service"
)

type RateLimiterApp struct {
	RateLimiterService *service.RateLimiterService
}
type RateLimiterServiceInterface interface {
	AllowRequest(ip, token string) bool
}

func (app *RateLimiterApp) AllowRequest(ip, token string) bool {
	return app.RateLimiterService.AllowRequest(ip, token)
}

func NewRateLimiterApp(rateLimiterService *service.RateLimiterService) *RateLimiterApp {
	return &RateLimiterApp{
		RateLimiterService: rateLimiterService,
	}
}
