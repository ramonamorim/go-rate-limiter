package middleware

import (
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/ramonamorim/go-rate-limiter/internal/application"
)

type RateLimiterMiddleware struct {
	RateLimiterApp application.RateLimiterServiceInterface
}

func NewRateLimiterMiddleware(rateLimiterApp application.RateLimiterServiceInterface) *RateLimiterMiddleware {
	return &RateLimiterMiddleware{
		RateLimiterApp: rateLimiterApp,
	}
}

func (rlm *RateLimiterMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := ExtractClientIP(r)
		token := r.Header.Get("API_KEY")

		log.Printf("Incoming request - IP: %s, Token: %s", ip, token)

		if rlm.RateLimiterApp.AllowRequest(ip, token) {
			next.ServeHTTP(w, r)
		} else {
			log.Printf("Request denied - IP: %s, Token: %s", ip, token)
			http.Error(w, "You have reached the maximum number of allowed requests", http.StatusTooManyRequests)
		}
	})
}

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
