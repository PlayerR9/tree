package tree

// TreeNoder is an interface for nodes in the tree.
type TreeNoder interface {
	// String returns the string representation of the node.
	//
	// Returns:
	//   - string: the string representation of the node.
	String() string

	// IsLeaf checks if the node is a leaf.
	//
	// Returns:
	//   - bool: True if the node is a leaf, false otherwise.
	IsLeaf() bool
}
