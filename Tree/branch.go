package Tree

import (
	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

// BranchIterator is the pull-based iterator for the branch.
type BranchIterator[T any] struct {
	// from_node is the node from which the branch starts.
	from_node *TreeNode[T]

	// to_node is the node to which the branch ends.
	to_node *TreeNode[T]

	// current is the current node of the iterator.
	current *TreeNode[T]
}

// Consume implements the common.Iterater interface.
//
// This scans from the root node to the leaf node.
//
// *TreeNode[T] is never nil.
func (bi *BranchIterator[T]) Consume() (*TreeNode[T], error) {
	uc.Assert(bi.current != nil, "BranchIterator: current is nil")

	value := bi.current

	if bi.current == bi.to_node {
		return nil, uc.NewErrExhaustedIter()
	}

	bi.current = bi.current.FirstChild
	return value, nil
}

// Restart implements the common.Iterater interface.
func (bi *BranchIterator[T]) Restart() {
	bi.current = bi.from_node
}

// Branch represents a branch in a tree.
type Branch[T any] struct {
	// from_node is the node from which the branch starts.
	from_node *TreeNode[T]

	// to_node is the node to which the branch ends.
	to_node *TreeNode[T]
}

// Copy implements the uc.Copier interface.
func (b *Branch[T]) Copy() uc.Copier {
	from_copy := b.from_node.Copy().(*TreeNode[T])
	to_copy := b.to_node.Copy().(*TreeNode[T])

	b_copy := &Branch[T]{
		from_node: from_copy,
		to_node:   to_copy,
	}

	return b_copy
}

// Iterator implements the uc.Iterable interface.
func (b *Branch[T]) Iterator() uc.Iterater[*TreeNode[T]] {
	iter := &BranchIterator[T]{
		from_node: b.from_node,
		current:   b.from_node,
	}

	return iter
}

// Slice implements the uc.Slicer interface.
func (b *Branch[T]) Slice() []*TreeNode[T] {
	var slice []*TreeNode[T]

	n := b.from_node
	for n != b.to_node {
		slice = append(slice, n)

		n = n.FirstChild
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
func NewBranch[T any](n *TreeNode[T]) *Branch[T] {
	branch := &Branch[T]{
		to_node: n,
	}

	node := n

	for {
		parent := node.GetParent()
		if parent == nil {
			break
		}

		node = parent
	}

	branch.from_node = node

	return branch
}
