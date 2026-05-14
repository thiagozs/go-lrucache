package lrucache

import (
	"fmt"
	"sync"
)

// Node represents a doubly linked list node used in the LRU cache.
type Node struct {
	key, value int
	prev, next *Node
}

// LRUCache stores integer key-value pairs using an LRU eviction policy.
type LRUCache struct {
	capacity int
	cache    map[int]*Node
	head     *Node
	tail     *Node
	mu       sync.RWMutex
}

const defaultCapacity = 128

type cacheConfig struct {
	capacity int
}

// Option configures LRUCache construction.
type Option func(*cacheConfig)

// WithCapacity sets a custom cache capacity.
// Negative values are normalized to zero.
func WithCapacity(capacity int) Option {
	return func(cfg *cacheConfig) {
		cfg.capacity = capacity
	}
}

// NewLRUCache creates a new cache using options.
//
// By default it uses an automatic capacity value.
// Use WithCapacity to override it explicitly.
// Negative capacities are normalized to zero.
func NewLRUCache(options ...Option) *LRUCache {
	cfg := cacheConfig{capacity: defaultCapacity}
	for _, opt := range options {
		if opt != nil {
			opt(&cfg)
		}
	}

	if cfg.capacity < 0 {
		cfg.capacity = 0
	}

	head := &Node{}
	tail := &Node{}

	head.next = tail
	tail.prev = head

	return &LRUCache{
		capacity: cfg.capacity,
		cache:    make(map[int]*Node, cfg.capacity),
		head:     head,
		tail:     tail,
	}
}

// Helpers do not lock by themselves; callers handle synchronization.
func (l *LRUCache) add(node *Node) {
	node.prev = l.head
	node.next = l.head.next

	l.head.next.prev = node
	l.head.next = node
}

func (l *LRUCache) remove(node *Node) {
	prev := node.prev
	next := node.next

	prev.next = next
	next.prev = prev
}

func (l *LRUCache) moveToFront(node *Node) {
	l.remove(node)
	l.add(node)
}

func (l *LRUCache) removeLRU() *Node {
	node := l.tail.prev
	l.remove(node)
	return node
}

// Get returns the value for a key and promotes it to most recently used.
// It returns -1 when the key is not present.
// Optimized: Uses RLock for initial lookup, then Lock only if found to minimize contention.
func (l *LRUCache) Get(key int) int {
	// First pass: RLock (read-only, allows concurrent readers)
	l.mu.RLock()
	node, ok := l.cache[key]
	l.mu.RUnlock()

	if !ok {
		return -1
	}

	// Second pass: Lock (write, because moveToFront modifies list)
	// Double-check pattern: re-verify node is still in cache
	l.mu.Lock()
	defer l.mu.Unlock()

	// Revalidate: node might have been evicted during the RUnlock window
	if currentNode, exists := l.cache[key]; exists &&
		currentNode == node {
		l.moveToFront(node)
		return node.value
	}

	return -1
}

// Put inserts or updates a value and evicts the least recently used item when needed.
func (l *LRUCache) Put(key int, value int) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if node, ok := l.cache[key]; ok {
		node.value = value
		l.moveToFront(node)
		return
	}

	node := &Node{key: key, value: value}
	l.cache[key] = node
	l.add(node)

	if len(l.cache) > l.capacity {
		lru := l.removeLRU()
		delete(l.cache, lru.key)
	}
}

// Debug prints the cache contents from most recently used to least recently used.
func (l *LRUCache) Debug() {
	l.mu.RLock()
	defer l.mu.RUnlock()

	fmt.Print("Cache [MRU → LRU]: ")
	curr := l.head.next
	for curr != nil && curr != l.tail {
		fmt.Printf("(%d:%d) ", curr.key, curr.value)
		curr = curr.next
	}
	fmt.Println()
}
