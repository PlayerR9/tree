package tree

import (
	"errors"
	"fmt"

	uc "github.com/PlayerR9/lib_units/common"
)

// BranchIterator is the pull-based iterator for the branch.
type BranchIterator[T Noder] struct {
	// from_node is the node from which the branch starts.
	from_node T

	// to_node is the node to which the branch ends.
	to_node T

	// current is the current node of the iterator.
	current T
}

// Consume implements the common.Iterater interface.
//
// This scans from the root node to the leaf node.
//
// Errors:
//   - *common.ErrExhaustedIter: If the iterator has reached the end of the branch.
//   - error: If the first child of the current node is not of the correct type or an impossible case occurs.
func (bi *BranchIterator[T]) Consume() (T, error) {
	value := bi.current

	if Noder(bi.current) == Noder(bi.to_node) {
		return *new(T), uc.NewErrExhaustedIter()
	}

	fc := bi.current.GetFirstChild()
	if fc == nil {
		return *new(T), errors.New("impossible case: no children but not reached the branch's end point")
	}

	tmp, ok := fc.(T)
	if !ok {
		return *new(T), fmt.Errorf("first child should be of type %T, got %T", *new(T), fc)
	}

	bi.current = tmp

	return value, nil
}

// Restart implements the common.Iterater interface.
func (bi *BranchIterator[T]) Restart() {
	bi.current = bi.from_node
}

// Branch represents a branch in a tree.
type Branch[T Noder] struct {
	// from_node is the node from which the branch starts.
	from_node T

	// to_node is the node to which the branch ends.
	to_node T
}

// Copy implements the uc.Copier interface.
func (b *Branch[T]) Copy() uc.Copier {
	from_copy := b.from_node.Copy().(T)
	to_copy := b.to_node.Copy().(T)

	b_copy := &Branch[T]{
		from_node: from_copy,
		to_node:   to_copy,
	}

	return b_copy
}

// Iterator implements the uc.Iterable interface.
func (b *Branch[T]) Iterator() uc.Iterater[T] {
	iter := &BranchIterator[T]{
		from_node: b.from_node,
		current:   b.from_node,
	}

	return iter
}

// Slice implements the uc.Slicer interface.
func (b *Branch[T]) Slice() []T {
	var slice []T

	n := b.from_node
	for Noder(n) != Noder(b.to_node) {
		slice = append(slice, n)

		fc := n.GetFirstChild()
		uc.Assert(fc != nil, "first child should not be nil")

		tmp, ok := fc.(T)
		uc.AssertF(ok, "first child should be of type %T, got %T", *new(T), fc)

		n = tmp
	}

	slice = append(slice, b.to_node)

	return slice
}

// NewBranch works like GetAncestors but includes the node itself.
//
// The nodes are returned as a slice where [0] is the root node
// and [len(branch)-1] is the leaf node.
//
// Returns:
//   - *Branch: The branch from the node to the root.
//   - error: An error if the creation fails (i.e., the node is not of type T).
func NewBranch[T Noder](n T) (*Branch[T], error) {
	branch := &Branch[T]{
		to_node: n,
	}

	node := n

	for {
		parent := node.GetParent()
		if parent == nil {
			break
		}

		tmp, ok := parent.(T)
		if !ok {
			return nil, fmt.Errorf("parent should be of type %T, got %T", *new(T), parent)
		}

		node = tmp
	}

	branch.from_node = node

	return branch, nil
}
