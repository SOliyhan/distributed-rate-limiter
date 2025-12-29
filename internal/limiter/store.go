package limiter

import "sync"

type BucketStore struct {
	buckets  map[string]*TokenBucket
	capacity int
	rate     float64
	mu       sync.Mutex
}

func NewBucketStore(capacity int, rate float64) *BucketStore {
	return &BucketStore{
		buckets:  make(map[string]*TokenBucket),
		capacity: capacity,
		rate:     rate,
	}
}

func (bs *BucketStore) Get(key string) *TokenBucket {
	bs.mu.Lock()
	defer bs.mu.Unlock()

	if bucket, ok := bs.buckets[key]; ok {
		return bucket
	}

	bucket := NewTokenBucket(bs.capacity, bs.rate)
	bs.buckets[key] = bucket
	return bucket
}
