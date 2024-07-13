package treenode

import (
	"fmt"

	uc "github.com/PlayerR9/MyGoLib/Units/common"
	uto "github.com/PlayerR9/MyGoLib/Utility/object"
)

// Noder is an interface that represents a node in a tree.
type Noder interface {
	// SetParent sets the parent of the node.
	//
	// Parameters:
	//   - parent: The parent node.
	//
	// Returns:
	//   - bool: True if the parent is set, false otherwise.
	SetParent(parent Noder) bool

	// GetParent returns the parent of the node.
	//
	// Returns:
	//   - Noder: The parent node.
	GetParent() Noder

	// LinkWithParent links the parent with the children. It also links the children
	// with each other.
	//
	// Parameters:
	//   - parent: The parent node.
	//   - children: The children nodes.
	//
	// Behaviors:
	//   - Only valid children are linked while the rest are ignored.
	LinkChildren(children []Noder)

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

	// GetLeaves returns all the leaves of the tree rooted at the node.
	//
	// Should be a DFS traversal.
	//
	// Returns:
	//   - []Noder: A slice of pointers to the leaves of the tree.
	//
	// Behaviors:
	//   - The leaves are returned in the order of a DFS traversal.
	GetLeaves() []Noder

	// GetAncestors returns all the ancestors of the node.
	//
	// This excludes the node itself.
	//
	// Returns:
	//   - []Noder: A slice of pointers to the ancestors of the node.
	//
	// Behaviors:
	//   - The ancestors are returned in the opposite order of a DFS traversal.
	//     Therefore, the first element is the parent of the node.
	GetAncestors() []Noder

	// GetFirstChild returns the first child of the node.
	//
	// Returns:
	//   - Noder: The first child of the node. Nil if the node has no children.
	GetFirstChild() Noder

	// DeleteChild removes the given child from the children of the node.
	//
	// Parameters:
	//   - target: The child to remove.
	//
	// Returns:
	//   - []Noder: A slice of pointers to the children of the node. Nil if the node has no children.
	DeleteChild(target Noder) []Noder

	// Size returns the number of nodes in the tree rooted at n.
	//
	// Returns:
	//   - size: The number of nodes in the tree.
	Size() int

	// AddChild adds a new child to the node with the given data.
	//
	// Parameters:
	//   - child: The child to add.
	//
	// Behaviors:
	//   - If the child is not valid, it is ignored.
	AddChild(child Noder)

	// removeNode removes the node from the tree and shifts the children up
	// in the space occupied by the node.
	//
	// Returns:
	//   - []Noder: A slice of pointers to the children of the node if
	//     the node is the root. Nil otherwise.
	RemoveNode() []Noder

	uc.Iterable[Noder]
	uc.Copier
	uto.Cleaner
	fmt.Stringer
}
