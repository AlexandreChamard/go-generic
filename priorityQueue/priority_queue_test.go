package priorityqueue

import (
	"fmt"
	"testing"
)

func TestPriorityQueue(t *testing.T) {
	pqueue := NewPriorityQueue(func(a, b int) bool { return a < b })

	for i := 9; i >= 0; i-- {
		pqueue.Push(i)
	}
	for i := 0; i < 10; i++ {
		fmt.Println(pqueue)
		if pqueue.Empty() {
			t.Fatalf("%d: pqueue.Empty(): expected %v got %v", i, false, pqueue.Empty())
		}
		if pqueue.Size() != 10-i {
			t.Fatalf("%d: pqueue.Size(): expected %d got %d", i, 10-i, pqueue.Size())
		}
		if pqueue.Front() != i {
			t.Fatalf("%d: pqueue.Front(): expected %d got %d", i, i, pqueue.Front())
		}
		pqueue.Pop()
	}
	if !pqueue.Empty() {
		t.Fatalf("pqueue.Empty(): should be empty at the end")
	}
}
