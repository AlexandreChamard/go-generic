package algorithm

import (
	. "github.com/AlexandreChamard/go-generic/builtin"
)

func Min[T Ord](a, b T) T {
	if a < b {
		return a
	} else {
		return b
	}
}

func Min_[T Ord](xs ...T) T {
	min := xs[0]
	for i := 1; i < len(xs); i++ {
		if xs[i] < min {
			min = xs[i]
		}
	}
	return min
}

func Max[T Ord](a, b T) T {
	if a >= b {
		return a
	} else {
		return b
	}
}

func Max_[T Ord](xs ...T) T {
	max := xs[0]
	for i := 1; i < len(xs); i++ {
		if xs[i] > max {
			max = xs[i]
		}
	}
	return max
}

func Sort[T any](slice []T, comp Ordf[T]) {
	QuickSort(slice, comp, 0, len(slice)-1)
}

func QuickSort[T any](slice []T, comp Ordf[T], low, high int) {
	if len(slice) <= 1 {
		return
	}
	if low < high {
		// pi is partitioning index, slice[p] is now at right place
		pi := Partition(slice, comp, low, high)
		QuickSort(slice, comp, low, pi-1)  // Before pi
		QuickSort(slice, comp, pi+1, high) // After pi
	}
}

func Partition[T any](slice []T, comp Ordf[T], low, high int) int {
	i := (low - 1)       // index of smaller element
	pivot := slice[high] // pivot

	for j := low; j < high; j++ {
		// if current element is smaller than or equal to pivot
		if comp(slice[j], pivot) {
			// increment index of smaller element
			i = i + 1
			slice[i], slice[j] = slice[j], slice[i]
		}
	}
	slice[i+1], slice[high] = slice[high], slice[i+1]
	return i + 1
}

func SliceEqual[T any](a, b []T, comp Compf[T]) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if comp(a[i], b[i]) {
			return false
		}
	}
	return true
}

func SliceEqual_[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] == b[i] {
			return false
		}
	}
	return true
}

func All[T any](slice []T, comp Eqf[T]) bool {
	for _, a := range slice {
		if !comp(a) {
			return false
		}
	}
	return true
}

func Any[T any](slice []T, comp Eqf[T]) bool {
	for _, a := range slice {
		if comp(a) {
			return true
		}
	}
	return false
}

func Unique[T any](slice []T, comp Compf[T]) bool {
	for i := range slice {
		for j := i + 1; j < len(slice); j++ {
			if comp(slice[i], slice[j]) {
				return false
			}
		}
	}
	return true
}

func Unique_[T comparable](slice []T) bool {
	m := make(map[T]bool, len(slice))
	for _, a := range slice {
		if m[a] {
			return false
		}
		m[a] = true
	}
	return true
}
