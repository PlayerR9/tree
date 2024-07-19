package common

// ErrNodeNotPartOfTree is an error that is returned when a node is not part of a tree.
type ErrNodeNotPartOfTree struct{}

// Error implements the error interface.
//
// Message: "node is not part of the tree".
func (e *ErrNodeNotPartOfTree) Error() string {
	return "node is not part of the tree"
}

// NewErrNodeNotPartOfTree creates a new ErrNodeNotPartOfTree.
//
// Returns:
//   - *ErrNodeNotPartOfTree: The newly created error.
func NewErrNodeNotPartOfTree() *ErrNodeNotPartOfTree {
	e := &ErrNodeNotPartOfTree{}
	return e
}
