package tree

import (
	"iter"
)

// Branch represents a branch in a tree.
type Branch[T interface {
	Child() iter.Seq[T]
	GetParent() (T, bool)
	TreeNoder
}] struct {
	// from_node is the node from which the branch starts.
	from_node T

	// to_node is the node to which the branch ends.
	to_node T
}

// Node is a method that scans the branch from the top to the bottom.
//
// Returns:
//   - iter.Seq[T]: A sequence of nodes from the top to the bottom.
func (b Branch[T]) Node() iter.Seq[T] {
	fn := func(yield func(T) bool) {
		for n := b.from_node; n != b.to_node; {
			for child := range n.Child() {
				if !yield(child) {
					return
				}

				break
			}
		}
	}

	return fn
}

// NewBranch works like GetAncestors but includes the node itself.
//
// The nodes are returned as a slice where [0] is the root node
// and [len(branch)-1] is the leaf node.
//
// Parameters:
//   - node: The node to get the ancestors of.
//
// Returns:
//   - *Branch: The branch from the node to the root.
//   - error: An error if the creation fails (i.e., the node is not of type T).
func NewBranch[T interface {
	Child() iter.Seq[T]
	GetParent() (T, bool)
	TreeNoder
}](node T) (*Branch[T], error) {
	branch := &Branch[T]{
		to_node: node,
	}

	for {
		parent, ok := node.GetParent()
		if !ok {
			break
		}

		node = parent
	}

	branch.from_node = node

	return branch, nil
}
