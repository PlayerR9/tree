package tree

import (
	"iter"
	"slices"
)

// TreeNoder is an interface for nodes in the tree.
type TreeNoder interface {
	comparable

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

	// IsSingleton checks whether the node is a singleton.
	//
	// Returns:
	//   - bool: True if the node is a singleton, false otherwise.
	IsSingleton() bool
}

// DeepCopy is a method that deep copies the node.
//
// Parameters:
//   - node: The node to copy.
//
// Returns:
//   - T: The copied node.
func DeepCopy[T interface {
	Child() iter.Seq[T]
	Copy() T
	LinkChildren(children []T)
	TreeNoder
}](node T) T {
	n := node.Copy()

	var children []T

	for child := range node.Child() {
		child_copy := DeepCopy(child)
		children = append(children, child_copy)
	}

	n.LinkChildren(children)

	return n
}

// RootOf returns the root of the given node.
//
// Parameters:
//   - node: The node to get the root of.
//
// Returns:
//   - T: The root of the given node.
func RootOf[T interface {
	GetParent() (T, bool)
	TreeNoder
}](node T) T {
	for {
		parent, ok := node.GetParent()
		if !ok {
			break
		}

		node = parent
	}

	return node
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
func GetNodeLeaves[T interface {
	Child() iter.Seq[T]
	Copy() T
	TreeNoder
}](node T) []T {
	var leaves []T

	stack := []T{node}

	for len(stack) > 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if top.IsLeaf() {
			leaves = append(leaves, top)
		} else {
			for child := range top.Child() {
				stack = append(stack, child)
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
func GetNodeSize[T interface {
	Child() iter.Seq[T]
	TreeNoder
}](node T) int {
	var size int

	stack := []T{node}

	for len(stack) > 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		size++

		for child := range top.Child() {
			stack = append(stack, child)
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
func GetNodeAncestors[T interface {
	GetParent() (T, bool)
	TreeNoder
}](node T) []T {
	var ancestors []T

	for {
		parent, ok := node.GetParent()
		if !ok {
			break
		}

		ancestors = append(ancestors, parent)

		node = parent
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
//   - T: The common ancestor.
//   - bool: True if the nodes have a common ancestor, false otherwise.
func FindCommonAncestor[T interface {
	GetParent() (T, bool)
	TreeNoder
}](n1, n2 T) (T, bool) {
	if n1 == n2 {
		return n1, true
	}

	ancestors1 := GetNodeAncestors(n1)
	ancestors2 := GetNodeAncestors(n2)

	if len(ancestors1) > len(ancestors2) {
		ancestors1, ancestors2 = ancestors2, ancestors1
	}

	for _, node := range ancestors1 {
		if slices.Contains(ancestors2, node) {
			return node, true
		}
	}

	return *new(T), false
}

// Cleanup is used to delete all the children of the given node.
//
// Parameters:
//   - node: The node to delete the children of.
func Cleanup[T interface {
	Cleanup() []T
	TreeNoder
}](node T) {
	queue := node.Cleanup()

	for len(queue) > 0 {
		first := queue[0]
		queue = queue[1:]

		queue = append(queue, first.Cleanup()...)
	}
}
