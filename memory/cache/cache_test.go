package cache

import (
	"testing"
	"time"
)

func dumpCacheData[Key comparable, Value any](t *testing.T, cache *cache[Key, Value]) {
	t.Log(cache.data)
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
	dumpCacheData(t, cache.cache)

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
	dumpCacheData(t, cache.cache)
	time.Sleep(time.Second/2 + 10*time.Millisecond)
	t.Log("expected [3 0]")
	dumpCacheData(t, cache.cache)
	time.Sleep(time.Second)
	t.Log("expected []")
	dumpCacheData(t, cache.cache)
}

func TestCacheSetResetTheTimer(t *testing.T) {
	cache := newCache[string, int](NewCacheOptions().
		CacheDuration(100 * time.Millisecond).
		PurgeTimer(100 * time.Millisecond))

	cache.Set("0", 0)
	time.Sleep(50 * time.Millisecond)
	cache.mutex.RLock()
	if cache.Size() == 0 {
		t.Errorf("the cache should not be empty at this point")
	}
	cache.mutex.RUnlock()

	cache.Set("0", 0)
	time.Sleep(50 * time.Millisecond)
	cache.mutex.RLock()
	if cache.Size() == 0 {
		t.Errorf("the cache should not be empty at this point")
	}
	cache.mutex.RUnlock()

	cache.Set("0", 0)
	time.Sleep(50 * time.Millisecond)
	cache.mutex.RLock()
	if cache.Size() == 0 {
		t.Errorf("the cache should not be empty at this point")
	}
	cache.mutex.RUnlock()

	cache.Set("0", 0)
	time.Sleep(50 * time.Millisecond)
	cache.mutex.RLock()
	if cache.Size() == 0 {
		t.Errorf("the cache should not be empty at this point")
	}
	cache.mutex.RUnlock()

	dumpCacheData(t, cache.cache)

	// add a little of time to be sure
	time.Sleep(200 * time.Millisecond)

	cache.mutex.RLock()
	if cache.Size() != 0 {
		t.Errorf("the cache should be empty at this point")
	}
	cache.mutex.RUnlock()
}

func TestCacheSetWithDeadlineResetTheTimerWhenNeeded(t *testing.T) {
	cache := newCache[string, int](NewCacheOptions().CacheDuration(100 * time.Millisecond))

	/*
		set 0
		wait 50ms
		delete expired items
		check if not delete
		wait 50ms
		delete expired items
		check if delete

		set 0
		set 0 10ms
		wait 50ms
		delete expired items
		check if delete

		set 0
		set 0 200ms
		wait 100ms
		delete expired items
		check if not delete
	*/

	cache.Set("0", 0)
	time.Sleep(50 * time.Millisecond)
	cache.DeleteExpiredItems()
	if cache.Size() == 0 {
		t.Errorf("the cache should not be empty at this point")
	}

	time.Sleep(50 * time.Millisecond)
	cache.DeleteExpiredItems()
	if cache.Size() != 0 {
		t.Errorf("the cache should be empty at this point")
	}

	cache.Set("0", 0)
	cache.SetWithExpiration("0", 0, 10*time.Millisecond)
	time.Sleep(50 * time.Millisecond)
	cache.DeleteExpiredItems()
	if cache.Size() != 0 {
		t.Errorf("the cache should be empty at this point")
	}

	cache.Set("0", 0)
	cache.SetWithExpiration("0", 0, 200*time.Millisecond)
	time.Sleep(100 * time.Millisecond)
	cache.DeleteExpiredItems()
	if cache.Size() == 0 {
		t.Errorf("the cache should not be empty at this point")
	}
}

func TestCacheStop(t *testing.T) {
	cache := newCache[string, int](NewCacheOptions().CacheDuration(100 * time.Millisecond))

	cache.Set("0", 0)

	cache.Stop()
	time.Sleep(10 * time.Millisecond)
	dumpCacheData(t, cache.cache)
}

func TestCacheSizeLimit(t *testing.T) {
	cache := newCache[string, int](NewCacheOptions().NoExpiration().MaxSize(5))

	cache.Set("0", 0)
	cache.Set("1", 1)
	cache.Set("2", 2)
	cache.Set("3", 3)
	cache.Set("4", 4)
	if cache.Size() != 5 {
		t.Fatalf("cache should have %d elements got %d", 5, cache.Size())
	}
	cache.Set("5", 5)
	if cache.Size() != 5 {
		t.Fatalf("cache should have %d elements got %d", 5, cache.Size())
	}
}

func TestCacheBigExpiration(t *testing.T) {
	cache := newCache[string, int](NewCacheOptions().CacheDuration(time.Second / 2).PurgeTimer(1000 * time.Millisecond))

	cache.SetWithExpiration("0", 0, 2500*time.Millisecond)
	time.Sleep(time.Second)
	cache.mutex.RLock()
	if cache.Size() != 1 {
		t.Fatalf("cache should have %d elements got %d", 1, cache.Size())
	}
	cache.mutex.RUnlock()

	cache.Get("0")
	time.Sleep(time.Second)
	cache.mutex.RLock()
	if cache.Size() != 1 {
		t.Fatalf("cache should have %d elements got %d", 1, cache.Size())
	}
	cache.mutex.RUnlock()

	time.Sleep(2 * time.Second)
	cache.mutex.RLock()
	if cache.Size() != 0 {
		t.Fatalf("cache should have %d elements got %d", 0, cache.Size())
	}
	cache.mutex.RUnlock()
}

var k int

func benchmarkCache(b *testing.B, n int, cacheSize int) {
	cache := NewCache[int, int](NewCacheOptions().CacheDuration(1 * time.Minute).MaxSize(cacheSize))

	lastSize := 0
	for round := 0; round < 100; round++ {
		for i := 0; i < n; i++ {
			cache.Set(i, i)
			if lastSize > cache.Size() {
				// b.Logf("last size %d now %d", lastSize, cache.Size())
			}
			lastSize = cache.Size()
		}
		for i := 0; i < n; i++ {
			k, _ = cache.Get(i)
		}
	}
}

func BenchmarkCache_100_1(b *testing.B)    { benchmarkCache(b, 100, 100/2) }
func BenchmarkCache_100_2(b *testing.B)    { benchmarkCache(b, 100, 100) }
func BenchmarkCache_200_1(b *testing.B)    { benchmarkCache(b, 200, 200/2) }
func BenchmarkCache_200_2(b *testing.B)    { benchmarkCache(b, 200, 200) }
func BenchmarkCache_1000_1(b *testing.B)   { benchmarkCache(b, 1000, 1000/2) }
func BenchmarkCache_1000_2(b *testing.B)   { benchmarkCache(b, 1000, 1000) }
func BenchmarkCache_10000_1(b *testing.B)  { benchmarkCache(b, 10000, 10000/2) }
func BenchmarkCache_10000_2(b *testing.B)  { benchmarkCache(b, 10000, 10000) }
func BenchmarkCache_100000_1(b *testing.B) { benchmarkCache(b, 100000, 100000/2) }
func BenchmarkCache_100000_2(b *testing.B) { benchmarkCache(b, 100000, 100000) }
