package Tree

// ErrMissingRoot is an error that is returned when the root of a tree is missing.
type ErrMissingRoot struct{}

// Error implements the error interface.
//
// Message: "missing root".
func (e *ErrMissingRoot) Error() string {
	return "missing root"
}

// NewErrMissingRoot creates a new ErrMissingRoot.
//
// Returns:
//   - *ErrMissingRoot: The newly created error.
func NewErrMissingRoot() *ErrMissingRoot {
	e := &ErrMissingRoot{}
	return e
}

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
