package Tree

import (
	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

// stackElement is a stack element.
type stackElement[T any] struct {
	// prev is the previous node.
	prev *TreeNode[T]

	// elem is the current node.
	elem *TreeNode[T]

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
func new_stack_element[T any](prev, data *TreeNode[T], info uc.Copier) *stackElement[T] {
	se := &stackElement[T]{
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
func (se *stackElement[T]) get_data() (*TreeNode[T], bool) {
	if se.elem == nil {
		return nil, false
	}

	return se.elem, true
}

// get_info returns the info of the stack element.
//
// Returns:
//   - common.Copier: The info of the stack element.
func (se *stackElement[T]) get_info() uc.Copier {
	return se.info
}

// link_to_prev links the current node to the previous node.
//
// Returns:
//   - bool: True if the link is successful, otherwise false.
func (se *stackElement[T]) link_to_prev() bool {
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
func (se *stackElement[T]) get_elem() *TreeNode[T] {
	return se.elem
}
