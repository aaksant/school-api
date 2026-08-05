package middlewares

import (
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// Per-IP token bucket implementation

// This rate limiter combines IP-based client tracking with the
// token bucket algorithm provided by golang.org/x/time/rate.

// Each client IP address is associated with its own rate.Limiter,
// allowing every client to maintain an independent token bucket.
// A request is permitted only if the client's limiter has an
// available token; otherwise, the request is rejected with 429.

// New limiters are created when a client is first seen and
// stored in an in memory map. Each visitor also records its last
// activity time, allowing inactive entries to be periodically
// removed by a background cleanup goroutine to prevent unbounded
// memory growth.

// Compared to a fixed window counter, this approach provides
// smoother rate limiting, supports configurable burst capacity,
// and avoids request spikes at window boundaries. It also ensures
// that one client's traffic does not consume another client's
// available request quota.

var (
	cleanupInterval = time.Minute
	visitorTtl      = 3 * time.Minute
)

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type ipRateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	rate     rate.Limit
	burst    int
}

func (iprl *ipRateLimiter) getVisitor(ip string) *rate.Limiter {
	iprl.mu.Lock()
	defer iprl.mu.Unlock()

	v, isVisitorExists := iprl.visitors[ip]
	if !isVisitorExists {
		limiter := rate.NewLimiter(
			iprl.rate,
			iprl.burst,
		)
		iprl.visitors[ip] = &visitor{
			limiter:  limiter,
			lastSeen: time.Now(),
		}
		return limiter
	}

	v.lastSeen = time.Now()
	return v.limiter
}

func (iprl *ipRateLimiter) cleanup() {
	for {
		time.Sleep(cleanupInterval)
		iprl.mu.Lock()
		for ip, v := range iprl.visitors {
			if time.Since(v.lastSeen) > visitorTtl {
				delete(iprl.visitors, ip)
			}
		}
		iprl.mu.Unlock()
	}
}

func NewIpRateLimiter(rate rate.Limit, burst int) *ipRateLimiter {
	iprl := &ipRateLimiter{
		visitors: make(map[string]*visitor),
		rate:     rate,
		burst:    burst,
	}

	go iprl.cleanup()
	return iprl
}

func (iprl *ipRateLimiter) RateLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}

		limiter := iprl.getVisitor(ip)

		if !limiter.Allow() {
			w.Header().Set("Retry-After", "1")
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// IP counter implementation

// This rate limiter implements a fixed window counter algorithm.

// Each client IP address is associated with a request counter stored
// in an in memory map. Every incoming request increments the counter
// for that IP. If the counter exceeds the configured limit, the request
// is rejected with 429.

// All counters are reset simultaneously after a fixed time interval by
// a background goroutine, starting a new counting window.

// Limitations:
//   - Counters are stored only in memory and are lost when the server
//     restarts.
//   - The fixed window algorithm allows request bursts near window
//     boundaries.
//   - A mutex is required to protect concurrent access to the visitor
//     map. The mutex should only guard the map operations and must not
//     be held while processing the HTTP request.

// type rateLimiter struct {
// 	mu        sync.Mutex
// 	visitors  map[string]int
// 	limit     int
// 	resetTime time.Duration
// }

// func (rl *rateLimiter) reset() {
// 	for {
// 		time.Sleep(rl.resetTime)
// 		rl.mu.Lock()
// 		rl.visitors = make(map[string]int)
// 		rl.mu.Unlock()
// 	}
// }

// func NewRateLimiter(limit int, resetTime time.Duration) *rateLimiter {
// 	rl := &rateLimiter{
// 		visitors:  make(map[string]int),
// 		limit:     limit,
// 		resetTime: resetTime,
// 	}

// 	go rl.reset()
// 	return rl
// }

// func (rl *rateLimiter) DeferredRateLimiter(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		rl.mu.Lock()
// 		defer rl.mu.Unlock()
// 		visitorIp := r.RemoteAddr
// 		rl.visitors[visitorIp]++

// 		if rl.visitors[visitorIp] > rl.limit {
// 			http.Error(w, "Too many request", http.StatusTooManyRequests)
// 			return
// 		}

// 		next.ServeHTTP(w, r)
// 	})
// }

// func (rl *rateLimiter) RateLimiter(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		rl.mu.Lock()
// 		visitorIP := r.RemoteAddr
// 		rl.visitors[visitorIP]++
// 		count := rl.visitors[visitorIP]
// 		rl.mu.Unlock()

// 		if count > rl.limit {
// 			http.Error(w, "Too many requests", http.StatusTooManyRequests)
// 			return
// 		}

// 		next.ServeHTTP(w, r)
// 	})
// }

// Token bucket implementation

// This rate limiter uses the golang.org/x/time/rate package, which
// implements the token bucket algorithm.

// Tokens are added to the bucket at a fixed rate, up to a configured
// burst capacity. Each incoming request consumes one token. If a token
// is available, the request is allowed immediately. Otherwise, the
// request is rejected with 429.

// In this example: rate.NewLimiter(10, 20)

// allows up to 10 requests per second with an initial burst of up to
// 20 requests. After the burst capacity is exhausted, new requests are
// permitted only as tokens are replenished.

// Compared to a fixed window counter, the token bucket algorithm
// provides smoother rate limiting and avoids large request bursts at
// window boundaries.

// func RateLimiter(next http.Handler) http.Handler {
// 	rl := rate.NewLimiter(10, 20)

// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		if !rl.Allow() {
// 			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
// 			return
// 		}
// 		next.ServeHTTP(w, r)
// 	})
// }
