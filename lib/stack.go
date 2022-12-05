package lib

// Stack is your standard First In Last Out stack.
type Stack[T any] struct {
	Data []T
}

func (s *Stack[T]) Push(e T) {
	s.Data = append(s.Data, e)
}

func (s *Stack[T]) Pop() (T, bool) {
	if len(s.Data) == 0 {
		return *new(T), false
	}
	end := len(s.Data) - 1
	e := s.Data[end]
	s.Data = s.Data[:end]
	return e, true
}

func (s *Stack[T]) Peek() (T, bool) {
	if len(s.Data) == 0 {
		return *new(T), false
	}
	return s.Data[len(s.Data)-1], true
}

func (s *Stack[T]) Len() int {
	return len(s.Data)
}
