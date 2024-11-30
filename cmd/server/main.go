package main

import (
	"log"
	"net/http"

	"github.com/ramonamorim/go-rate-limiter/internal/application"
	"github.com/ramonamorim/go-rate-limiter/internal/domain/service"
	"github.com/ramonamorim/go-rate-limiter/internal/infra/config"
	"github.com/ramonamorim/go-rate-limiter/internal/infra/limiter"
	"github.com/ramonamorim/go-rate-limiter/internal/infra/middleware"
)

func main() {
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}
	log.Printf("Redis configuration: host=%s, port=%s, password=%s", cfg.RedisHost, cfg.RedisPort, cfg.RedisPassword)

	redisAddr := cfg.RedisHost + ":" + cfg.RedisPort
	ipILimiter := limiter.NewRedis(redisAddr, cfg.RedisPassword)
	tokenILimiter := limiter.NewRedis(redisAddr, cfg.RedisPassword)

	rateLimiter := limiter.NewLimiter(
		cfg.IPMaxRequestsPerSecond,
		cfg.IPBlockDurationSeconds,
		cfg.TokenMaxRequestsPerSecond,
		cfg.TokenBlockDurationSeconds,
		ipILimiter,
		tokenILimiter,
	)

	rateLimiterService := service.NewRateLimiterService(rateLimiter)
	rateLimiterApp := application.NewRateLimiterApp(rateLimiterService)
	rateLimiterMiddleware := middleware.NewRateLimiterMiddleware(rateLimiterApp)

	mux := http.NewServeMux()
	mux.Handle("/", rateLimiterMiddleware.Handler(http.HandlerFunc(handler)))

	log.Println("Server running on port 8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}
