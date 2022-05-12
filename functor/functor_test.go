package functor

import (
	"fmt"
	"testing"
)

func TestFunctor(t *testing.T) {

	functors := []Functor{
		makeFunctor1(
			func(i int) {
				t.Log(i)
			},
			42,
		),
		makeFunctor2(
			func(i int, s string) {
				t.Log(i, s)
			},
			1, "foobar",
		),
		makeFunctor5(
			func(i int, s string, b bool, pf *float64, a any) {
				t.Log(i, s, b, pf, a)
			},
			1, "foobar", true, func(f float64) *float64 { return &f }(1.23), "it's working very well",
		),
	}

	for _, f := range functors {
		f.Apply()
	}
}

func BenchmarkFunctor(b *testing.B) {
	b.ReportAllocs()
	b.SetBytes(1)
	b.StartTimer()

	functors := []Functor{}
	for i := 0; i < 1000; i++ {

		functors = append(functors,
			makeFunctor1(
				func(i int) {
					fmt.Println(i)
				},
				42,
			),
			makeFunctor2(
				func(i int, s string) {
					fmt.Println(i, s)
				},
				1, "foobar",
			),
			makeFunctor5(
				func(i int, s string, b bool, pf *float64, a any) {
					fmt.Println(i, s, b, pf, a)
				},
				1, "foobar", true, func(f float64) *float64 { return &f }(1.23), "it's working very well",
			),
		)
	}
	for _, f := range functors {
		f.Apply()
	}

	b.StopTimer()
}
