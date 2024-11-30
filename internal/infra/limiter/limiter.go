package limiter

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

type Limiter struct {
	IPMaxRequestsPerSecond    int
	IPBlockDurationSeconds    int
	TokenMaxRequestsPerSecond int
	TokenBlockDurationSeconds int
	IPILimiter                ILimiter
	TokenILimiter             ILimiter
}

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

func (l *Limiter) AllowRequest(ip, token string) bool {
	if token != "" {
		return l.checkRateLimit(fmt.Sprintf("ratelimit:token:%s", token), l.TokenMaxRequestsPerSecond, l.TokenBlockDurationSeconds, l.TokenILimiter)
	}
	return l.checkRateLimit(fmt.Sprintf("ratelimit:ip:%s", ip), l.IPMaxRequestsPerSecond, l.IPBlockDurationSeconds, l.IPILimiter)
}

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

func (l *Limiter) getRequestCount(key string, limiter ILimiter) (int, error) {
	countStr, err := limiter.Get(key)
	if err != nil {
		return 0, err
	}

	count, _ := strconv.Atoi(countStr)
	log.Printf("Key: %s, Current Count: %d", key, count)
	return count, nil
}

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
