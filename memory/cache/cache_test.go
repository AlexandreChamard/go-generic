package cache

import (
	"fmt"
	"testing"
	"time"
)

// redefine newCache to check internal data on running
func newCache[Key comparable, Value any](options *cacheOptions) *cache[Key, Value] {
	if options == nil {
		return nil
	}
	cache := &cache[Key, Value]{
		data:          map[Key]mapData[Value]{},
		cacheDuration: options.cacheDuration,
		maxSize:       options.maxSize,

		running:  true,
		pingChan: make(chan bool),
	}

	cache.pqueue = NewPriorityQueue(
		func(a pqueueData[Key], b pqueueData[Key]) bool {
			return a.deleteAt.Before(b.deleteAt)
		},
		func(a pqueueData[Key], b pqueueData[Key]) {
			// used to track the position of all key in the pqueue
			aData := cache.data[a.key]
			bData := cache.data[b.key]
			cache.data[a.key] = mapData[Value]{
				value: aData.value,
				pos:   bData.pos,
			}
			cache.data[b.key] = mapData[Value]{
				value: bData.value,
				pos:   aData.pos,
			}
		})

	go func() {
		for cache.running {
			cache.mutex.RLock()
			if cache.pqueue.Empty() || cache.pqueue.Front().deleteAt == (time.Time{}) {
				cache.mutex.RUnlock()
				select {
				case <-cache.pingChan:
					// Nothing to do
				}
			} else {
				nextDeletion := time.Until(cache.pqueue.Front().deleteAt)
				cache.mutex.RUnlock()
				fmt.Println("next deletion in", nextDeletion.Milliseconds(), "ms")
				select {
				case <-cache.pingChan:
					// Nothing to do

				case <-time.After(nextDeletion):
					fmt.Println("call for deletion")
					cache.deleteExpiredItems()
				}
			}
		}
		fmt.Println("stop the gofunc")
	}()
	return cache
}

func dumpCacheData[Key comparable, Value any](t *testing.T, cache *cache[Key, Value]) {
	t.Log(cache.data)
	t.Log(*cache.pqueue)
}

func TestTrueCache(t *testing.T) {
	cache := NewCache[string, int](NewCacheOptions().CacheDuration(time.Second))

	t.Log("start test")
	t.Log("insert 0")
	cache.Set("0", 0)
	time.Sleep(100 * time.Millisecond)
	t.Log("insert 1")
	cache.Set("1", 1)
	time.Sleep(100 * time.Millisecond)
	t.Log("insert 2")
	cache.Set("2", 2)
	time.Sleep(100 * time.Millisecond)
	t.Log("insert 3")
	cache.Set("3", 3)
	time.Sleep(100 * time.Millisecond)
	t.Log("expected [0 1 2 3]")

	t.Log("insert 1")
	cache.Set("1", 1)
	t.Log("insert 2")
	cache.Set("2", 2)
	time.Sleep(time.Second / 2)
	t.Log("insert 3")
	cache.Set("3", 3)
	t.Log("insert 0")
	cache.Set("0", 0)
	t.Log("expected [1 2 3 0]")
	time.Sleep(time.Second/2 + 10*time.Millisecond)
	t.Log("expected [3 0]")
	time.Sleep(time.Second)
	t.Log("expected []")
}

