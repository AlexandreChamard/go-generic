package stack

type Stack[T any] interface {
	Empty() bool
	Size() int
	Top() T
	Push(T)
	Pop()
}

func NewStack[T any]() Stack[T] { return &stack[T]{} }

type stack[T any] []T

func (this stack[T]) Empty() bool  { return this.Size() == 0 }
func (this stack[T]) Size() int    { return len(this) }
func (this stack[T]) Top() T       { return this[len(this)-1] }
func (this stack[T]) TopPtr() *T   { return &this[len(this)-1] }
func (this *stack[T]) Push(info T) { *this = append(*this, info) }
func (this *stack[T]) Pop()        { *this = (*this)[:len(*this)-1] }
