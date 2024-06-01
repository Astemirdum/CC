package pkg

type Stack[T any] interface {
	Pop() T
	Push(value T)
	Top() T
	Len() int
}

func NewStack[T any]() Stack[T] {
	return &stack[T]{}
}

type stack[T any] struct {
	arr []T
}

func (s *stack[T]) Pop() T {
	value := s.arr[len(s.arr)-1]
	s.arr = s.arr[:len(s.arr)-1]
	return value
}

func (s *stack[T]) Push(value T) {
	s.arr = append(s.arr, value)
}

func (s *stack[T]) Top() T {
	return s.arr[len(s.arr)-1]
}

func (s *stack[T]) Len() int {
	return len(s.arr)
}
