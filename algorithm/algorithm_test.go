package algorithm

import (
	"testing"

	. "github.com/AlexandreChamard/go_generics/builtin"
)

func TestSort(t *testing.T) {
	sliceInt := []int{42, 12, 4, -5, 0, 69, 20}
	Sort(sliceInt, Less[int])

	intResult := []int{-5, 0, 4, 12, 20, 42, 69}
	if SliceEqual(sliceInt, intResult, Equal[int]) {
		t.Fatalf("expected %v got %v", intResult, sliceInt)
	}
}

func TestUnique(t *testing.T) {
	type a struct {
		b int
	}

	eqf := func(a, b a) bool { return a.b == b.b }

	slice1 := []a{{42}, {12}, {4}, {-5}, {0}, {69}, {20}}
	slice2 := []a{{42}, {12}, {4}, {-5}, {0}, {69}, {4}, {20}}
	if !Unique(slice1, eqf) {
		t.Fatalf("Unique(%v): expected %v got %v", slice1, true, false)
	}
	if Unique(slice2, eqf) {
		t.Fatalf("Unique(%v): expected %v got %v", slice2, false, true)
	}
}

func TestUnique_(t *testing.T) {
	type a struct {
		b int
	}

	slice1 := []a{{42}, {12}, {4}, {-5}, {0}, {69}, {20}}
	slice2 := []a{{42}, {12}, {4}, {-5}, {0}, {69}, {4}, {20}}
	if !Unique_(slice1) {
		t.Fatalf("Unique(%v): expected %v got %v", slice1, true, false)
	}
	if Unique_(slice2) {
		t.Fatalf("Unique(%v): expected %v got %v", slice2, false, true)
	}
}
