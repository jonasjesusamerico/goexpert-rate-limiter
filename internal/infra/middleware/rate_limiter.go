package middleware

import (
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/jonasjesusamerico/goexpert-rate-limiter/internal/application"
)

// RateLimiterMiddleware define a estrutura para aplicar o rate limiting
type RateLimiterMiddleware struct {
	RateLimiterApp application.RateLimiterServiceInterface
}

// NewRateLimiterMiddleware cria uma nova instância do middleware de rate limiting
func NewRateLimiterMiddleware(rateLimiterApp application.RateLimiterServiceInterface) *RateLimiterMiddleware {
	return &RateLimiterMiddleware{
		RateLimiterApp: rateLimiterApp,
	}
}

// Handler aplica o rate limiting ao manipulador HTTP fornecido
func (rlm *RateLimiterMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := ExtractClientIP(r)
		token := r.Header.Get("API_KEY")

		log.Printf("Incoming request - IP: %s, Token: %s", ip, token)

		// Verifica se a requisição é permitida pelo Rate Limiter
		if rlm.RateLimiterApp.AllowRequest(ip, token) {
			next.ServeHTTP(w, r)
		} else {
			log.Printf("Request denied - IP: %s, Token: %s", ip, token)
			http.Error(w, "You have reached the maximum number of allowed requests", http.StatusTooManyRequests)
		}
	})
}

// extractClientIP obtém o IP do cliente a partir da solicitação
func ExtractClientIP(r *http.Request) string {
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		return strings.Split(forwarded, ",")[0]
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ""
	}
	return host
}
