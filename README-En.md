# LRU Cache in Go

### Overview

This project implements a thread-safe **LRU (Least Recently Used)** cache in Go.

It combines:
- a map for fast key lookup
- a doubly linked list to track usage order (MRU -> LRU)
- `sync.RWMutex` for concurrent access safety

When capacity is exceeded, the least recently used item is automatically evicted.

### Usage
Prerequisite: Go installed.

```go
package main

import (
	"fmt"

	lrucache "github.com/thiagozs/go-lrucache"
)

func main() {
	cache := lrucache.NewLRUCache(lrucache.WithCapacity(2))
	cache.Put(1, 10)
	cache.Put(2, 20)
	cache.Put(3, 30)

	fmt.Println(cache.Get(1)) // -1 (evicted)
	fmt.Println(cache.Get(2)) // 20
	fmt.Println(cache.Get(3)) // 30
}
```

### Run tests
```bash
go test ./...
```

### Project structure
- [lrucache.go](lrucache.go): library implementation.
- [lrucache_test.go](lrucache_test.go): package unit tests.
- [examples/main.go](examples/main.go): executable example consuming the library.

### Main API
- `NewLRUCache(options ...Option) *LRUCache`: creates a cache with automatic default capacity.
- `WithCapacity(capacity int) Option`: sets capacity explicitly when needed.
- `Get(key int) int`: returns the value and promotes key to MRU; returns `-1` if missing.
- `Put(key int, value int)`: inserts/updates a value; evicts when over capacity.
- `Debug()`: prints current state in MRU -> LRU order.

### Internal design
- `head` and `tail` are sentinel nodes.
- Inserts/promotions happen right after `head` (MRU side).
- The LRU node is right before `tail`.

Write flow (`Put`):
1. If key exists, update value and move to front.
2. If missing, add to front.
3. If `len(cache) > capacity`, remove LRU and delete from map.

Read flow (`Get`):
1. If key exists, move to front and return value.
2. If missing, return `-1`.

### Complexity
- `Get`: O(1)
- `Put`: O(1)
- Space: O(capacity)

### Concurrency
Public methods lock around operations:
- `Get` and `Put`: `mu.Lock()` (both mutate LRU order)
- `Debug`: `mu.RLock()`

Internal helpers (`add`, `remove`, `moveToFront`, `removeLRU`) do not lock by themselves; synchronization is handled by caller public methods.

### `examples` folder
- The repository includes [examples/main.go](examples/main.go) as a runnable example.

### Current limitations
- No explicit validation for `capacity <= 0`.
