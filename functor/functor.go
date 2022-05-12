package functor

type Functor interface {
	Apply()
}

type functor struct {
	f func()
}

func (this functor) Apply() { this.f() }

func makeFunctor1[T1 any](f func(T1), a1 T1) functor {
	return functor{func() { f(a1) }}
}

func makeFunctor2[T1, T2 any](f func(T1, T2), a1 T1, a2 T2) functor {
	return functor{func() { f(a1, a2) }}
}

func makeFunctor3[T1, T2, T3 any](f func(T1, T2, T3), a1 T1, a2 T2, a3 T3) functor {
	return functor{func() { f(a1, a2, a3) }}
}

func makeFunctor4[T1, T2, T3, T4 any](f func(T1, T2, T3, T4), a1 T1, a2 T2, a3 T3, a4 T4) functor {
	return functor{func() { f(a1, a2, a3, a4) }}
}

func makeFunctor5[T1, T2, T3, T4, T5 any](f func(T1, T2, T3, T4, T5), a1 T1, a2 T2, a3 T3, a4 T4, a5 T5) functor {
	return functor{func() { f(a1, a2, a3, a4, a5) }}
}
