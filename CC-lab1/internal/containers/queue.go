package containers

type Queue interface {
	Push(value any)
	Pop() any
	Len() int
}

type queue struct {
	arr []any
}

func NewQueue() Queue {
	return new(queue)
}

func (s *queue) Push(value any) {
	s.arr = append(s.arr, value)
}

func (s *queue) Pop() any {
	value := s.arr[0]
	s.arr = s.arr[1:]
	return value
}

func (s *queue) Len() int {
	return len(s.arr)
}
