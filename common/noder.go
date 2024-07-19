package common

import (
	"fmt"

	uc "github.com/PlayerR9/MyGoLib/Units/common"
	utob "github.com/PlayerR9/MyGoLib/Utility/object"
)

// Noder is an interface that represents a node in a tree.
type Noder interface {
	// IsLeaf returns true if the node is a leaf.
	//
	// Returns:
	//   - bool: True if the node is a leaf, false otherwise.
	IsLeaf() bool

	// GetParent returns the parent of the node.
	//
	// Returns:
	//   - Noder: The parent of the node. Nil if the node has no parent.
	GetParent() Noder

	// IsSingleton returns true if the node is a singleton (i.e., has only one child).
	//
	// Returns:
	//   - bool: True if the node is a singleton, false otherwise.
	IsSingleton() bool

	// DeleteChild deletes the child from the children of the node while
	// returning the children of the target node.
	//
	// No nil nodes are returned.
	//
	// Parameters:
	//   - target: The child to remove.
	//
	// Returns:
	//   - []Noder: A slice of the children of the target node.
	//
	// If target is nil or is not of the correct type, it does nothing and returns nil.
	DeleteChild(target Noder) []Noder

	/*
		// Size returns the number of nodes in the tree rooted at n.
		//
		// Returns:
		//   - size: The number of nodes in the tree.
		Size() int
	*/

	// GetFirstChild returns the first child of the node.
	//
	// Returns:
	//   - Noder: The first child of the node. Nil if the node has no children.
	GetFirstChild() Noder

	// ToTree creates a tree from the node.
	//
	// Returns:
	//   - Treer: The tree containing the node. Never returns nil.
	ToTree() Treer

	// AddChild adds a child to the node.
	//
	// Parameters:
	//   - child: The child to add.
	//
	// If child is nil or not of the correct type, it does nothing.
	AddChild(child Noder)

	uc.Copier

	uc.Iterable[Noder]
	utob.Cleaner

	fmt.Stringer
}
