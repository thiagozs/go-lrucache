package lrucache

import (
	"testing"
)

// Benchmark para operação Get em cache hit
func BenchmarkGet_CacheHit(b *testing.B) {
	cache := NewLRUCache(WithCapacity(1000))
	cache.Put(1, 100)

	for b.Loop() {
		cache.Get(1)
	}
}

// Benchmark para operação Get em cache miss
func BenchmarkGet_CacheMiss(b *testing.B) {
	cache := NewLRUCache(WithCapacity(1000))
	cache.Put(1, 100)

	for b.Loop() {
		cache.Get(2)
	}
}

// Benchmark para operação Put com novo item
func BenchmarkPut_NewItem(b *testing.B) {
	cache := NewLRUCache(WithCapacity(10000))

	for i := 0; b.Loop(); i++ {
		cache.Put(i, i*10)
	}
}

// Benchmark para operação Put atualizando item existente
func BenchmarkPut_ExistingItem(b *testing.B) {
	cache := NewLRUCache(WithCapacity(1000))
	cache.Put(1, 100)

	for i := 0; b.Loop(); i++ {
		cache.Put(1, i*10)
	}
}

// Benchmark para operação Put com evicção
func BenchmarkPut_WithEviction(b *testing.B) {
	cache := NewLRUCache(WithCapacity(100))

	// Preenche o cache
	for i := range 100 {
		cache.Put(i, i*10)
	}

	for i := 0; b.Loop(); i++ {
		cache.Put((i%1000)+100, (i%1000)*10)
	}
}

// Benchmark para padrão Get/Put alternado
func BenchmarkMixed_GetPut(b *testing.B) {
	cache := NewLRUCache(WithCapacity(500))

	for i := 0; b.Loop(); i++ {
		key := i % 1000
		cache.Put(key, key*10)
		cache.Get((key + 1) % 1000)
	}
}

// Benchmark para padrão típico de cache (80% reads, 20% writes)
func BenchmarkTypical_80Read_20Write(b *testing.B) {
	cache := NewLRUCache(WithCapacity(1000))

	// Pré-populate
	for i := range 100 {
		cache.Put(i, i*10)
	}

	for i := 0; b.Loop(); i++ {
		if i%5 < 4 {
			cache.Get(i % 100)
		} else {
			cache.Put(i%1000, (i%1000)*10)
		}
	}
}

// Benchmark com concorrência
func BenchmarkConcurrent_16Goroutines(b *testing.B) {
	cache := NewLRUCache(WithCapacity(1000))

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := i % 500
			cache.Put(key, key*10)
			cache.Get((key + 1) % 500)
			i++
		}
	})
}

// Benchmark moveToFront (operação crítica)
func BenchmarkMoveToFront(b *testing.B) {
	cache := NewLRUCache(WithCapacity(10000))

	// Popula com alguns items
	for i := range 100 {
		cache.Put(i, i*10)
	}

	// Get promove item para frente (moveToFront)

	for i := 0; b.Loop(); i++ {
		cache.Get(i % 100)
	}
}

// Benchmark Debug output (não é crítico de performance)
func BenchmarkDebug(b *testing.B) {
	cache := NewLRUCache(WithCapacity(100))
	for i := range 50 {
		cache.Put(i, i*10)
	}

	for b.Loop() {
		cache.Debug()
	}
}
