package middleware

import (
	"net/http"
	"sync"
	"time"

	"lms/internal/utils"
)

type clientRecord struct {
	lastSeen time.Time
	tokens   int
}

type RateLimiter struct {
	mu      sync.Mutex
	clients map[string]*clientRecord
	limit   int
	window  time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		clients: make(map[string]*clientRecord),
		limit:   limit,
		window:  window,
	}
}

func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	record, exists := rl.clients[ip]

	if !exists || now.Sub(record.lastSeen) > rl.window {
		rl.clients[ip] = &clientRecord{lastSeen: now, tokens: rl.limit - 1}
		return true
	}

	if record.tokens > 0 {
		record.tokens--
		record.lastSeen = now
		return true
	}

	return false
}

// RateLimitMiddleware blocks IPs exceeding rate thresholds
func RateLimitMiddleware(rl *RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			if !rl.Allow(ip) {
				utils.RateLimitExceeded(w, "Terlalu banyak permintaan. Silakan tunggu beberapa saat.")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
