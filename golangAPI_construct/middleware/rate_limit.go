package middleware

import (
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"golangAPI_construct/responses"
)

type bucket struct {
	tokens     float64
	lastRefill time.Time
}

type rateLimiter struct {
	mu    sync.Mutex
	bkts  map[string]*bucket
	rate  float64
	burst float64
	ttl   time.Duration
}

func newRateLimiter(rate, burst float64, ttl time.Duration) *rateLimiter {
	rl := &rateLimiter{
		bkts:  make(map[string]*bucket),
		rate:  rate,
		burst: burst,
		ttl:   ttl,
	}
	go rl.janitor()
	return rl
}

func (r *rateLimiter) janitor() {
	t := time.NewTicker(r.ttl)
	for range t.C {
		now := time.Now()
		r.mu.Lock()
		for k, b := range r.bkts {
			if now.Sub(b.lastRefill) > r.ttl {
				delete(r.bkts, k)
			}
		}
		r.mu.Unlock()
	}
}

func (r *rateLimiter) allow(key string) bool {
	now := time.Now()
	r.mu.Lock()
	defer r.mu.Unlock()
	b, ok := r.bkts[key]
	if !ok {
		b = &bucket{tokens: r.burst, lastRefill: now}
		r.bkts[key] = b
	}
	elapsed := now.Sub(b.lastRefill).Seconds()
	if elapsed > 0 {
		b.tokens += elapsed * r.rate
		if b.tokens > r.burst {
			b.tokens = r.burst
		}
		b.lastRefill = now
	}
	if b.tokens >= 1 {
		b.tokens -= 1
		return true
	}
	return false
}

func RateLimit() func(http.Handler) http.Handler {
	rps := 5.0
	burst := 10.0
	if v := os.Getenv("RATE_LIMIT_RPS"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil && f > 0 {
			rps = f
		}
	}
	if v := os.Getenv("RATE_LIMIT_BURST"); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil && f >= 1 {
			burst = f
		}
	}
	rl := newRateLimiter(rps, burst, 5*time.Minute)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := clientIP(r)
			if !rl.allow(ip) {
				responses.Fail(w, r, responses.NewAppError(http.StatusTooManyRequests, "RATE_LIMIT", "too many requests"))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func clientIP(r *http.Request) string {
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		// 只取第一個
		host := fwd
		if idx := len(fwd); idx > 0 {
			host = fwd
		}
		if ip := net.ParseIP(host); ip != nil {
			return ip.String()
		}
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
