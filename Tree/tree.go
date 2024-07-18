package Tree

import (
	uc "github.com/PlayerR9/MyGoLib/Units/common"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
)

// Tree is a generic data structure that represents a tree.
type Tree[T any] struct {
	// root is the root of the tree.
	root *TreeNode[T]

	// leaves is the leaves of the tree.
	leaves []*TreeNode[T]

	// size is the number of nodes in the tree.
	size int
}

// Cleanup implements the object.Cleaner interface.
func (t *Tree[T]) Cleanup() {
	root := t.root
	if root == nil {
		return
	}

	root.Cleanup()

	t.root = nil
}

// Copy implements the common.Copier interface.
func (t *Tree[T]) Copy() uc.Copier {
	root := t.root
	if root == nil {
		tree := &Tree[T]{
			root:   nil,
			leaves: nil,
			size:   0,
		}

		return tree
	}

	var tree *Tree[T]

	root_copy := root.Copy().(*TreeNode[T])

	tree = &Tree[T]{
		root:   root_copy,
		leaves: root_copy.GetLeaves(),
		size:   t.size,
	}

	return tree
}

// NewTree creates a new tree with the given value as the root.
//
// Parameters:
//   - data: The value of the root.
//
// Returns:
//   - *Tree: A pointer to the newly created tree.
func NewTree[T any](root *TreeNode[T]) *Tree[T] {
	if root == nil {
		tree := &Tree[T]{
			root:   nil,
			leaves: nil,
			size:   0,
		}

		return tree
	}

	var leaves []*TreeNode[T]
	var size int

	ok := root.IsLeaf()
	if ok {
		leaves = []*TreeNode[T]{root}
		size = 1
	} else {
		leaves = root.GetLeaves()
		size = root.Size()
	}

	tree := &Tree[T]{
		root:   root,
		leaves: leaves,
		size:   size,
	}

	return tree
}

// SetChildren sets the children of the root of the tree.
//
// Parameters:
//   - children: The children to set.
//
// Returns:
//   - error: An error of type *ErrMissingRoot if the tree does not have a root.
func (t *Tree[T]) SetChildren(children []*Tree[T]) error {
	children = us.SliceFilter(children, FilterNonNilTree)
	if len(children) == 0 {
		return nil
	}

	root := t.root
	if root == nil {
		return NewErrMissingRoot()
	}

	var leaves, sub_children []*TreeNode[T]

	t.size = 1

	for _, child := range children {
		leaves = append(leaves, child.leaves...)
		t.size += child.Size()

		croot := child.root
		croot.Parent = root

		sub_children = append(sub_children, croot)
	}

	root.LinkChildren(sub_children)

	t.leaves = leaves

	return nil
}

// IsSingleton returns true if the tree has only one node.
//
// Returns:
//   - bool: True if the tree has only one node, false otherwise.
func (t *Tree[T]) IsSingleton() bool {
	return t.size == 1
}

// Size returns the number of nodes in the tree.
//
// Returns:
//   - int: The number of nodes in the tree.
func (t *Tree[T]) Size() int {
	return t.size
}

// Root returns the root of the tree.
//
// Returns:
//   - T: The root of the tree. Nil if the tree does not have a root.
func (t *Tree[T]) Root() *TreeNode[T] {
	return t.root
}

/*

// GetChildren returns all the children of the tree in a DFS order.
//
// Returns:
//   - children: A slice of the values of the children of the tree.
//
// Behaviors:
//   - The root is the first element in the slice.
//   - If the tree does not have a root, it returns nil.
func (t *Tree) GetChildren() (children []T) {
	root := t.root
	if root == nil {
		return nil
	}

	S := Stacker.NewLinkedStack(root)

	for {
		node, ok := S.Pop()
		if !ok {
			break
		}

		children = append(children, node.Data)

		for i := 0; i < len(node.children); i++ {
			current := node.children[i]

			S.Push(current)
		}
	}

	return children
}
*/

// GetLeaves returns all the leaves of the tree.
//
// Returns:
//   - []T: A slice of the leaves of the tree. Nil if the tree does not have a root.
//
// Behaviors:
//   - It returns the leaves that are stored in the tree. Make sure to call
//     any update function before calling this function if the tree has been modified
//     unexpectedly.
func (t *Tree[T]) GetLeaves() []*TreeNode[T] {
	return t.leaves
}

// PruneBranches removes all the children of the node that satisfy the given filter.
// The filter is a function that takes the value of a node and returns a boolean.
// If the filter returns true for a child, the child is removed along with its children.
//
// Parameters:
//   - filter: The filter to apply.
//
// Returns:
//   - bool: True if the whole tree can be deleted, false otherwise.
//
// Behaviors:
//   - If the root satisfies the filter, the tree is cleaned up.
//   - It is a recursive function.
func (t *Tree[T]) PruneBranches(filter us.PredicateFilter[*TreeNode[T]]) bool {
	if filter == nil {
		return false
	}

	root := t.root
	if root == nil {
		return true
	}

	highest, ok := rec_prune_func(filter, nil, root)
	if ok {
		return true
	}

	t.leaves = highest.GetLeaves()
	t.size = highest.Size()

	return false
}

