package cache

import (
	"runtime"
	"sync"
	"time"

	. "github.com/AlexandreChamard/go-generic/algorithm"
	priorityqueue "github.com/AlexandreChamard/go-generic/priorityQueue"
)

type Cache[Key comparable, Value any] interface {
	Set(Key, Value)
	Get(Key) (value Value, ok bool)

	Size() int
	DeleteExpiredItems()
	FlushKOldest(int)
	Flush()

	Stop()
}

type cacheOptions struct {
	cacheDuration time.Duration // No expiration on 0
	maxSize       int           // No size limit on 0
	cacheOffset   int           // 10% of the maxSize on <=0 ; Remove that amount of items on reaching the maxSize
	purgeTimer    time.Duration // No auto purge on 0
}

func NewCacheOptions() *cacheOptions {
	return &cacheOptions{
		cacheDuration: 30 * time.Second,
		purgeTimer:    time.Minute,
	}
}

func (this *cacheOptions) CacheDuration(cacheDuration time.Duration) *cacheOptions {
	this.cacheDuration = cacheDuration
	return this
}

func (this *cacheOptions) NoExpiration() *cacheOptions {
	return this.CacheDuration(0)
}

func (this *cacheOptions) MaxSize(maxSize int) *cacheOptions {
	this.maxSize = maxSize
	return this
}

func (this *cacheOptions) NoSizeLimit() *cacheOptions {
	return this.MaxSize(0)
}

func (this *cacheOptions) CacheOffset(cacheOffset int) *cacheOptions {
	this.cacheOffset = cacheOffset
	return this
}

func (this *cacheOptions) DefaultCacheOffset() *cacheOptions {
	return this.CacheOffset(0)
}

func (this *cacheOptions) PurgeTimer(purgeTimer time.Duration) *cacheOptions {
	this.purgeTimer = purgeTimer
	return this
}

func (this *cacheOptions) NoPurge() *cacheOptions {
	return this.PurgeTimer(0)
}

// options must be define - use the NewCacheOptions to get the default values
func NewCache[Key comparable, Value any](options *cacheOptions) Cache[Key, Value] {
	return newCache[Key, Value](options)
}

type cache[Key comparable, Value any] struct {
	data          map[Key]mapData[Value]
	cacheDuration time.Duration
	maxSize       int
	cacheOffset   int

	runningChan chan bool
	mutex       sync.RWMutex
}

type mapData[Value any] struct {
	value    Value
	deleteAt int64
}

type cacheWrapper[Key comparable, Value any] struct {
	*cache[Key, Value]
}

func newCache[Key comparable, Value any](options *cacheOptions) *cacheWrapper[Key, Value] {
	if options == nil {
		return nil
	}
	cache := &cache[Key, Value]{
		data:          map[Key]mapData[Value]{},
		cacheDuration: options.cacheDuration,
		maxSize:       options.maxSize,
		cacheOffset:   Min(options.cacheOffset, options.maxSize),
	}

	if cache.cacheOffset <= 0 {
		cache.cacheOffset = cache.maxSize / 10
	}

	cacheW := &cacheWrapper[Key, Value]{cache}

	if options.purgeTimer > 0 {
		cache.runningChan = make(chan bool)
		go cache.runJanitor(options.purgeTimer)
		runtime.SetFinalizer(cacheW, func(cache *cacheWrapper[Key, Value]) {
			cache.Stop()
		})
	}
	return cacheW
}

func (this *cache[Key, Value]) Set(key Key, value Value) {
	this.set(key, value, this.cacheDuration, false)
}

func (this *cache[Key, Value]) SetWithExpiration(key Key, value Value, expiration time.Duration) {
	this.set(key, value, expiration, true)
}

func (this *cache[Key, Value]) set(key Key, value Value, expiration time.Duration, forceExpiration bool) {
	this.mutex.Lock()
	if data, ok := this.data[key]; ok && !forceExpiration {
		// If the key is aready defined, just reset its deletion time
		this.resetExpiration(key, mapData[Value]{
			value:    value,
			deleteAt: data.deleteAt,
		}, this.cacheDuration)
	} else {
		// Delete oldest elements if the cache size is exceeded
		if this.maxSize > 0 && this.Size() >= this.maxSize {
			this.flushKOldest(this.cacheOffset + 1) // remove the k oldest items in the cache (10% by default)
		}

		// If not, insert it in the map and the pqueue
		this.data[key] = mapData[Value]{
			value:    value,
			deleteAt: computeExpirationTimestamp(expiration),
		}
	}
	this.mutex.Unlock()
}

func (this *cache[Key, Value]) Get(key Key) (Value, bool) {
	this.mutex.RLock()
	data, ok := this.data[key]
	this.mutex.RUnlock()

	if ok && data.deleteAt != 0 && data.deleteAt <= time.Now().UnixNano() {
		return struct{ v Value }{}.v, false
	}

	return data.value, ok
}

func (this *cache[Key, Value]) resetExpiration(key Key, data mapData[Value], expiration time.Duration) {
	deleteAt := computeExpirationTimestamp(expiration)

	if deleteAt > data.deleteAt {
		this.data[key] = mapData[Value]{
			value:    data.value,
			deleteAt: deleteAt,
		}
	}
}

func (this *cache[Key, Value]) DeleteExpiredItems() {
	this.mutex.Lock()

	now := time.Now().UnixNano()
	for key, data := range this.data {
		if data.deleteAt != 0 && data.deleteAt < now {
			delete(this.data, key)
		}
	}
	this.mutex.Unlock()
}

func (this *cache[Key, Value]) Size() int {
	return len(this.data)
}

type keyDeleteTuple[Key any] struct {
	key      Key
	deleteAt int64
}

func (this *cache[Key, Value]) FlushKOldest(n int) {
	if n <= 0 {
		return
	}

	this.mutex.Lock()
	this.flushKOldest(n)
	this.mutex.Unlock()
}

func (this *cache[Key, Value]) flushKOldest(n int) {
	if n <= 0 {
		return
	}
	if n >= this.Size() {
		this.flush()
		return
	}

	queue := priorityqueue.NewPriorityQueue(func(a, b keyDeleteTuple[Key]) bool {
		return a.deleteAt < b.deleteAt
	})

	for key, value := range this.data {
		queue.Push(keyDeleteTuple[Key]{
			key:      key,
			deleteAt: value.deleteAt,
		})
		if queue.Size() > n {
			queue.Pop()
		}
	}

	for !queue.Empty() {
		delete(this.data, queue.Front().key)
		queue.Pop()
	}
}

func (this *cache[Key, Value]) Flush() {
	this.mutex.Lock()
	this.flush()
	this.mutex.Unlock()
}

func (this *cache[Key, Value]) flush() {
	this.data = make(map[Key]mapData[Value])
}

func (this *cache[Key, Value]) Stop() {
	this.mutex.Lock()
	this.flush()
	this.mutex.Unlock()

	close(this.runningChan)
}

func computeExpirationTimestamp(expiration time.Duration) int64 {
	if expiration > 0 {
		return time.Now().Add(expiration).UnixNano()
	}
	return 0
}

// should ONLY be called by the cache builder (see NewCache function)
func (this *cache[Key, Value]) runJanitor(purgeTimer time.Duration) {
	for {
		select {
		case <-this.runningChan:
			return
		case <-time.After(purgeTimer):
			this.DeleteExpiredItems()
		}
	}
}
