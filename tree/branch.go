package tree

import (
	"fmt"
	"iter"

	dbg "github.com/PlayerR9/go-debug/assert"
)

// Branch represents a branch in a tree.
type Branch[T Noder] struct {
	// from_node is the node from which the branch starts.
	from_node T

	// to_node is the node to which the branch ends.
	to_node T
}

// Copy is a method that returns a copy of the branch.
//
// Returns:
//   - *Branch: A copy of the branch.
func (b *Branch[T]) Copy() *Branch[T] {
	from_copy := b.from_node.Copy().(T)
	to_copy := b.to_node.Copy().(T)

	b_copy := &Branch[T]{
		from_node: from_copy,
		to_node:   to_copy,
	}

	return b_copy
}

// Iterator implements the uc.Iterable interface.
func (b *Branch[T]) Iterator() iter.Seq[T] {
	fn := func(yield func(T) bool) {
		for n := b.from_node; Noder(n) != Noder(b.to_node); {
			fc := n.GetFirstChild()
			dbg.AssertNotNil(fc, "fc")

			tmp := dbg.AssertConv[T](fc, "fc")

			if !yield(tmp) {
				return
			}
		}
	}

	return fn
}

// Slice implements the uc.Slicer interface.
func (b *Branch[T]) Slice() []T {
	var slice []T

	n := b.from_node
	for Noder(n) != Noder(b.to_node) {
		slice = append(slice, n)

		fc := n.GetFirstChild()
		// uc.Assert(fc != nil, "first child should not be nil")

		tmp := fc.(T)
		// uc.AssertF(ok, "first child should be of type %T, got %T", *new(T), fc)

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
