package pairs

type Pair[T, U any] struct {
	Left  T
	Right U
}

func New[T, U any](left T, right U) Pair[T, U] {
	return Pair[T, U]{
		Left:  left,
		Right: right,
	}
}

func Flipped[T, U any](p Pair[T, U]) Pair[U, T] {
	return Pair[U, T]{
		Left:  p.Right,
		Right: p.Left,
	}
}

func ToSlice[T any](p Pair[T, T]) []T {
	return []T{p.Left, p.Right}
}

func ToArray[T any](p Pair[T, T]) [2]T {
	return [2]T{p.Left, p.Right}
}
