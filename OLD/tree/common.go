package tree

import (
	"errors"
	"fmt"
	"iter"
	"slices"

	gcers "github.com/PlayerR9/go-commons/errors"
	gcslc "github.com/PlayerR9/go-commons/slices"
)

// FindBranchingPoint returns the first node in the path from n to the root
// such that has more than one sibling.
//
// Parameters:
//   - n: The node to start the search.
//
// Returns:
//   - Noder: The branching point. Nil if no branching point was found.
//   - Noder: The parent of the branching point. Nil if n is nil.
//   - bool: True if the node has a branching point, false otherwise.
//
// Behaviors:
//   - If there is no branching point, it returns the root of the tree. However,
//     if n is nil, it returns nil, nil, false and if the node has no parent, it
//     returns nil, n, false.
func FindBranchingPoint[T interface {
	Child() iter.Seq[T]
	BackwardChild() iter.Seq[T]
	Copy() T
	GetParent() (T, bool)
	LinkChildren(children []T)
	Noder
}](n T) (T, *T, bool) {
	parent, ok := n.GetParent()
	if !ok {
		return n, nil, false
	}

	var has_branching_point bool

	for !has_branching_point {
		grand_parent, ok := parent.GetParent()
		if !ok {
			break
		}

		ok = parent.IsSingleton()
		if !ok {
			has_branching_point = true
		} else {
			n = parent
			parent = grand_parent
		}
	}

	return n, &parent, has_branching_point
}

// DeleteBranchContaining deletes the branch containing the given node.
//
// Parameters:
//   - tree: The tree to delete the branch from.
//   - n: The node to delete.
//
// Returns:
//   - error: An error if the node is not a part of the tree.
func DeleteBranchContaining[T interface {
	BackwardChild() iter.Seq[T]
	Child() iter.Seq[T]
	Cleanup() []T
	Copy() T
	DeleteChild(child T) []T
	GetParent() (T, bool)
	LinkChildren(children []T)
	Noder
}](tree *Tree[T], n T) error {
	if tree == nil {
		return gcers.NewErrNilParameter("tree")
	}

	child, tmp, hasBranching := FindBranchingPoint(n)
	if tmp == nil {
		return errors.New("not a branch")
	}

	parent := *tmp

	if !hasBranching {
		if parent != tree.root {
			return NewErrNodeNotPartOfTree()
		}

		tree.Cleanup()

		return nil
	}

	children := parent.DeleteChild(child)

	for i := 0; i < len(children); i++ {
		current := children[i]

		Cleanup(current)
	}

	children = children[:0:0]

	tree.RegenerateLeaves()

	return nil
}

// PruneTree prunes the tree using the given filter.
//
// Parameters:
//   - tree: The tree to prune.
//   - filter: The filter to use to prune the tree. Must return true iff the node
//     should be pruned.
//
// Returns:
//   - bool: False if the whole tree can be deleted, true otherwise.
func Prune[T interface {
	BackwardChild() iter.Seq[T]
	Child() iter.Seq[T]
	Cleanup() []T
	Copy() T
	DeleteChild(child T) []T
	GetParent() (T, bool)
	LinkChildren(children []T)
	Noder
}](tree *Tree[T], filter func(node T) bool) bool {
	if tree == nil {
		return false
	}

	for tree.size != 0 {
		target, ok := tree.SearchNodes(filter)
		if !ok {
			return true
		}

		err := DeleteBranchContaining(tree, target)
		if err != nil {
			panic(fmt.Sprintf("failed to delete branch: %v", err))
		}
	}

	return false
}

// ExtractBranch extracts the branch of the tree that contains the given leaf.
//
// Parameters:
//   - tree: The tree to search.
//   - leaf: The leaf to extract the branch from.
//   - delete: If true, the branch is deleted from the tree.
//
// Returns:
//   - *Branch[T]: A pointer to the branch extracted. Nil if the leaf is not a part
//     of the tree. Nil if the leaf is not a part of the tree and delete is false.
//
// Behaviors:
//   - If delete is true, then the branch is deleted from the tree.
func ExtractBranch[T interface {
	BackwardChild() iter.Seq[T]
	Child() iter.Seq[T]
	Cleanup() []T
	Copy() T
	DeleteChild(child T) []T
	GetParent() (T, bool)
	LinkChildren(children []T)
	Noder
}](tree *Tree[T], leaf T, delete bool) *Branch[T] {
	if tree == nil {
		return nil
	}

	_, ok := leaf.GetParent()
	if !ok {
		return nil
	}

	found := slices.Contains(tree.leaves, leaf)
	if !found {
		return nil
	}

	branch, err := NewBranch[T](leaf)
	if err != nil {
		panic(err.Error())
	}

	if !delete {
		return branch
	}

	child, parent, ok := FindBranchingPoint(leaf)
	if parent == nil {
		panic("IMPOSSIBLE: should be a branch")
	}

	if ok {
		_ = (*parent).DeleteChild(child)
		tree.RegenerateLeaves()
	} else {
		tree.Cleanup()
	}

	return branch
}

