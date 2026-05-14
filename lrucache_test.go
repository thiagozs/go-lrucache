package lrucache

import (
	"sync"
	"testing"
)

func TestLRUCachePutGetAndEvictionOrder(t *testing.T) {
	cache := NewLRUCache(WithCapacity(2))

	cache.Put(1, 10)
	cache.Put(2, 20)

	if got := cache.Get(1); got != 10 {
		t.Fatalf("expected key 1 to return 10, got %d", got)
	}

	// Accessing key 1 makes key 2 the LRU. Next insertion must evict key 2.
	cache.Put(3, 30)

	if got := cache.Get(2); got != -1 {
		t.Fatalf("expected key 2 to be evicted, got %d", got)
	}
	if got := cache.Get(1); got != 10 {
		t.Fatalf("expected key 1 to stay in cache with value 10, got %d", got)
	}
	if got := cache.Get(3); got != 30 {
		t.Fatalf("expected key 3 to return 30, got %d", got)
	}
}

func TestLRUCacheUpdateExistingKey(t *testing.T) {
	cache := NewLRUCache(WithCapacity(2))

	cache.Put(1, 10)
	cache.Put(2, 20)
	cache.Put(1, 15) // update + move to MRU
	cache.Put(3, 30) // should evict key 2

	if got := cache.Get(1); got != 15 {
		t.Fatalf("expected updated value 15 for key 1, got %d", got)
	}
	if got := cache.Get(2); got != -1 {
		t.Fatalf("expected key 2 to be evicted, got %d", got)
	}
	if got := cache.Get(3); got != 30 {
		t.Fatalf("expected key 3 to return 30, got %d", got)
	}
}

func TestLRUCacheZeroCapacity(t *testing.T) {
	cache := NewLRUCache(WithCapacity(0))

	cache.Put(1, 10)

	if got := cache.Get(1); got != -1 {
		t.Fatalf("expected key 1 to be unavailable for zero-capacity cache, got %d", got)
	}
}

func TestLRUCacheNegativeCapacityBehavesAsZero(t *testing.T) {
	cache := NewLRUCache(WithCapacity(-5))

	cache.Put(1, 10)

	if got := cache.Get(1); got != -1 {
		t.Fatalf("expected key 1 to be unavailable for negative-capacity cache, got %d", got)
	}
}

func TestLRUCacheConcurrentAccess(t *testing.T) {
	cache := NewLRUCache(WithCapacity(5))

	const workers = 8
	const iterations = 200

	var wg sync.WaitGroup
	for w := range workers {
		w := w
		wg.Go(func() {
			for i := range iterations {
				key := (w + i) % 12
				cache.Put(key, key*10)
				_ = cache.Get((key + 1) % 12)
			}
		})
	}

	wg.Wait()

	cache.mu.RLock()
	defer cache.mu.RUnlock()

	if len(cache.cache) > cache.capacity {
		t.Fatalf("cache size exceeded capacity: size=%d capacity=%d", len(cache.cache), cache.capacity)
	}

	count := 0
	for node := cache.head.next; node != nil && node != cache.tail; node = node.next {
		count++
		if count > cache.capacity {
			t.Fatalf("linked list contains more nodes than capacity: count=%d capacity=%d", count, cache.capacity)
		}
		if _, ok := cache.cache[node.key]; !ok {
			t.Fatalf("node key %d found in list but missing in map", node.key)
		}
	}

	if count != len(cache.cache) {
		t.Fatalf("list/map size mismatch: list=%d map=%d", count, len(cache.cache))
	}
}

func TestLRUCacheDefaultCapacityIsApplied(t *testing.T) {
	cache := NewLRUCache()

	if cache.capacity != defaultCapacity {
		t.Fatalf("expected default capacity %d, got %d", defaultCapacity, cache.capacity)
	}
}