func TestCache(t *testing.T) {
	cache := newCache[string, int](NewCacheOptions().CacheDuration(time.Second))

	t.Log("start test")
	t.Log("insert 0")
	cache.Set("0", 0)
	time.Sleep(100 * time.Millisecond)
	t.Log("insert 1")
	cache.Set("1", 1)
	time.Sleep(100 * time.Millisecond)
	t.Log("insert 2")
	cache.Set("2", 2)
	time.Sleep(100 * time.Millisecond)
	t.Log("insert 3")
	cache.Set("3", 3)
	time.Sleep(100 * time.Millisecond)
	t.Log("expected [0 1 2 3]")
	dumpCacheData(t, cache)

	t.Log("insert 1")
	cache.Set("1", 1)
	t.Log("insert 2")
	cache.Set("2", 2)
	time.Sleep(time.Second / 2)
	t.Log("insert 3")
	cache.Set("3", 3)
	t.Log("insert 0")
	cache.Set("0", 0)
	t.Log("expected [1 2 3 0]")
	dumpCacheData(t, cache)
	time.Sleep(time.Second/2 + 10*time.Millisecond)
	t.Log("expected [3 0]")
	dumpCacheData(t, cache)
	time.Sleep(time.Second)
	t.Log("expected []")
	dumpCacheData(t, cache)
}
func TestCacheSetResetTheTimer(t *testing.T) {
	cache := newCache[string, int](NewCacheOptions().CacheDuration(100 * time.Millisecond))

	cache.Set("0", 0)
	time.Sleep(50 * time.Millisecond)
	if cache.pqueue.Empty() {
		t.Errorf("the cache should not be empty at this point")
	}

	cache.Set("0", 0)
	time.Sleep(50 * time.Millisecond)
	if cache.pqueue.Empty() {
		t.Errorf("the cache should not be empty at this point")
	}

	cache.Set("0", 0)
	time.Sleep(50 * time.Millisecond)
	if cache.pqueue.Empty() {
		t.Errorf("the cache should not be empty at this point")
	}

	cache.Set("0", 0)
	time.Sleep(50 * time.Millisecond)
	if cache.pqueue.Empty() {
		t.Errorf("the cache should not be empty at this point")
	}

	dumpCacheData(t, cache)

	// add a little of time to be sure
	time.Sleep(52 * time.Millisecond)

	if !cache.pqueue.Empty() {
		t.Errorf("the cache should be empty at this point")
	}
}

func TestCacheGetResetTheTimer(t *testing.T) {
	cache := newCache[string, int](NewCacheOptions().CacheDuration(100 * time.Millisecond))

	cache.Set("0", 0)
	time.Sleep(50 * time.Millisecond)
	if cache.pqueue.Empty() {
		t.Errorf("the cache should not be empty at this point")
	}

	cache.Get("0")
	time.Sleep(50 * time.Millisecond)
	if cache.pqueue.Empty() {
		t.Errorf("the cache should not be empty at this point")
	}

	cache.Get("0")
	time.Sleep(50 * time.Millisecond)
	if cache.pqueue.Empty() {
		t.Errorf("the cache should not be empty at this point")
	}

	cache.Get("0")
	time.Sleep(50 * time.Millisecond)
	if cache.pqueue.Empty() {
		t.Errorf("the cache should not be empty at this point")
	}

	dumpCacheData(t, cache)

	// add a little of time to be sure
	time.Sleep(52 * time.Millisecond)
	if !cache.pqueue.Empty() {
		t.Errorf("the cache should be empty at this point")
	}
}

func TestCacheStop(t *testing.T) {
	cache := newCache[string, int](NewCacheOptions().CacheDuration(100 * time.Millisecond))

	cache.Set("0", 0)

	cache.Stop()
	time.Sleep(10 * time.Millisecond)
	dumpCacheData(t, cache)
}

func TestCacheSizeLimit(t *testing.T) {
	cache := newCache[string, int](NewCacheOptions().NoExpiration().MaxSize(5))

	cache.Set("0", 0)
	cache.Set("1", 1)
	cache.Set("2", 2)
	cache.Set("3", 3)
	cache.Set("4", 4)
	if len(cache.data) != 5 {
		t.Fatalf("cache should have %d elements got %d", 5, len(cache.data))
	}
	cache.Set("5", 5)
	if len(cache.data) != 5 {
		t.Fatalf("cache should have %d elements got %d", 5, len(cache.data))
	}
}

func TestCacheBigExpiration(t *testing.T) {
	cache := newCache[string, int](NewCacheOptions().CacheDuration(time.Second / 2))

	cache.SetWithExpiration("0", 0, 2500*time.Millisecond)
	time.Sleep(time.Second)
	if len(cache.data) != 1 {
		t.Fatalf("cache should have %d elements got %d", 1, len(cache.data))
	}
	cache.Get("0")
	time.Sleep(time.Second)
	if len(cache.data) != 1 {
		t.Fatalf("cache should have %d elements got %d", 1, len(cache.data))
	}

	time.Sleep(time.Second)
	if len(cache.data) != 0 {
		t.Fatalf("cache should have %d elements got %d", 0, len(cache.data))
	}
}
