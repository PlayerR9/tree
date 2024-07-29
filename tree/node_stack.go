// Code generated with go generate. DO NOT EDIT.
package tree

import (
	"strconv"
	"strings"

	"github.com/PlayerR9/lib_units/common"
)

// stack_node_N is a node in the linked stack.
type stack_node_N[N Noder] struct {
	value N
	next *stack_node_N[N]
}

// LinkedNStack is a stack of N values implemented without a maximum capacity
// and using a linked list.
type LinkedNStack[N Noder] struct {
	front *stack_node_N[N]
	size int
}

// NewLinkedNStack creates a new linked stack.
//
// Returns:
//   - *LinkedNStack[N]: A pointer to the newly created stack. Never returns nil.
func NewLinkedNStack[N Noder]() *LinkedNStack[N] {
	return &LinkedNStack[N]{
		size: 0,
	}
}

// Push implements the stack.Stacker interface.
//
// Always returns true.
func (s *LinkedNStack[N]) Push(value N) bool {
	node := &stack_node_N[N]{
		value: value,
	}

	if s.front != nil {
		node.next = s.front
	}

	s.front = node
	s.size++

	return true
}

// PushMany implements the stack.Stacker interface.
//
// Always returns the number of values pushed onto the stack.
func (s *LinkedNStack[N]) PushMany(values []N) int {
	if len(values) == 0 {
		return 0
	}

	node := &stack_node_N[N]{
		value: values[0],
	}

	if s.front != nil {
		node.next = s.front
	}

	s.front = node

	for i := 1; i < len(values); i++ {
		node := &stack_node_N[N]{
			value: values[i],
			next:  s.front,
		}

		s.front = node
	}

	s.size += len(values)
	
	return len(values)
}

// Pop implements the stack.Stacker interface.
func (s *LinkedNStack[N]) Pop() (N, bool) {
	if s.front == nil {
		return *new(N), false
	}

	to_remove := s.front
	s.front = s.front.next

	s.size--
	to_remove.next = nil

	return to_remove.value, true
}

// Peek implements the stack.Stacker interface.
func (s *LinkedNStack[N]) Peek() (N, bool) {
	if s.front == nil {
		return *new(N), false
	}

	return s.front.value, true
}

// IsEmpty implements the stack.Stacker interface.
func (s *LinkedNStack[N]) IsEmpty() bool {
	return s.front == nil
}

// Size implements the stack.Stacker interface.
func (s *LinkedNStack[N]) Size() int {
	return s.size
}

// Iterator implements the stack.Stacker interface.
func (s *LinkedNStack[N]) Iterator() common.Iterater[N] {
	var builder common.Builder[N]

	for node := s.front; node != nil; node = node.next {
		builder.Add(node.value)
	}

	return builder.Build()
}

// Clear implements the stack.Stacker interface.
func (s *LinkedNStack[N]) Clear() {
	if s.front == nil {
		return
	}

	prev := s.front

	for node := s.front.next; node != nil; node = node.next {
		prev = node
		prev.next = nil
	}

	prev.next = nil

	s.front = nil
	s.size = 0
}

// GoString implements the stack.Stacker interface.
func (s *LinkedNStack[N]) GoString() string {
	values := make([]string, 0, s.size)
	for node := s.front; node != nil; node = node.next {
		values = append(values, common.StringOf(node.value))
	}

	var builder strings.Builder

	builder.WriteString("LinkedNStack[N][size=")
	builder.WriteString(strconv.Itoa(s.size))
	builder.WriteString(", values=[")
	builder.WriteString(strings.Join(values, ", "))
	builder.WriteString(" →]]")

	return builder.String()
}

// Slice implements the stack.Stacker interface.
//
// The 0th element is the top of the stack.
func (s *LinkedNStack[N]) Slice() []N {
	slice := make([]N, 0, s.size)

	for node := s.front; node != nil; node = node.next {
		slice = append(slice, node.value)
	}

	return slice
}

// Copy implements the stack.Stacker interface.
//
// The copy is a shallow copy.
func (s *LinkedNStack[N]) Copy() common.Copier {
	if s.front == nil {
		return &LinkedNStack[N]{}
	}

	s_copy := &LinkedNStack[N]{
		size: s.size,
	}

	node_copy := &stack_node_N[N]{
		value: s.front.value,
	}

	s_copy.front = node_copy

	prev := node_copy

	for node := s.front.next; node != nil; node = node.next {
		node_copy := &stack_node_N[N]{
			value: node.value,
		}

		prev.next = node_copy

		prev = node_copy
	}

	return s_copy
}

// Capacity implements the stack.Stacker interface.
//
// Always returns -1.
func (s *LinkedNStack[N]) Capacity() int {
	return -1
}

// IsFull implements the stack.Stacker interface.
//
// Always returns false.
func (s *LinkedNStack[N]) IsFull() bool {
	return false
}