// InsertBranch inserts the given branch into the tree.
//
// Parameters:
//   - tree: The tree to insert the branch into.
//   - branch: The branch to insert.
//
// Returns:
//   - T: The updated tree.
//   - error: An error if the insertion fails.
func InsertBranch[T interface {
	AddChild(child T)
	BackwardChild() iter.Seq[T]
	Child() iter.Seq[T]
	Cleanup() []T
	Copy() T
	GetFirstChild() (T, bool)
	GetParent() (T, bool)
	LinkChildren(children []T)
	Noder
}](tree *Tree[T], branch *Branch[T]) (*Tree[T], error) {
	if branch == nil {
		return tree, nil
	} else if tree == nil {
		return NewTree[T](branch.from_node), nil
	}

	root := tree.root

	var from T

	from = branch.from_node

	if root != from {
		return tree, nil
	}

	var ok bool

	for from != branch.to_node {
		from, ok = from.GetFirstChild()
		if !ok {
			panic("IMPOSSIBLE: somehow the 'to_node' is not in the 'from' branch")
		}

		var next T
		var found bool

		for !found {
			c, ok := root.GetFirstChild()
			if !ok {
				break
			}

			if c == from {
				next = c
				found = true
			}
		}

		if !found {
			break
		}

		// from is a child of the root. Keep going
		root = next
	}

	// From this point onward, anything from 'from' up to 'to' must be
	// added in the tree as new children.
	root.AddChild(from)

	tree.RegenerateLeaves()

	return tree, nil
}

// rec_prune_func is an helper function that removes all the children of the
// node that satisfy the given filter including all of their children.
//
// Parameters:
//   - filter: The filter to apply.
//   - n: The node to prune.
//
// Returns:
//   - T: A pointer to the highest ancestor of the pruned node.
//   - bool: True if the node satisfies the filter, false otherwise.
//
// Behaviors:
//   - This function is recursive.
func rec_prune_first_func[T interface {
	Child() iter.Seq[T]
	Cleanup() []T
	DeleteChild(node T) []T
	GetParent() (T, bool)
	Noder
}](filter func(node T) bool, n T) (T, bool) {
	ok := filter(n)

	if ok {
		// Delete all children
		n.Cleanup()

		return n, true
	}

	next, close := iter.Pull(n.Child())
	defer close()

	// Handle first node

	node, ok := next()
	if !ok {
		panic("implement me")
	}

	high, ok := rec_prune_first_func(filter, node)
	if ok {
		_ = n.DeleteChild(node)

		return high, true
	}

	// Handle other nodes
	highest := high

	for {
		node, ok := next()
		if !ok {
			break
		}

		high, ok := rec_prune_func(filter, highest, node)
		if !ok {
			continue
		}

		_ = n.DeleteChild(node)

		highest, ok = FindCommonAncestor(highest, high)
		if !ok {
			panic("could not find common ancestor")
		}
	}

	return highest, false
}

// rec_prune_func is an helper function that removes all the children of the
// node that satisfy the given filter including all of their children.
//
// Parameters:
//   - filter: The filter to apply.
//   - n: The node to prune.
//
// Returns:
//   - T: A pointer to the highest ancestor of the pruned node.
//   - bool: True if the node satisfies the filter, false otherwise.
//
// Behaviors:
//   - This function is recursive.
func rec_prune_func[T interface {
	Child() iter.Seq[T]
	Cleanup() []T
	DeleteChild(node T) []T
	GetParent() (T, bool)
	Noder
}](filter func(node T) bool, highest T, n T) (T, bool) {
	ok := filter(n)

	if ok {
		// Delete all children
		_ = n.Cleanup()

		ancestors, ok := FindCommonAncestor(highest, n)
		if !ok {
			panic("could not find common ancestor")
		}

		return ancestors, true
	}

	for child := range n.Child() {
		high, ok := rec_prune_func(filter, highest, child)
		if !ok {
			continue
		}

		_ = n.DeleteChild(child)

		highest, ok = FindCommonAncestor(highest, high)
		if !ok {
			panic("could not find common ancestor")
		}
	}

	return highest, false
}

