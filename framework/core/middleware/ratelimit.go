package middleware

import (
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	rate     int
	burst    int
}

type visitor struct {
	tokens   int
	lastTime time.Time
}

func NewRateLimiter(rate, burst int) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		rate:     rate,
		burst:    burst,
	}
	go rl.cleanup(5 * time.Minute)
	return rl
}

func (rl *RateLimiter) cleanup(interval time.Duration) {
	for {
		time.Sleep(interval)
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastTime) > interval {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		rl.mu.Lock()
		v, exists := rl.visitors[ip]
		if !exists {
			v = &visitor{tokens: rl.burst, lastTime: time.Now()}
			rl.visitors[ip] = v
		}
		elapsed := time.Since(v.lastTime)
		v.tokens += int(elapsed.Seconds()) * rl.rate
		if v.tokens > rl.burst {
			v.tokens = rl.burst
		}
		v.lastTime = time.Now()
		if v.tokens <= 0 {
			rl.mu.Unlock()
			w.Header().Set("Retry-After", "1")
			http.Error(w, `{"error":{"code":"TOO_MANY_REQUESTS","message":"rate limit exceeded"}}`, http.StatusTooManyRequests)
			return
		}
		v.tokens--
		rl.mu.Unlock()
		next.ServeHTTP(w, r)
	})
}
