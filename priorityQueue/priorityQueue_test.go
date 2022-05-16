package priorityqueue

import (
	"fmt"
	"testing"
)

func TestPriorityQueue(t *testing.T) {
	pqueue := NewPriorityQueue(func(a, b int) bool { return a < b })

	n := 10000

	for i := n - 1; i >= 0; i-- {
		pqueue.Push(i)
	}
	for i := 0; i < n; i++ {
		fmt.Println(pqueue)
		if pqueue.Empty() {
			t.Fatalf("%d: pqueue.Empty(): expected %v got %v", i, false, pqueue.Empty())
		}
		if pqueue.Size() != n-i {
			t.Fatalf("%d: pqueue.Size(): expected %d got %d", i, n-i, pqueue.Size())
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

func BenchmarkPriorityQueue(b *testing.B) {
	b.ReportAllocs()
	b.SetBytes(1)
	b.StartTimer()

	pqueue := NewPriorityQueue(func(a, b int) bool { return a < b })

	for n := 1; n <= (1 << 20); n = n << 1 {
		for i := n - 1; i >= 0; i-- {
			pqueue.Push(i)
		}
		for i := 0; i < n; i++ {
			if pqueue.Empty() {
				b.Errorf("%d: %v", i, pqueue)
				b.Fatalf("%d: pqueue.Empty(): expected %v got %v", i, false, pqueue.Empty())
			}
			if pqueue.Size() != n-i {
				b.Errorf("%d: %v", i, pqueue)
				b.Fatalf("%d: pqueue.Size(): expected %d got %d", i, n-i, pqueue.Size())
			}
			if pqueue.Front() != i {
				b.Errorf("%d: %v", i, pqueue)
				b.Fatalf("%d: pqueue.Front(): expected %d got %d", i, i, pqueue.Front())
			}
			pqueue.Pop()
		}
		if !pqueue.Empty() {
			b.Fatalf("pqueue.Empty(): should be empty at the end")
		}
	}
	b.StopTimer()
}
