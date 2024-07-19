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
