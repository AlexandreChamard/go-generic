package memory

func Memoize1[Ret any, T1 comparable](f func(T1) Ret) func(T1) Ret {
	memoizer := make(map[T1]Ret)

	return func(t T1) Ret {
		if ret, ok := memoizer[t]; ok {
			return ret
		}
		ret := f(t)
		memoizer[t] = ret
		return ret
	}
}

func Memoize2[Ret any, T1, T2 comparable](f func(T1, T2) Ret) func(T1, T2) Ret {
	memoizer := make(map[struct {
		t1 T1
		t2 T2
	}]Ret)

	return func(t1 T1, t2 T2) Ret {
		args := struct {
			t1 T1
			t2 T2
		}{t1, t2}
		if ret, ok := memoizer[args]; ok {
			return ret
		}
		ret := f(t1, t2)
		memoizer[args] = ret
		return ret
	}
}

func Memoize3[Ret any, T1, T2, T3 comparable](f func(T1, T2, T3) Ret) func(T1, T2, T3) Ret {
	memoizer := make(map[struct {
		t1 T1
		t2 T2
		t3 T3
	}]Ret)

	return func(t1 T1, t2 T2, t3 T3) Ret {
		args := struct {
			t1 T1
			t2 T2
			t3 T3
		}{t1, t2, t3}
		if ret, ok := memoizer[args]; ok {
			return ret
		}
		ret := f(t1, t2, t3)
		memoizer[args] = ret
		return ret
	}
}

func Memoize4[Ret any, T1, T2, T3, T4 comparable](f func(T1, T2, T3, T4) Ret) func(T1, T2, T3, T4) Ret {
	memoizer := make(map[struct {
		t1 T1
		t2 T2
		t3 T3
		t4 T4
	}]Ret)

	return func(t1 T1, t2 T2, t3 T3, t4 T4) Ret {
		args := struct {
			t1 T1
			t2 T2
			t3 T3
			t4 T4
		}{t1, t2, t3, t4}
		if ret, ok := memoizer[args]; ok {
			return ret
		}
		ret := f(t1, t2, t3, t4)
		memoizer[args] = ret
		return ret
	}
}

func Memoize5[Ret any, T1, T2, T3, T4, T5 comparable](f func(T1, T2, T3, T4, T5) Ret) func(T1, T2, T3, T4, T5) Ret {
	memoizer := make(map[struct {
		t1 T1
		t2 T2
		t3 T3
		t4 T4
		t5 T5
	}]Ret)

	return func(t1 T1, t2 T2, t3 T3, t4 T4, t5 T5) Ret {
		args := struct {
			t1 T1
			t2 T2
			t3 T3
			t4 T4
			t5 T5
		}{t1, t2, t3, t4, t5}
		if ret, ok := memoizer[args]; ok {
			return ret
		}
		ret := f(t1, t2, t3, t4, t5)
		memoizer[args] = ret
		return ret
	}
}
