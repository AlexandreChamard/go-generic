package queue

import (
	"testing"
)

func TestQueue(t *testing.T) {
	queue := MakeQueue[int]() // Const value

	for i := 0; i < 10; i++ {
		queue.Push(i)
	}
	for i := 0; i < 10; i++ {
		if queue.Empty() {
			t.Fatalf("%d: queue.Empty(): expected %v got %v", i, false, queue.Empty())
		}
		if queue.Size() != 10-i {
			t.Fatalf("%d: queue.Size(): expected %d got %d", i, 10-i, queue.Size())
		}
		if queue.Front() != i {
			t.Fatalf("%d: queue.Front(): expected %d got %d", i, i, queue.Front())
		}
		if queue.Back() != 9 {
			t.Fatalf("%d: queue.Back(): expected %d got %d", i, 9, queue.Back())
		}
		queue.Pop()
	}
}
