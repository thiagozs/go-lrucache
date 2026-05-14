package main

import (
	"fmt"

	"github.com/thiagozs/go-lrucache"
)

func main() {
	cache := lrucache.NewLRUCache(lrucache.WithCapacity(2))

	cache.Put(1, 10)
	cache.Put(2, 20)

	fmt.Println(cache.Get(1)) // returns 10
	cache.Put(3, 30)          // evicts key 2
	fmt.Println(cache.Get(2))
	cache.Put(4, 40) // evicts key 3
	fmt.Println(cache.Get(3))
	fmt.Println(cache.Get(4))

	cache.Debug()

}
