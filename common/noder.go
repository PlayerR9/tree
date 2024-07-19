package common

import (
	"fmt"

	lls "github.com/PlayerR9/MyGoLib/ListLike/Stacker"
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

	// AddChild adds a child to the node.
	//
	// Parameters:
	//   - child: The child to add.
	//
	// If child is nil or not of the correct type, it does nothing.
	AddChild(child Noder)

	// LinkChildren links the given children to the node.
	//
	// Parameters:
	//   - children: The children to link.
	//
	// Children that are either nil or not of the correct type are ignored.
	LinkChildren(children []Noder)

	// RemoveNode removes the node from the tree while shifting the children up one level to
	// maintain the tree structure.
	//
	// Also, the returned children can be used to create a forest of trees if the root node
	// is removed.
	//
	// Returns:
	//   - []Noder: A slice of pointers to the children of the node iff the node is the root.
	//     Nil otherwise.
	//
	// Example:
	//
	//	// Given the tree:
	//	1
	//	├── 2
	//	└── 3
	//		├── 4
	//		└── 5
	//	└── 6
	//
	//	// The tree after removing node 3:
	//
	//	1
	//	├── 2
	//	└── 4
	//	└── 5
	//	└── 6
	RemoveNode() []Noder

	uc.Copier
	uc.Iterable[Noder]
	utob.Cleaner
	fmt.Stringer
}

// GetNodeLeaves returns the leaves of the given node.
//
// This is expensive as leaves are not stored and so, every time this function is called,
// it has to do a DFS traversal to find the leaves. Thus, it is recommended to call
// this function once and then store the leaves somewhere if needed.
//
// Despite the above, this function does not use recursion and is safe to use.
//
// Finally, no nil nodes are returned.
func GetNodeLeaves[N Noder](node N) []N {
	stack := lls.NewLinkedStack(node)

	var leaves []N

	for {
		top, ok := stack.Pop()
		if !ok {
			break
		}

		if top.IsLeaf() {
			leaves = append(leaves, top)
		} else {
			iter := top.Iterator()
			uc.Assert(iter != nil, "Iterator should not be nil")

			for {
				c, err := iter.Consume()
				if err != nil {
					break
				}

				tmp, ok := c.(N)
				uc.AssertF(ok, "child should be of type %T, got %T", *new(N), c)

				stack.Push(tmp)
			}
		}
	}

	return leaves
}

// Size implements the *TreeNode[T] interface.
//
// This is expensive as it has to traverse the whole tree to find the size of the tree.
// Thus, it is recommended to call this function once and then store the size somewhere if needed.
//
// Despite the above, this function does not use recursion and is safe to use.
//
// Finally, the traversal is done in a depth-first manner.
func GetNodeSize(node Noder) int {
	if node == nil {
		return 0
	}

	stack := lls.NewLinkedStack(node)

	var size int

	for {
		top, ok := stack.Pop()
		if !ok {
			break
		}

		size++

		iter := top.Iterator()
		uc.Assert(iter != nil, "Iterator should not be nil")

		for {
			c, err := iter.Consume()
			if err != nil {
				break
			}

			stack.Push(c)
		}
	}

	return size
}
