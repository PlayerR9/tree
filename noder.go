package treenode

import (
	"fmt"

	uc "github.com/PlayerR9/MyGoLib/Units/common"
	uto "github.com/PlayerR9/MyGoLib/Utility/object"
)

// Noder is an interface that represents a node in a tree.
type Noder interface {
	// IsLeaf returns true if the node is a leaf.
	//
	// Returns:
	//   - bool: True if the node is a leaf, false otherwise.
	IsLeaf() bool

	// IsSingleton returns true if the node is a singleton (i.e., has only one child).
	//
	// Returns:
	//   - bool: True if the node is a singleton, false otherwise.
	IsSingleton() bool

	// Size returns the number of nodes in the tree rooted at n.
	//
	// Returns:
	//   - size: The number of nodes in the tree.
	Size() int

	uc.Copier
	uto.Cleaner
	fmt.Stringer
}
