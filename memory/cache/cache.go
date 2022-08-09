package cache

import (
	"sync"
	"time"
)

type Cache[Key comparable, Value any] interface {
	Set(Key, Value)
	Get(Key) (value Value, ok bool)

	Clear()

	Stop()
}

type cacheOptions struct {
	cacheDuration time.Duration // No expiration on 0
	maxSize       int           // No size limit on 0
}

func NewCacheOptions() *cacheOptions {
	return &cacheOptions{
		cacheDuration: 30 * time.Second,
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

// options must be define - use the NewCacheOptions to get the default values
func NewCache[Key comparable, Value any](options *cacheOptions) Cache[Key, Value] {
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
				select {
				case <-cache.pingChan:
					// Nothing to do

				case <-time.After(nextDeletion):
					cache.deleteExpiredItems()
				}
			}
		}
	}()
	return cache
}

type cache[Key comparable, Value any] struct {
	data          map[Key]mapData[Value]
	pqueue        *priorityQueue[pqueueData[Key]]
	cacheDuration time.Duration
	maxSize       int

	running  bool
	pingChan chan bool
	mutex    sync.RWMutex
}

type mapData[Value any] struct {
	value Value
	pos   int // helper to find the key in the pqueue
}

type pqueueData[Key any] struct {
	key      Key
	deleteAt time.Time
}

func (this *cache[Key, Value]) Set(key Key, value Value) {
	this.set(key, value, this.cacheDuration)
}

func (this *cache[Key, Value]) SetWithExpiration(key Key, value Value, expiration time.Duration) {
	this.set(key, value, expiration)
}

func (this *cache[Key, Value]) set(key Key, value Value, expiration time.Duration) {
	if !this.running {
		return
	}

	this.mutex.Lock()

	if data, ok := this.data[key]; ok {
		// if the key is aready defined, just reset its deletion time
		this.resetExpiration(key, data, this.cacheDuration)
	} else {
		if this.maxSize > 0 && len(this.data) >= this.maxSize {
			this.deleteElements(len(this.data) - this.maxSize + 1)
		}
		// if not, insert it in the map and the pqueue
		this.data[key] = mapData[Value]{
			value: value,
			pos:   this.pqueue.Size(),
		}
		this.pqueue.Push(pqueueData[Key]{
			key:      key,
			deleteAt: computeExpirationTimestamp(expiration),
		})
	}
	this.mutex.Unlock()
	select {
	case this.pingChan <- true:
	default:
	}
}

func (this *cache[Key, Value]) Get(key Key) (Value, bool) {
	this.mutex.RLock()
	data, ok := this.data[key]
	this.mutex.RUnlock()

	// reset the key duration in the cache
	// can e done asynchronously to not wait for it
	go func() {
		this.mutex.Lock()
		this.resetExpiration(key, data, this.cacheDuration)
		this.mutex.Unlock()
	}()

	return data.value, ok
}

func (this *cache[Key, Value]) resetExpiration(key Key, data mapData[Value], expiration time.Duration) {
	deletedAt := computeExpirationTimestamp(expiration)

	// rebalance the tree with the new deletion time
	this.pqueue.balancedBinTree[data.pos] = pqueueData[Key]{
		key:      key,
		deleteAt: deletedAt,
	}
	this.pqueue.balanceUp(data.pos)
	// get another time the data because the position may have changed
	this.pqueue.balanceDown(this.data[key].pos)
}

func (this *cache[Key, Value]) deleteExpiredItems() {
	this.mutex.Lock()
	for !this.pqueue.Empty() && toDelete(this.pqueue.Front().deleteAt) {
		this.deleteElements(1)
	}
	this.mutex.Unlock()
}

func (this *cache[Key, Value]) deleteElements(n int) {
	for ; !this.pqueue.Empty() && n > 0; n-- {
		key := this.pqueue.Front().key
		this.pqueue.Pop()
		delete(this.data, key)
	}
}

func (this *cache[Key, Value]) Clear() {
	if !this.running {
		return
	}

	this.mutex.Lock()
	this.clear()
	this.mutex.Unlock()
}

