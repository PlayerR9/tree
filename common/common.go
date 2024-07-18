package common

import (
	"slices"
	"strings"

	lls "github.com/PlayerR9/MyGoLib/ListLike/Stacker"
	us "github.com/PlayerR9/MyGoLib/Units/common"
)

/*
// TreeFormatter is a formatter that formats the tree.
type TreeFormatter struct {
	// spacing is the spacing between nodes.
	spacing string

	// leaf_prefix is the prefix for leaves.
	leaf_prefix string

	// node_prefix is the prefix for nodes.
	node_prefix string
}

// WithSpacing sets the spacing between nodes.
//
// If spacing is an empty string, it is set to three spaces.
//
// Parameters:
//   - spacing: The spacing between nodes.
//
// Returns:
//   - ffs.Option: The option function.
func WithSpacing(spacing string) ffs.Option {
	size := utf8.RuneCountInString(spacing)
	if size <= 1 {
		spacing = "   "
	}

	p1 := strings.Repeat("─", size-1)
	p2 := strings.Repeat(spacing, size)

	return func(s ffs.Settinger) {
		tf, ok := s.(*TreeFormatter)
		if !ok {
			return
		}

		var builder strings.Builder

		builder.WriteRune('|')
		builder.WriteString(p2)

		tf.spacing = builder.String()
		builder.Reset()

		builder.WriteRune('├')
		builder.WriteString(p1)
		builder.WriteRune(' ')

		tf.leaf_prefix = builder.String()
		builder.Reset()

		builder.WriteRune('└')
		builder.WriteString(p1)
		builder.WriteRune(' ')

		tf.node_prefix = builder.String()
	}
}

*/

const (
	// DefaultSpacing is the default spacing between nodes.
	DefaultSpacing string = "|   "

	// DefaultLeafPrefix is the default prefix for leaves.
	DefaultLeafPrefix string = "├── "

	// DefaultNodePrefix is the default prefix for nodes.
	DefaultNodePrefix string = "└── "
)

// FString implements the FString.FStringer interface.
//
// By default, it uses a three-space indentation.
//
// Format:
//
//	root
//	├── node1
//	│   ├── node2
//	│   └── node3
//	└── node4
//	|   └── node5
//	|
//	| // ...
//	// ...
//
// Parameters:
//   - t: The tree to format.
//
// Returns:
//   - string: The formatted string.
func FString(t Treer) string {
	iter := NewDFSIterator(t)

	elem, err := iter.Consume()
	if err != nil {
		return ""
	}

	var builder strings.Builder

	// Deal with root.
	str := elem.Node.String()

	builder.WriteString(str)

	// Deal with children.

	for {
		node, err := iter.Consume()
		if err != nil {
			break
		}

		builder.WriteRune('\n')

		builder.WriteString(strings.Repeat(DefaultSpacing, node.Depth))

		ok := node.Node.IsLeaf()
		if ok {
			builder.WriteString(DefaultLeafPrefix)
		} else {
			builder.WriteString(DefaultNodePrefix)
		}

		builder.WriteString(node.Node.String())
	}

	return builder.String()
}

// RegenerateLeaves regenerates the leaves of the tree. No op if the tree is nil.
//
// Parameters:
//   - tree: The tree to regenerate the leaves of.
//
// Behaviors:
//   - The leaves are updated in a DFS order.
//   - Expensive operation; use it only when necessary (i.e., leaves changed unexpectedly.)
//   - This also updates the size of the tree.
func RegenerateLeaves(tree Treer) {
	if tree == nil {
		return
	}

	root := tree.Root()
	if root == nil {
		return
	}

	var leaves []Noder

	iter := NewDFSIterator(tree)

	var size int

	for {
		elem, err := iter.Consume()
		if err != nil {
			break
		}

		size++

		ok := elem.Node.IsLeaf()
		if ok {
			leaves = append(leaves, elem.Node)
		}
	}

	tree.SetLeaves(leaves, size)
}

// UpdateLeaves updates the leaves of the tree. No op if the tree is nil.
//
// Parameters:
//   - tree: The tree to update.
//
// Behaviors:
//   - The leaves are updated in a DFS order.
//   - Less expensive than RegenerateLeaves. However, if nodes has been deleted
//     from the tree, this may give unexpected results.
//   - This also updates the size of the tree.
func UpdateLeaves(tree Treer) {
	if tree == nil {
		return
	}

	leaves := tree.GetLeaves()
	if len(leaves) == 0 {
		tree.SetLeaves(nil, 0)
		return
	}

	var new_leaves []Noder
	size := tree.Size() - len(leaves)

	stack := lls.NewArrayStack(leaves...)

	for {
		top, ok := stack.Pop()
		if !ok {
			break
		}

		size++

		ok = top.IsLeaf()
		if ok {
			new_leaves = append(new_leaves, top)
		}
	}

	tree.SetLeaves(new_leaves, size)
}

// HasChild returns true if the tree has the given child in any of its nodes
// in a BFS order.
//
// Parameters:
//   - tree: The tree to filter.
//   - filter: The filter to apply. Must return true iff the node is the one we are looking for.
//     This function must assume node is never nil.
//
// Returns:
//   - bool: True if the tree has the child, false otherwise.
//
// If either tree or filter is nil, false is returned.
func HasChild(tree Treer, filter func(node Noder) bool) bool {
	if tree == nil || filter == nil {
		return false
	}

	iter := NewBFSIterator(tree)

	for {
		value, err := iter.Consume()
		if err != nil {
			break
		}

		ok := filter(value.Node)
		if ok {
			return true
		}
	}

	return false
}

