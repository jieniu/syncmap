// A thread safe map implementation for Golang
package syncmap

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// syncMap wraps built-in map by using RWMutex for concurrent safe.
type syncMap64 struct {
	items map[uint64]interface{}
	sync.RWMutex
}

// SyncMap keeps a slice of *syncMap with length of `shardCount`.
// Using a slice of syncMap instead of a large one is to avoid lock bottlenecks.
type SyncMap64 struct {
	shardCount uint8
	shards     []*syncMap64
}

// Create a new SyncMap with default shard count.
func New64() *SyncMap64 {
	return NewWithShard64(defaultShardCount)
}

// Create a new SyncMap with given shard count.
// NOTE: shard count must be power of 2, default shard count will be used otherwise.
func NewWithShard64(shardCount uint8) *SyncMap64 {
	if !isPowerOfTwo(shardCount) {
		shardCount = defaultShardCount
	}
	m := new(SyncMap64)
	m.shardCount = shardCount
	m.shards = make([]*syncMap64, m.shardCount)
	for i, _ := range m.shards {
		m.shards[i] = &syncMap64{items: make(map[uint64]interface{})}
	}
	return m
}

// Find the specific shard with the given key
func (m *SyncMap64) locate(key uint64) *syncMap64 {
	strkey := fmt.Sprintf("%d", key)
	return m.shards[bkdrHash(strkey)&uint32((m.shardCount-1))]
}

// Retrieves a value
func (m *SyncMap64) Get(key uint64) (value interface{}, ok bool) {
	shard := m.locate(key)
	shard.RLock()
	value, ok = shard.items[key]
	shard.RUnlock()
	return
}

// Sets value with the given key
func (m *SyncMap64) Set(key uint64, value interface{}) {
	shard := m.locate(key)
	shard.Lock()
	shard.items[key] = value
	shard.Unlock()
}

// Removes an item
func (m *SyncMap64) Delete(key uint64) {
	shard := m.locate(key)
	shard.Lock()
	delete(shard.items, key)
	shard.Unlock()
}

// Pop delete and return a random item in the cache
func (m *SyncMap64) Pop() (uint64, interface{}) {
	if m.Size() == 0 {
		panic("syncmap64: map is empty")
	}

	var (
		key   uint64
		value interface{}
		found = false
		n     = int(m.shardCount)
	)

	for !found {
		idx := rand.Intn(n)
		shard := m.shards[idx]
		shard.Lock()
		if len(shard.items) > 0 {
			found = true
			for key, value = range shard.items {
				break
			}
			delete(shard.items, key)
		}
		shard.Unlock()
	}

	return key, value
}

// Whether SyncMap has the given key
func (m *SyncMap64) Has(key uint64) bool {
	_, ok := m.Get(key)
	return ok
}

// Returns the number of items
func (m *SyncMap64) Size() int {
	size := 0
	for _, shard := range m.shards {
		shard.RLock()
		size += len(shard.items)
		shard.RUnlock()
	}
	return size
}

// Wipes all items from the map
func (m *SyncMap64) Flush() int {
	size := 0
	for _, shard := range m.shards {
		shard.Lock()
		size += len(shard.items)
		shard.items = make(map[uint64]interface{})
		shard.Unlock()
	}
	return size
}

// Returns a channel from which each key in the map can be read
func (m *SyncMap64) IterKeys() <-chan uint64 {
	ch := make(chan uint64)
	go func() {
		for _, shard := range m.shards {
			shard.RLock()
			for key, _ := range shard.items {
				ch <- key
			}
			shard.RUnlock()
		}
		close(ch)
	}()
	return ch
}

// Item is a pair of key and value
type Item64 struct {
	Key   uint64
	Value interface{}
}

// Return a channel from which each item (key:value pair) in the map can be read
func (m *SyncMap64) IterItems() <-chan Item64 {
	ch := make(chan Item64)
	go func() {
		for _, shard := range m.shards {
			shard.RLock()
			for key, value := range shard.items {
				ch <- Item64{key, value}
			}
			shard.RUnlock()
		}
		close(ch)
	}()
	return ch
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
