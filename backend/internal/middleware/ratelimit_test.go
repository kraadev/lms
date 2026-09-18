package middleware

import (
	"testing"
	"time"
)

func TestRateLimiter(t *testing.T) {
	rl := NewRateLimiter(3, 1*time.Second)
	ip := "192.168.1.100"

	// 1-3 allowed
	if !rl.Allow(ip) || !rl.Allow(ip) || !rl.Allow(ip) {
		t.Errorf("expected first 3 requests to be allowed")
	}

	// 4th rejected
	if rl.Allow(ip) {
		t.Errorf("expected 4th request within window to be blocked")
	}

	// Different IP allowed
	if !rl.Allow("192.168.1.101") {
		t.Errorf("expected different IP to be allowed")
	}
}
