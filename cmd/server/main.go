package main

import (
	"log"
	"net/http"

	"github.com/SOliyhan/distributed-rate-limiter/internal/limiter"
	"github.com/SOliyhan/distributed-rate-limiter/internal/middleware"
)

func main() {
	store := limiter.NewBucketStore(10, 1)

	mux := http.NewServeMux()
	
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Distributed Rate Limiter is running"))
	})
	
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	
	handler := middleware.RateLimit(store)(mux)

	log.Println("server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
