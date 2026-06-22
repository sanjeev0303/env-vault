package handler

import (
	"net/http"
	"sync"
	"time"

	"env-vault/server/internal/domain"
)

type rateLimiter struct {
	tokens    int
	lastSeen  time.Time
}

type RateLimiter struct {
	mu           sync.Mutex
	visitors     map[string]*rateLimiter
	rate         int
	burst        int
	cleanupDelay time.Duration
}

// NewRateLimiter creates a new token bucket rate limiter.
// rate: tokens added per second
// burst: maximum tokens bucket can hold
func NewRateLimiter(rate, burst int) *RateLimiter {
	rl := &RateLimiter{
		visitors:     make(map[string]*rateLimiter),
		rate:         rate,
		burst:        burst,
		cleanupDelay: 3 * time.Minute,
	}

	go rl.cleanupLoop()
	return rl
}

func (rl *RateLimiter) cleanupLoop() {
	for {
		time.Sleep(rl.cleanupDelay)
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > rl.cleanupDelay {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// Limit is the middleware that applies the rate limit
func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
			ip = forwarded
		}

		rl.mu.Lock()
		v, exists := rl.visitors[ip]
		if !exists {
			rl.visitors[ip] = &rateLimiter{
				tokens:   rl.burst - 1,
				lastSeen: time.Now(),
			}
			rl.mu.Unlock()
			next.ServeHTTP(w, r)
			return
		}

		now := time.Now()
		elapsed := now.Sub(v.lastSeen).Seconds()
		v.tokens += int(elapsed * float64(rl.rate))
		if v.tokens > rl.burst {
			v.tokens = rl.burst
		}

		if v.tokens <= 0 {
			v.lastSeen = now
			rl.mu.Unlock()
			RespondWithError(w, domain.ErrRateLimited)
			return
		}

		v.tokens--
		v.lastSeen = now
		rl.mu.Unlock()

		next.ServeHTTP(w, r)
	})
}
