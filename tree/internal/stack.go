package internal

// Stack is a stack of elements of type T.
type Stack[T any] struct {
	// elems is the underlying slice of elements.
	elems []T
}

// Push adds an element to the stack.
//
// Parameters:
//   - elem: the element to add.
//
// If the receiver is nil, then nothing is done.
func (s *Stack[T]) Push(elem T) {
	if s == nil {
		return
	}

	s.elems = append(s.elems, elem)
}

// Pop removes an element from the stack.
//
// Returns:
//   - T: the element removed.
//   - bool: whether the element was removed.
//
// If the receiver is nil or the stack is empty, then the element is not removed.
func (s *Stack[T]) Pop() (T, bool) {
	if s == nil || len(s.elems) == 0 {
		return *new(T), false
	}

	elem := s.elems[len(s.elems)-1]
	s.elems = s.elems[:len(s.elems)-1]

	return elem, true
}
