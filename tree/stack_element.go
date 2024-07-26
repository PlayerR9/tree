package tree

import (
	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

// stack_element is a stack element.
type stack_element struct {
	// prev is the previous node.
	prev Noder

	// elem is the current node.
	elem Noder

	// info is the info of the current node.
	info uc.Copier
}

// new_stack_element creates a new stack element.
//
// Parameters:
//   - prev: The previous node.
//   - data: The current node.
//   - info: The info of the current node.
//
// Returns:
//   - *stackElement: A pointer to the stack element.
func new_stack_element(prev, data Noder, info uc.Copier) *stack_element {
	se := &stack_element{
		prev: prev,
		elem: data,
		info: info,
	}

	return se
}

// get_data returns the data of the stack element.
//
// Returns:
//   - Tree.*TreeNode[T]: The data of the stack element.
//   - bool: True if the data is valid, otherwise false.
func (se *stack_element) get_data() (Noder, bool) {
	if se.elem == nil {
		return nil, false
	}

	return se.elem, true
}

// get_info returns the info of the stack element.
//
// Returns:
//   - common.Copier: The info of the stack element.
func (se *stack_element) get_info() uc.Copier {
	return se.info
}

// link_to_prev links the current node to the previous node.
//
// Returns:
//   - bool: True if the link is successful, otherwise false.
func (se *stack_element) link_to_prev() bool {
	if se.prev == nil {
		return false
	}

	se.prev.AddChild(se.elem)

	return true
}

// get_elem returns the current node.
//
// Returns:
//   - Tree.*TreeNode[T]: The current node.
func (se *stack_element) get_elem() Noder {
	return se.elem
}
