package containers

// TODO: generic
type Stack interface {
	Pop() any
	Push(value any)
	Top() any
	Len() int
}

func NewStack() Stack {
	return &stack{}
}

type stack struct {
	arr []any
}

func (s *stack) Pop() any {
	value := s.arr[len(s.arr)-1]
	s.arr = s.arr[:len(s.arr)-1]
	return value
}

func (s *stack) Push(value any) {
	s.arr = append(s.arr, value)
}

func (s *stack) Top() any {
	return s.arr[len(s.arr)-1]
}

func (s *stack) Len() int {
	return len(s.arr)
}
