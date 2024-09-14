package tree

import "iter"

// Tree is a tree.
type Tree[T interface {
	BackwardChild() iter.Seq[T]

	TreeNoder
}] struct {
	// root is the root of the tree.
	root T
}

// String implements the Stringer interface.
func (t Tree[T]) String() string {
	return PrintTree(t.root)
}

// NewTree creates a new tree.
//
// Parameters:
//   - root: the root node of the tree.
//
// Returns:
//   - Tree: the new tree.
func NewTree[T interface {
	BackwardChild() iter.Seq[T]
	TreeNoder
}](root T) Tree[T] {
	return Tree[T]{
		root: root,
	}
}
