package bintree

// Binary Tree

type BinaryTree[T any] interface {
	Set(T)
	Get() T

	Left() BinaryTree[T]
	Right() BinaryTree[T]

	SetLeft(BinaryTree[T])
	SetRight(BinaryTree[T])
}

func NewBinaryTree[T any]() BinaryTree[T] { return &binaryTree[T]{} }

type binaryTree[T any] struct {
	value       T
	left, right BinaryTree[T]
}

func (this *binaryTree[T]) Set(value T) { this.value = value }
func (this binaryTree[T]) Get() T       { return this.value }

func (this binaryTree[T]) Left() BinaryTree[T]  { return this.left }
func (this binaryTree[T]) Right() BinaryTree[T] { return this.right }

func (this *binaryTree[T]) SetLeft(left BinaryTree[T])   { this.left = left }
func (this *binaryTree[T]) SetRight(right BinaryTree[T]) { this.right = right }