// SkipFunc removes all the children of the tree that satisfy the given filter
// without removing any of their children. Useful for removing unwanted nodes from the tree.
//
// Parameters:
//   - filter: The filter to apply.
//
// Returns:
//   - []*Tree: A slice of pointers to the trees obtained after removing the nodes.
//
// Behaviors:
//   - If this function returns only one tree, this is the updated tree. But, if
//     it returns more than one tree, then we have deleted the root of the tree and
//     obtained a forest.
func (t *Tree[T]) SkipFilter(filter us.PredicateFilter[*TreeNode[T]]) (forest []*Tree[T]) {
	frontier := make([]*TreeNode[T], len(t.leaves))
	copy(frontier, t.leaves)

	seen := make(map[*TreeNode[T]]bool)
	var leaves []*TreeNode[T]

	f := func(n *TreeNode[T]) bool {
		return !seen[n]
	}

	for len(frontier) > 0 {
		leaf := frontier[0]
		seen[leaf] = true

		// Remove any node that has been seen from the frontier.
		frontier = us.SliceFilter(frontier, f)

		ok := filter(leaf)

		parent := leaf.Parent

		if !ok {
			if parent == nil {
				// We reached the root
				frontier = frontier[1:]
			} else {
				ok := leaf.IsLeaf()
				if ok {
					leaves = append(leaves, leaf)
				}

				if !seen[parent] {
					frontier[0] = parent
				} else {
					frontier = frontier[1:]
				}
			}
		} else {
			children := leaf.RemoveNode()

			if len(children) != 0 {
				// We obtained a forest as we reached the root

				for i := 0; i < len(children); i++ {
					child := children[i]

					tree := NewTree(child)

					forest = append(forest, tree)
				}

				// We reached the root
				frontier = frontier[1:]
			} else {
				if !seen[parent] {
					frontier[0] = parent
				} else {
					frontier = frontier[1:]
				}

				t.size--
			}
		}
	}

	if len(forest) == 0 {
		t.leaves = leaves

		forest = []*Tree[T]{t}
	}

	return
}

// replaceLeafWithTree is a helper function that replaces a leaf with a tree.
//
// Parameters:
//   - at: The index of the leaf to replace.
//   - children: The children of the leaf.
//
// Behaviors:
//   - The leaf is replaced with the children.
//   - The size of the tree is updated.
func (t *Tree[T]) replaceLeafWithTree(at int, values []*TreeNode[T]) {
	leaf := t.leaves[at]

	// Make the subtree
	leaf.LinkChildren(values)

	// Update the size of the tree
	t.size += len(values) - 1

	// Replace the current leaf with the leaf's children
	sub_leaves := leaf.GetLeaves()

	if at == len(t.leaves)-1 {
		t.leaves = append(t.leaves[:at], sub_leaves...)
	} else if at == 0 {
		t.leaves = append(sub_leaves, t.leaves[at+1:]...)
	} else {
		t.leaves = append(t.leaves[:at], append(sub_leaves, t.leaves[at+1:]...)...)
	}
}

// ProcessLeaves applies the given function to the leaves of the tree and replaces
// the leaves with the children returned by the function.
//
// Parameters:
//   - f: The function to apply to the leaves.
//
// Returns:
//   - error: An error returned by the function.
//
// Behaviors:
//   - The function is applied to the leaves in order.
//   - The function must return a slice of values of type T.
//   - If the function returns an error, the process stops and the error is returned.
//   - The leaves are replaced with the children returned by the function.
func (t *Tree[T]) ProcessLeaves(f uc.EvalManyFunc[*TreeNode[T], *TreeNode[T]]) error {
	for i, leaf := range t.leaves {
		children, err := f(leaf)
		if err != nil {
			return err
		}

		if len(children) != 0 {
			t.replaceLeafWithTree(i, children)
		}
	}

	return nil
}

// GetDirectChildren returns the direct children of the root of the tree.
//
// Children are never nil.
//
// Returns:
//   - []T: A slice of the direct children of the root. Nil if the tree does not have a root.
func (t *Tree[T]) GetDirectChildren() []*TreeNode[T] {
	root := t.root
	if root == nil {
		return nil
	}

	var children []*TreeNode[T]

	for c := root.FirstChild; c != nil; c = c.NextSibling {
		uc.Assert(c != nil, "Unexpected nil child")

		children = append(children, c)
	}

	return children
}
