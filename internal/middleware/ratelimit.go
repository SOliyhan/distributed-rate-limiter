package middleware

import (
	"net"
	"net/http"

	"github.com/SOliyhan/distributed-rate-limiter/internal/limiter"
)

func RateLimitRedis(
	redisLimiter *limiter.RedisLimiter,
	fallback *limiter.BucketStore,
) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, _, _ := net.SplitHostPort(r.RemoteAddr)

			allowed, err := redisLimiter.Allow(r.Context(), ip)
			if err != nil {
				// Redis down then fallback to memory
				if !fallback.Get(ip).Allow() {
					http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
					return
				}
				next.ServeHTTP(w, r)
				return
			}

			if !allowed {
				http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
