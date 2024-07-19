package tree

import (
	"fmt"
	"slices"

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

	// GetParent returns the parent of the node. The returned parent is guaranteed to be
	// of the same type as the node whenever it is not nil.
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
	// returning the children of the target node. Each returned child is guaranteed to be
	// of the same type as the target node and not nil.
	//
	// Parameters:
	//   - target: The child to remove.
	//
	// Returns:
	//   - []Noder: A slice of the children of the target node. Nil if either the target
	//     is nil or not of the correct type.
	DeleteChild(target Noder) []Noder

	// GetFirstChild returns the first child of the node. The returned child is guaranteed
	// to be of the same type as the node whenever it is not nil.
	//
	// Returns:
	//   - Noder: The first child of the node. Nil if the node has no children.
	GetFirstChild() Noder

	// AddChild adds the target child to the node. Because this function clears the parent and sibling
	// of the target, it does not add its relatives.
	//
	// Parameters:
	//   - child: The child to add.
	//
	// If child is nil or not of the correct type, it does nothing.
	AddChild(target Noder)

	// LinkChildren links the given children to the node. However, children that are either
	// nil or not of the correct type are ignored.
	//
	// Parameters:
	//   - children: The children to link.
	LinkChildren(children []Noder)

	// RemoveNode removes the node from the tree while shifting the children up one level to
	// maintain the tree structure. The returned children can be used to create a forest of
	// trees if the root node is removed.
	//
	// Finally, the returned children are guaranteed to be of the same type as the node and
	// not nil.
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

// GetAncestors is used to get all the ancestors of the given node. This excludes
// the node itself.
//
// Parameters:
//   - node: The node to get the ancestors of.
//
// Returns:
//   - []T: The ancestors of the node.
//
// This is expensive since ancestors are not stored and so, every time this
// function is called, it has to traverse the tree to find the ancestors. Thus, it is
// recommended to call this function once and then store the ancestors somewhere if needed.
//
// Despite the above, this function does not use recursion and is safe to use.
//
// Finally, no nil nodes are returned.
func GetNodeAncestors[N Noder](node N) []N {
	var ancestors []N

	for {
		parent := node.GetParent()
		if parent == nil {
			break
		}

		tmp, ok := parent.(N)
		uc.AssertF(ok, "parent should be of type %T, got %T", *new(N), parent)

		ancestors = append(ancestors, tmp)

		node = tmp
	}

	slices.Reverse(ancestors)

	return ancestors
}

// FindCommonAncestor returns the first common ancestor of the two nodes.
//
// This function is expensive as it calls GetNodeAncestors two times.
//
// Parameters:
//   - n1: The first node.
//   - n2: The second node.
//
// Returns:
//   - N: The common ancestor.
//   - bool: True if the nodes have a common ancestor, false otherwise.
func FindCommonAncestor[N Noder](n1, n2 N) (N, bool) {
	if Noder(n1) == Noder(n2) {
		return n1, true
	}

	ancestors1 := GetNodeAncestors(n1)
	ancestors2 := GetNodeAncestors(n2)

	if len(ancestors1) > len(ancestors2) {
		ancestors1, ancestors2 = ancestors2, ancestors1
	}

	for _, node := range ancestors1 {
		conv_node := Noder(node)

		ok := slices.ContainsFunc(ancestors2, func(other N) bool {
			return conv_node == Noder(other)
		})
		if ok {
			return node, true
		}
	}

	return *new(N), false
}
