package middleware

import (
	CryptoSHA256 "crypto/sha256"
	EncodingHex "encoding/hex"
	Net "net"
	Http "net/http"
	"strconv"
	Strings "strings"
	Sync "sync"
	Time "time"
)

type rateLimitEntry struct {
	windowStart Time.Time
	count       int
}

type RateLimiter struct {
	mu      Sync.Mutex
	entries map[string]rateLimitEntry
	limit   int
	window  Time.Duration
}

func NewRateLimiter(limit int, window Time.Duration) *RateLimiter {
	if limit < 1 {
		limit = 60
	}
	if window <= 0 {
		window = Time.Minute
	}
	return &RateLimiter{entries: make(map[string]rateLimitEntry), limit: limit, window: window}
}

func (limiter *RateLimiter) Middleware(next Http.Handler) Http.Handler {
	return Http.HandlerFunc(func(response Http.ResponseWriter, request *Http.Request) {
		key := clientKey(request)
		now := Time.Now()
		limiter.mu.Lock()
		entry := limiter.entries[key]
		if entry.windowStart.IsZero() || now.Sub(entry.windowStart) >= limiter.window {
			entry = rateLimitEntry{windowStart: now}
		}
		entry.count++
		allowed := entry.count <= limiter.limit
		limiter.entries[key] = entry
		limiter.mu.Unlock()
		if !allowed {
			response.Header().Set("Retry-After", strconv.Itoa(int(limiter.window.Seconds())))
			Http.Error(response, "Too many requests", Http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(response, request)
	})
}

func clientKey(request *Http.Request) string {
	host, _, err := Net.SplitHostPort(request.RemoteAddr)
	if err != nil {
		host = request.RemoteAddr
	}
	token := request.Header.Get(Authorization)
	sum := CryptoSHA256.Sum256([]byte(token))
	return host + ":" + EncodingHex.EncodeToString(sum[:])
}

func SecurityHeaders(allowedOrigins []string) func(Http.Handler) Http.Handler {
	return func(next Http.Handler) Http.Handler {
		return Http.HandlerFunc(func(response Http.ResponseWriter, request *Http.Request) {
			response.Header().Set("X-Content-Type-Options", "nosniff")
			response.Header().Set("X-Frame-Options", "DENY")
			response.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
			response.Header().Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'")
			response.Header().Set("X-Robots-Tag", "noindex, nofollow, noarchive")
			if request.Method == Http.MethodGet && request.Header.Get(Authorization) == "" {
				response.Header().Set("Cache-Control", "public, max-age=30, stale-while-revalidate=60")
			} else {
				response.Header().Set("Cache-Control", "no-store")
			}
			origin := request.Header.Get("Origin")
			if origin != "" && containsOrigin(allowedOrigins, origin) {
				response.Header().Set("Access-Control-Allow-Origin", origin)
				response.Header().Set("Vary", "Origin")
			}
			if request.Method == Http.MethodOptions {
				if origin == "" || !containsOrigin(allowedOrigins, origin) {
					Http.Error(response, "Origin not allowed", Http.StatusForbidden)
					return
				}
				response.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				response.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Correlation-ID")
				response.WriteHeader(Http.StatusNoContent)
				return
			}
			next.ServeHTTP(response, request)
		})
	}
}

func containsOrigin(allowed []string, origin string) bool {
	for _, candidate := range allowed {
		if Strings.EqualFold(candidate, origin) {
			return true
		}
	}
	return false
}
