package containers

type SetType interface {
	~byte | ~int
}

type Set[T SetType] map[T]bool

func NewSet[T SetType]() Set[T] {
	return make(Set[T])
}

func (s Set[T]) Add(value T) {
	s[value] = true
}

func (s Set[T]) Remove(value T) {
	delete(s, value)
}
func (s Set[T]) Contains(v T) bool {
	_, ok := s[v]
	return ok
}

func (s Set[T]) Size() int {
	return len(s)
}

func (s Set[T]) ToSet() map[T]bool {
	return s
}

func (s Set[T]) Unite(a Set[T]) Set[T] {
	newSet := make(Set[T])
	for k := range s {
		newSet.Add(k)
	}
	for k := range a {
		newSet.Add(k)
	}
	return newSet
}

func (s Set[T]) Subtract(o Set[T]) Set[T] {
	newSet := make(Set[T])
	for k := range s {
		if !o.Contains(k) {
			newSet.Add(k)
		}
	}
	return newSet
}

func (s Set[T]) Equals(a Set[T]) bool {
	if s.Size() != a.Size() {
		return false
	}
	for k := range s {
		if !a.Contains(k) {
			return false
		}
	}
	return true
}
