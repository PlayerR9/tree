package Tree

// FilterNonNilTree is a filter that returns true if the tree is not nil.
//
// Parameters:
//   - tree: The tree to filter.
//
// Returns:
//   - bool: True if the tree is not nil, false otherwise.
func FilterNonNilTree[T any](tree *Tree[T]) bool {
	if tree == nil {
		return false
	}

	return tree.root != nil
}
