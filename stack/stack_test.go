package stack

import (
	"fmt"
	"testing"
)

func TestStack(t *testing.T) {
	stack := NewStack[int]() // Const value

	for i := 0; i < 10; i++ {
		stack.Push(i)
	}
	fmt.Println(stack)
	for i := 0; i < 10; i++ {
		if stack.Empty() {
			t.Fatalf("%d: stack.Empty(): expected %v got %v", i, false, stack.Empty())
		}
		if stack.Size() != 10-i {
			t.Fatalf("%d: stack.Size(): expected %d got %d", i, 10-i, stack.Size())
		}
		if stack.Top() != 9-i {
			t.Fatalf("%d: stack.Top(): expected %d got %d", i, 9-i, stack.Top())
		}
		stack.Pop()
	}
	if !stack.Empty() {
		t.Fatalf("stack.Empty(): should be empty at the end")
	}
}
