package utils

type Int interface {
	~int | ~int32 | ~int64
}
type Uint interface {
	~uint | ~uint32 | ~uint64
}

type AtomicAble[I Int | Uint] interface {
	Load() I
	Add(I) I
	CompareAndSwap(I, I) bool
}

func FetchAndAdd[I Int | Uint, T AtomicAble[I]](a T, delta I) I {
	for {
		v := a.Load()
		if a.CompareAndSwap(v, v+delta) {
			return v
		}
	}
}
