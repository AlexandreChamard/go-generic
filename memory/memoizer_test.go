package memory

import (
	"fmt"
	"testing"
)

func TestMemoizer1(t *testing.T) {
	f1 := func(i int) int {
		fmt.Println("call f1", i)
		return i
	}

	mf1 := Memoize1(f1)
	runTest := func(test string, got int, expected int) {
		if expected != got {
			t.Errorf("%s: expect %d got %d", test, expected, got)
		}
	}

	runTest("mf1(1)", mf1(1), 1)
	runTest("mf1(2)", mf1(2), 2)
	runTest("mf1(3)", mf1(3), 3)
	runTest("memoized mf1(1)", mf1(1), 1)
	runTest("memoized mf1(2)", mf1(2), 2)
	runTest("memoized mf1(3)", mf1(3), 3)
}

func BenchmarkMemoizeFibonnacci(t *testing.B) {
	var mfib func(int64) int64
	fib := func(i int64) int64 {
		if i <= 0 {
			return 0
		}
		if i < 3 {
			return 1
		}
		return mfib(i-1) + mfib(i-2)
	}
	mfib = Memoize1(fib)

	runTest := func(test string, got int64, expected int64) {
		if expected != got {
			t.Errorf("%s: expect %d got %d", test, expected, got)
		}
	}

	runTest("fib(1)", mfib(1), 1)
	runTest("fib(2)", mfib(2), 1)
	runTest("fib(3)", mfib(3), 2)
	runTest("fib(4)", mfib(4), 3)
	runTest("fib(5)", mfib(5), 5)
	runTest("fib(6)", mfib(6), 8)
	runTest("fib(30)", mfib(30), 832040)
	runTest("fib(40)", mfib(40), 102334155)
	runTest("fib(80)", mfib(80), 23416728348467685)

	runTest("fib(-10)", mfib(-10), 0)
	runTest("fib(0)", mfib(0), 0)
}
