package queue

type Queue[T any] interface {
	Empty() bool
	Size() int
	Front() T
	Back() T
	Push(T)
	Pop()
}

func NewQueue[T any]() Queue[T] { return &queue[T]{} }

type queue[T any] []T

func (this queue[T]) Empty() bool  { return this.Size() == 0 }
func (this queue[T]) Size() int    { return len(this) }
func (this queue[T]) Front() T     { return this[0] }
func (this queue[T]) Back() T      { return this[len(this)-1] }
func (this *queue[T]) Push(info T) { *this = append(*this, info) }
func (this *queue[T]) Pop()        { *this = (*this)[1:] }
