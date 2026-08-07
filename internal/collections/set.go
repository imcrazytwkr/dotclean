package collections

type Set[V comparable] struct {
	content map[V]struct{}
}

func (s *Set[V]) init() {
	if s.content == nil {
		s.content = make(map[V]struct{})
	}
}

func (s *Set[V]) Add(value V) {
	s.init()
	s.content[value] = struct{}{}
}

func (s *Set[V]) Contains(value V) bool {
	_, ok := s.content[value]
	return ok
}

func (s *Set[V]) IsEmpty() bool {
	return s.Size() == 0
}

func (s *Set[V]) Size() int {
	return len(s.content)
}