// FilterChildren returns all the children of the tree that satisfy the given filter
// in a BFS order.
//
// Parameters:
//   - tree: The tree to filter.
//   - filter: The filter to apply. Must return true iff the node is the one we want to keep.
//     This function must assume node is never nil.
//
// Returns:
//   - []T: A slice of the children that satisfy the filter.
//   - bool: True if all the nodes are of type T, false otherwise.
//
// If either tree or filter is nil, an empty slice and false are returned.
func FilterChildren[T Noder](tree Treer, filter func(node T) bool) ([]T, bool) {
	if tree == nil || filter == nil {
		return nil, true
	}

	iter := NewBFSIterator(tree)

	var children []T

	for {
		value, err := iter.Consume()
		if err != nil {
			break
		}

		tmp, ok := value.Node.(T)
		if !ok {
			return nil, false
		}

		ok = filter(tmp)
		if ok {
			children = append(children, tmp)
		}
	}

	return children, true
}

// SearchNodes searches for the first node that satisfies the given filter in a BFS order.
//
// Parameters:
//   - tree: The tree to search.
//   - filter: The filter to apply. Must return true iff the node is the one we are looking for.
//     This function must assume node is never nil.
//
// Returns:
//   - T: The node that satisfies the filter.
//   - bool: True if the node was found, false otherwise.
//
// Nodes that are not of type T will be ignored. If either tree or filter is nil, false is returned.
func SearchNodes[T Noder](tree Treer, filter func(node T) bool) (T, bool) {
	if tree == nil || filter == nil {
		return *new(T), false
	}

	iter := NewBFSIterator(tree)

	for {
		value, err := iter.Consume()
		if err != nil {
			break
		}

		tmp, ok := value.Node.(T)
		if ok {
			ok = filter(tmp)
			if ok {
				return tmp, true
			}
		}
	}

	return *new(T), false
}

// DeleteBranchContaining deletes the branch containing the given node.
//
// Parameters:
//   - n: The node to delete.
//
// Returns:
//   - error: An error if the node is not a part of the tree.
func DeleteBranchContaining[T Noder](tree Treer, n *TreeNode[T]) error {
	if n == nil {
		return nil
	}

	root := t.root
	if root == nil {
		return NewErrNodeNotPartOfTree()
	}

	child, parent, hasBranching := FindBranchingPoint(n)
	if !hasBranching {
		if parent != root {
			return NewErrNodeNotPartOfTree()
		}

		t.Cleanup()
	}

	children := parent.DeleteChild(child)

	for i := 0; i < len(children); i++ {
		current := children[i]

		current.Cleanup()

		children[i] = nil
	}

	leaves, err := t.RegenerateLeaves()
	if err != nil {
		return err
	}

	t.leaves = leaves

	return nil
}

// PruneTree prunes the tree using the given filter.
//
// Parameters:
//   - filter: The filter to use to prune the tree.
//
// Returns:
//   - bool: True if no nodes were pruned, false otherwise.
func (t *Tree[T]) Prune(filter us.PredicateFilter[*TreeNode[T]]) bool {
	for t.Size() != 0 {
		target := t.SearchNodes(filter)
		if target == nil {
			return true
		}

		t.DeleteBranchContaining(target)
	}

	return false
}

// ExtractBranch extracts the branch of the tree that contains the given leaf.
//
// Parameters:
//   - leaf: The leaf to extract the branch from.
//   - delete: If true, the branch is deleted from the tree.
//
// Returns:
//   - *Branch[T]: A pointer to the branch extracted. Nil if the leaf is not a part
//     of the tree.
func (t *Tree[T]) ExtractBranch(leaf *TreeNode[T], delete bool) (*Branch[T], error) {
	found := slices.Contains(t.leaves, leaf)
	if !found {
		return nil, nil
	}

	branch := NewBranch(leaf)

	if !delete {
		return branch, nil
	}

	child, parent, ok := FindBranchingPoint(leaf)
	if !ok {
		parent.DeleteChild(child)
	}

	leaves, err := t.RegenerateLeaves()
	if err != nil {
		return nil, err
	}

	t.leaves = leaves

	return branch, nil
}

// InsertBranch inserts the given branch into the tree.
//
// Parameters:
//   - branch: The branch to insert.
//
// Returns:
//   - bool: True if the branch was inserted, false otherwise.
//   - error: An error if the insertion fails.
func (t *Tree[T]) InsertBranch(branch *Branch[T]) (bool, error) {
	if branch == nil {
		return true, nil
	}

	ref := t.root

	if ref == nil {
		otherTree := NewTree(branch.from_node)

		t.root = otherTree.root
		t.leaves = otherTree.leaves
		t.size = otherTree.size

		return true, nil
	}

	from := branch.from_node
	if ref != from {
		return false, nil
	}

	for from != branch.to_node {
		from = from.FirstChild

		var next *TreeNode[T]

		c := ref.FirstChild

		for c != nil && next == nil {
			if c == from {
				next = c
			}

			c = c.FirstChild
		}

		if next == nil {
			break
		}

		// from is a child of the root. Keep going
		ref = next
	}

	// From this point onward, anything from 'from' up to 'to' must be
	// added in the tree as new children.
	ref.AddChild(from)

	prev_size := t.size

	_, err := t.RegenerateLeaves()
	if err != nil {
		return false, err
	}

	ok := t.size != prev_size
	return ok, nil
}