func (this *cache[Key, Value]) clear() {
	this.data = make(map[Key]mapData[Value])
	this.pqueue.Clear()
}

func (this *cache[Key, Value]) Stop() {
	if !this.running {
		return
	}

	this.mutex.Lock()
	defer this.mutex.Unlock()

	this.clear()
	this.running = false
	close(this.pingChan)
}

func computeExpirationTimestamp(expiration time.Duration) time.Time {
	if expiration > 0 {
		return time.Now().Add(expiration)
	}
	return time.Time{}
}

func toDelete(t time.Time) bool {
	return t != time.Time{} && t.Before(time.Now())
}

/*
Redefine a priorityQueue because we need to rebalance the tree sometimes
and we need to know the position from an external point of view
*/

func NewPriorityQueue[T any](comp func(a, b T) bool, swap func(a, b T)) *priorityQueue[T] {
	return &priorityQueue[T]{
		comp:            comp,            // true: a<b | false: a>=b
		swap:            swap,            // called every time two data in the pqueue are swaped
		balancedBinTree: make([]T, 0, 3), // arbitrary value
	}
}

type priorityQueue[T any] struct {
	comp            func(T, T) bool
	swap            func(T, T)
	balancedBinTree []T
}

func (this priorityQueue[T]) Empty() bool { return this.Size() == 0 }
func (this priorityQueue[T]) Size() int   { return len(this.balancedBinTree) }
func (this priorityQueue[T]) Front() T    { return this.balancedBinTree[0] }
func (this *priorityQueue[T]) Push(info T) {
	this.balancedBinTree = append(this.balancedBinTree, info)
	this.balanceUp(this.Size() - 1)
}
func (this *priorityQueue[T]) Pop() {
	s := this.balancedBinTree
	l := this.Size() - 1

	this.swap(s[l], s[0])

	s[0], s[l] = s[l], s[0]
	this.balancedBinTree = s[:l]
	this.balanceDown(0)
}
func (this *priorityQueue[T]) Clear() {
	this.balancedBinTree = make([]T, 0, 3)
}

func (this *priorityQueue[T]) balanceUp(n int) {
	if n == 0 {
		return
	}
	parent := this.parent(n)
	if this.comp(this.balancedBinTree[n], this.balancedBinTree[parent]) {
		this.swap(this.balancedBinTree[parent], this.balancedBinTree[n])
		this.balancedBinTree[n], this.balancedBinTree[parent] = this.balancedBinTree[parent], this.balancedBinTree[n]
		this.balanceUp(parent)
		return
	}
}

func (this *priorityQueue[T]) balanceDown(n int) {
	left := this.left(n)
	right := this.right(n)

	if left >= this.Size() {
		return
	}
	if right >= this.Size() {
		// no right, just check left
		if this.comp(this.balancedBinTree[left], this.balancedBinTree[n]) {
			this.swap(this.balancedBinTree[left], this.balancedBinTree[n])
			this.balancedBinTree[n], this.balancedBinTree[left] = this.balancedBinTree[left], this.balancedBinTree[n]
			this.balanceDown(left)
		}
		return
	}

	if this.comp(this.balancedBinTree[left], this.balancedBinTree[right]) {
		// left < right
		if this.comp(this.balancedBinTree[left], this.balancedBinTree[n]) {
			this.swap(this.balancedBinTree[left], this.balancedBinTree[n])
			this.balancedBinTree[n], this.balancedBinTree[left] = this.balancedBinTree[left], this.balancedBinTree[n]
			this.balanceDown(left)
			return
		}
	} else {
		// left >= right
		if this.comp(this.balancedBinTree[right], this.balancedBinTree[n]) {
			this.swap(this.balancedBinTree[right], this.balancedBinTree[n])
			this.balancedBinTree[n], this.balancedBinTree[right] = this.balancedBinTree[right], this.balancedBinTree[n]
			this.balanceDown(right)
			return
		}
	}
}

func (this priorityQueue[T]) parent(n int) int { return (n - 1) / 2 }
func (this priorityQueue[T]) left(n int) int   { return n*2 + 1 }
func (this priorityQueue[T]) right(n int) int  { return n*2 + 2 }
