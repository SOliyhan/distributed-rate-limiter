package middleware

import (
	"net"
	"net/http"

	"github.com/SOliyhan/distributed-rate-limiter/internal/limiter"
)

func RateLimit(store *limiter.BucketStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				http.Error(w, "invalid IP", http.StatusInternalServerError)
				return
			}

			bucket := store.Get(ip)
			if !bucket.Allow() {
				http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