// PruneFunc removes all the children of the node that satisfy the given filter
// including all of their children. If the filter is nil, nothing is removed.
//
// Parameters:
//   - filter: The filter to apply. Must return true iff the node should be pruned.
//
// Returns:
//   - bool: True if the node satisfies the filter, false otherwise.
//
// Behaviors:
//   - The root node is not pruned.
func PruneFunc[T interface {
	BackwardChild() iter.Seq[T]
	Child() iter.Seq[T]
	Cleanup() []T
	Copy() T
	DeleteChild(child T) []T
	GetParent() (T, bool)
	LinkChildren(children []T)
	Noder
}](tree *Tree[T], filter func(node T) bool) bool {
	if filter == nil {
		return false
	}

	highest, ok := rec_prune_first_func(filter, tree.root)
	if ok {
		return true
	}

	tree.leaves = GetNodeLeaves(highest)
	tree.size = GetNodeSize(highest)

	return false
}

// PruneBranches removes all the children of the node that satisfy the given filter.
// The filter is a function that takes the value of a node and returns a boolean.
// If the filter returns true for a child, the child is removed along with its children.
//
// Parameters:
//   - tree: The tree to prune.
//   - filter: The filter to apply. Must return true iff the node should be pruned.
//
// Returns:
//   - bool: True if the whole tree can be deleted, false otherwise.
//
// Behaviors:
//   - If the root satisfies the filter, the tree is cleaned up.
//   - It is a recursive function.
func PruneBranches[T interface {
	BackwardChild() iter.Seq[T]
	Child() iter.Seq[T]
	Cleanup() []T
	Copy() T
	DeleteChild(child T) []T
	GetParent() (T, bool)
	LinkChildren(children []T)
	Noder
}](tree *Tree[T], filter func(node T) bool) bool {
	if tree == nil || filter == nil {
		return false
	}

	highest, ok := rec_prune_first_func(filter, tree.root)
	if ok {
		return true
	}

	tree.size = GetNodeSize(highest)
	tree.leaves = GetNodeLeaves(highest)

	return false
}

// SkipFunc removes all the children of the tree that satisfy the given filter
// without removing any of their children. Useful for removing unwanted nodes from the tree.
//
// Parameters:
//   - tree: The tree to prune.
//   - filter: The filter to apply.
//
// Returns:
//   - forest: A slice of pointers to the trees obtained after removing the nodes.
//
// Behaviors:
//   - If this function returns only one tree, this is the updated tree. But, if
//     it returns more than one tree, then we have deleted the root of the tree and
//     obtained a forest.
func SkipFilter[T interface {
	BackwardChild() iter.Seq[T]
	Child() iter.Seq[T]
	Cleanup() []T
	Copy() T
	GetParent() (T, bool)
	LinkChildren(children []T)
	RemoveNode() []T
	Noder
}](tree *Tree[T], filter func(node T) bool) (forest []*Tree[T]) {
	if tree == nil {
		return nil
	} else if filter == nil {
		return []*Tree[T]{tree}
	}

	frontier := make([]T, 0, len(tree.leaves))
	for _, leaf := range tree.leaves {
		frontier = append(frontier, leaf)
	}

	seen := make(map[T]bool)
	var leaves []T

	f := func(n T) bool {
		return !seen[n]
	}

	for len(frontier) > 0 {
		leaf := frontier[0]
		seen[leaf] = true

		// Remove any node that has been seen from the frontier.
		frontier = gcslc.SliceFilter(frontier, f)

		ok := filter(leaf)

		parent, has_parent := leaf.GetParent()

		if !ok {
			if !has_parent {
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
				for _, child := range children {
					forest = append(forest, NewTree(child))
				}

				// We reached the root
				frontier = frontier[1:]
			} else {
				if !seen[parent] {
					frontier[0] = parent
				} else {
					frontier = frontier[1:]
				}

				tree.size--
			}
		}
	}

	if len(forest) == 0 {
		tree.leaves = leaves

		forest = []*Tree[T]{tree}
	}

	return
}
