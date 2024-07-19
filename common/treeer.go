package common

import (
	utob "github.com/PlayerR9/MyGoLib/Utility/object"
)

type Treer interface {
	// Root returns the root of the tree.
	//
	// Returns:
	//   - T: The root of the tree. Nil if the tree does not have a root.
	Root() Noder

	// Size returns the number of nodes in the tree.
	//
	// Returns:
	//   - int: The number of nodes in the tree.
	Size() int

	// GetLeaves returns all the leaves of the tree.
	//
	// Returns:
	//   - []Noder: A slice of the leaves of the tree. Nil if the tree does not have a root.
	//
	// Behaviors:
	//   - It returns the leaves that are stored in the tree. Make sure to call
	//     any update function before calling this function if the tree has been modified
	//     unexpectedly.
	GetLeaves() []Noder

	// SetLeaves is an internal function that sets the leaves of the tree. Nil
	// or invalid leaves are ignored.
	//
	// WARNING: Never call this function unless you know what you are doing.
	//
	// Parameters:
	//   - leaves: The leaves to set.
	//   - size: The size of the tree.
	SetLeaves(leaves []Noder, size int)

	utob.Cleaner
}
