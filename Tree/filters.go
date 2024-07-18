package Tree

import (
	tn "github.com/PlayerR9/treenode"
)

// FilterNonNilTree is a filter that returns true if the tree is not nil.
//
// Parameters:
//   - tree: The tree to filter.
//
// Returns:
//   - bool: True if the tree is not nil, false otherwise.
func FilterNonNilTree[T tn.Noder](tree *Tree[T]) bool {
	if tree == nil {
		return false
	}

	return tree.root != nil
}
