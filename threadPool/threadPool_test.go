package threadpool

import (
	"fmt"
	"testing"
	"time"

	. "github.com/AlexandreChamard/go_generics/functor"
)

func TestThreadPool(t *testing.T) {
	tp := NewThreadPool(100)

	f := func(i int) {
		time.Sleep(100 * time.Millisecond)
		// fmt.Println("coucou", i)
	}

	for n := 0; n < 100; n++ {
		tp.SubmitPriority(MakeFunctor1(f, n), n/10)
	}

	tp.Stop()
	tp.Wait()
}

func BenchmarkThreadPool(b *testing.B) {
	b.Log("Start benchmark")

	tp := NewThreadPool(1)

	f := func(i int) { fmt.Println("coucou", i) }

	for n := 0; n < 100; n++ {
		tp.SubmitPriority(MakeFunctor1(f, n), n/10)
	}

	// time.Sleep(1 * time.Second)

	tp.Stop()
	tp.Wait()
}