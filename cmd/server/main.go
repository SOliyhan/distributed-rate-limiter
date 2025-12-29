package main

import (
	"context"
	"log"
	"net/http"

	"github.com/SOliyhan/distributed-rate-limiter/internal/limiter"
	"github.com/SOliyhan/distributed-rate-limiter/internal/middleware"
	redisclient "github.com/SOliyhan/distributed-rate-limiter/internal/redis"
)

func main() {
	ctx := context.Background()

	redisClient := redisclient.NewClient()
	if err := redisclient.Ping(ctx, redisClient); err != nil {
		log.Fatal("redis not available")
	}

	redisLimiter := limiter.NewRedisLimiter(redisClient, 10, 1)
	fallback := limiter.NewBucketStore(10, 1)

	mux := http.NewServeMux()
	
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Distributed Rate Limiter is running"))
	})
	
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	handler := middleware.RateLimitRedis(redisLimiter, fallback)(mux)

	log.Println("server started on :8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", handler))
}
