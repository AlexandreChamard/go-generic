package priorityqueue

type PriorityQueue[T any] interface {
	Empty() bool
	Size() int
	Front() T
	Push(T)
	Pop()
}

func NewPriorityQueue[T any](comp func(a, b T) bool) PriorityQueue[T] {
	return &priorityQueue[T]{
		comp:            comp,            // true: a<b | false: a>=b
		balancedBinTree: make([]T, 0, 3), // arbitrary value
	}
}

type priorityQueue[T any] struct {
	comp            func(T, T) bool
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
	s[0], s[l] = s[l], s[0]
	this.balancedBinTree = s[:l]
	this.balanceDown(0)
}

func (this *priorityQueue[T]) balanceUp(n int) {
	if n == 0 {
		return
	}
	parent := this.parent(n)
	if this.comp(this.balancedBinTree[n], this.balancedBinTree[parent]) {
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
			this.balancedBinTree[n], this.balancedBinTree[left] = this.balancedBinTree[left], this.balancedBinTree[n]
			this.balanceDown(left)
		}
		return
	}

	if this.comp(this.balancedBinTree[left], this.balancedBinTree[right]) {
		// left < right
		if this.comp(this.balancedBinTree[left], this.balancedBinTree[n]) {
			this.balancedBinTree[n], this.balancedBinTree[left] = this.balancedBinTree[left], this.balancedBinTree[n]
			this.balanceDown(left)
			return
		}
	} else {
		// left >= right
		if this.comp(this.balancedBinTree[right], this.balancedBinTree[n]) {
			this.balancedBinTree[n], this.balancedBinTree[right] = this.balancedBinTree[right], this.balancedBinTree[n]
			this.balanceDown(right)
			return
		}
	}
}

func (this priorityQueue[T]) parent(n int) int { return (n - 1) / 2 }
func (this priorityQueue[T]) left(n int) int   { return n*2 + 1 }
func (this priorityQueue[T]) right(n int) int  { return n*2 + 2 }
