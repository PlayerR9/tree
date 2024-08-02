package tree

import (
	"slices"

	uc "github.com/PlayerR9/lib_units/common"
	us "github.com/PlayerR9/lib_units/slices"
)

// RegenerateLeaves regenerates the leaves of the tree. No op if the tree is nil.
//
// Parameters:
//   - tree: The tree to regenerate the leaves of.
//
// Behaviors:
//   - The leaves are updated in a DFS order.
//   - Expensive operation; use it only when necessary (i.e., leaves changed unexpectedly.)
//   - This also updates the size of the tree.
func RegenerateLeaves[N Noder](tree *Tree[N]) {
	if tree == nil {
		return
	}

	var leaves []N

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

	tree.leaves = leaves
	tree.size = size
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
func UpdateLeaves[N Noder](tree *Tree[N]) {
	if tree == nil {
		return
	}

	if len(tree.leaves) == 0 {
		tree.leaves = []N{tree.root}
		tree.size = 1

		return
	}

	var new_leaves []N
	size := tree.size - len(tree.leaves)

	lls := NewLinkedNStack[N]()

	lls.PushMany(tree.leaves)

	for {
		top, ok := lls.Pop()
		if !ok {
			break
		}

		size++

		ok = top.IsLeaf()
		if ok {
			new_leaves = append(new_leaves, top)
		}
	}

	tree.leaves = new_leaves
	tree.size = size
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
func HasChild[N Noder](tree *Tree[N], filter func(node N) bool) bool {
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
func FilterChildren[N Noder](tree *Tree[N], filter func(node N) bool) ([]N, bool) {
	if tree == nil || filter == nil {
		return nil, true
	}

	iter := NewBFSIterator(tree)

	var children []N

	for {
		value, err := iter.Consume()
		if err != nil {
			break
		}

		ok := filter(value.Node)
		if ok {
			children = append(children, value.Node)
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
func SearchNodes[N Noder](tree *Tree[N], filter func(node N) bool) (N, bool) {
	if tree == nil || filter == nil {
		return *new(N), false
	}

	iter := NewBFSIterator(tree)

	for {
		value, err := iter.Consume()
		if err != nil {
			break
		}

		ok := filter(value.Node)
		if ok {
			return value.Node, true
		}
	}

	return *new(N), false
}

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
func FindBranchingPoint(n Noder) (Noder, Noder, bool) {
	if n == nil {
		return nil, nil, false
	}

	parent := n.GetParent()
	if parent == nil {
		return nil, n, false
	}

	var has_branching_point bool

	for !has_branching_point {
		grand_parent := parent.GetParent()
		if grand_parent == nil {
			break
		}

		ok := parent.IsSingleton()
		if !ok {
			has_branching_point = true
		} else {
			n = parent
			parent = grand_parent
		}
	}

	return n, parent, has_branching_point
}

// DeleteBranchContaining deletes the branch containing the given node.
//
// Parameters:
//   - tree: The tree to delete the branch from.
//   - n: The node to delete.
//
// Returns:
//   - error: An error if the node is not a part of the tree.
func DeleteBranchContaining[N Noder](tree *Tree[N], n N) error {
	if tree == nil {
		uc.NewErrNilParameter("tree")
	}

	child, parent, hasBranching := FindBranchingPoint(n)
	if !hasBranching {
		if parent != Noder(tree.root) {
			return NewErrNodeNotPartOfTree()
		}

		tree.Cleanup()
	}

	children := parent.DeleteChild(child)

	for i := 0; i < len(children); i++ {
		current := children[i]

		current.Cleanup()

		children[i] = nil
	}

	RegenerateLeaves(tree)

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
func Prune[N Noder](tree *Tree[N], filter func(node N) bool) bool {
	if tree == nil {
		return false
	}

	for tree.size != 0 {
		target, ok := SearchNodes(tree, filter)
		if !ok {
			return true
		}

		DeleteBranchContaining(tree, target)
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
func ExtractBranch[N Noder](tree *Tree[N], leaf N, delete bool) *Branch[N] {
	if tree == nil {
		return nil
	}

	conv_node := Noder(leaf)

	found := slices.ContainsFunc(tree.leaves, func(node N) bool {
		return Noder(node) == conv_node
	})
	if !found {
		return nil
	}

	branch, _ := NewBranch[N](leaf)
	// uc.AssertErr(err, "NewBranch[%T](%s)", leaf, leaf.String())

	if !delete {
		return branch
	}

	child, parent, ok := FindBranchingPoint(leaf)
	if !ok {
		parent.DeleteChild(child)
	}

	RegenerateLeaves(tree)

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
func InsertBranch[N Noder](tree *Tree[N], branch *Branch[N]) (*Tree[N], error) {
	if branch == nil {
		return tree, nil
	} else if tree == nil {
		return NewTree[N](branch.from_node), nil
	}

	root := tree.root

	var from Noder

	from = branch.from_node

	if Noder(root) != from {
		return tree, nil
	}

	for from != Noder(branch.to_node) {
		from = from.GetFirstChild()

		var next N
		var found bool

		c := root.GetFirstChild()

		for c != nil && !found {
			if c == from {
				tmp := c.(N)
				// uc.AssertF(ok, "from should be of type %T, got %T", *new(N), c)

				next = tmp
				found = true
			}

			c = c.GetFirstChild()
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

	RegenerateLeaves(tree)

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
func rec_prune_first_func[N Noder](filter func(node N) bool, n N) (N, bool) {
	ok := filter(n)

	if ok {
		// Delete all children
		n.Cleanup()

		return n, true
	}

	iter := n.Iterator()
	// uc.Assert(iter != nil, "iter is nil")

	// Handle first node

	node, _ := iter.Consume()
	// uc.AssertErr(err, "iter.Consume()")

	tmp := node.(N)
	// uc.AssertF(ok, "node should be of type %T, got %T", *new(N), node)

	high, ok := rec_prune_first_func(filter, tmp)
	if ok {
		n.DeleteChild(node)

		return high, true
	}

	// Handle other nodes
	highest := high

	for {
		node, err := iter.Consume()
		if err != nil {
			break
		}

		tmp := node.(N)
		// uc.AssertF(ok, "node should be of type %T, got %T", *new(N), node)

		high, ok := rec_prune_func(filter, highest, tmp)
		if !ok {
			continue
		}

		n.DeleteChild(node)

		highest, _ = FindCommonAncestor(highest, high)
		// uc.Assert(ok, "could not find common ancestor")
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
func rec_prune_func[N Noder](filter func(node N) bool, highest N, n N) (N, bool) {
	ok := filter(n)

	if ok {
		// Delete all children
		n.Cleanup()

		ancestors, _ := FindCommonAncestor(highest, n)
		// uc.Assert(ok, "could not find common ancestor")

		return ancestors, true
	}

	iter := n.Iterator()
	// uc.Assert(iter != nil, "iter is nil")

	for {
		node, err := iter.Consume()
		if err != nil {
			break
		}

		tmp, _ := node.(N)
		// uc.AssertF(ok, "node should be of type %T, got %T", *new(N), node)

		high, ok := rec_prune_func(filter, highest, tmp)
		if !ok {
			continue
		}

		n.DeleteChild(node)

		highest, _ = FindCommonAncestor(highest, high)
		// uc.Assert(ok, "could not find common ancestor")
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
func PruneFunc[N Noder](tree *Tree[N], filter func(node N) bool) bool {
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

// rec_snake_traversal is an helper function that returns all the paths
// from n to the leaves of the tree rooted at n.
//
// Returns:
//   - [][]T: A slice of slices of elements.
//
// Behaviors:
//   - The paths are returned in the order of a BFS traversal.
//   - It is a recursive function.
func rec_snake_traversal[N Noder](n N) [][]N {
	ok := n.IsLeaf()
	if ok {
		return [][]N{
			{n},
		}
	}

	var result [][]N

	iter := n.Iterator()
	// uc.Assert(iter != nil, "iter is nil")

	for {
		node, err := iter.Consume()
		if err != nil {
			break
		}

		tmp := node.(N)
		// uc.AssertF(ok, "node should be of type %T, got %T", *new(N), node)

		subResults := rec_snake_traversal(tmp)

		for _, tmp := range subResults {
			tmp = append([]N{n}, tmp...)
			result = append(result, tmp)
		}
	}

	return result
}

// SnakeTraversal returns all the paths from the root to the leaves of the tree.
//
// Returns:
//   - [][]T: A slice of slices of elements. Nil if the tree is empty.
//
// Behaviors:
//   - The paths are returned in the order of a BFS traversal.
func SnakeTraversal[N Noder](tree *Tree[N]) [][]N {
	if tree == nil {
		return nil
	}

	sol := rec_snake_traversal(tree.root)
	return sol
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
func PruneBranches[N Noder](tree *Tree[N], filter func(node N) bool) bool {
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
func SkipFilter[N Noder](tree *Tree[N], filter func(node N) bool) (forest []*Tree[N]) {
	if tree == nil {
		return nil
	} else if filter == nil {
		return []*Tree[N]{tree}
	}

	frontier := make([]Noder, 0, len(tree.leaves))
	for _, leaf := range tree.leaves {
		frontier = append(frontier, leaf)
	}

	seen := make(map[Noder]bool)
	var leaves []N

	f := func(n Noder) bool {
		return !seen[n]
	}

	for len(frontier) > 0 {
		leaf := frontier[0]
		seen[leaf] = true

		// Remove any node that has been seen from the frontier.
		frontier = us.SliceFilter(frontier, f)

		tmp := leaf.(N)
		// uc.AssertF(ok, "leaf should be of type %T, got %T", *new(N), leaf)

		ok := filter(tmp)

		parent := leaf.GetParent()

		if !ok {
			if parent == nil {
				// We reached the root
				frontier = frontier[1:]
			} else {
				ok := leaf.IsLeaf()
				if ok {
					leaves = append(leaves, tmp)
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
					tmp := child.(N)
					// uc.AssertF(ok, "child should be of type %T, got %T", *new(N), child)

					forest = append(forest, NewTree(tmp))
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

		forest = []*Tree[N]{tree}
	}

	return
}

// replaceLeafWithTree is a helper function that replaces a leaf with a tree.
//
// Parameters:
//   - tree: The tree to replace.
//   - at: The index of the leaf to replace.
//   - children: The children of the leaf.
//
// Behaviors:
//   - The leaf is replaced with the children.
//   - The size of the tree is updated.
func replaceLeafWithTree[N Noder](tree *Tree[N], at int, values []Noder) {
	// uc.AssertParam("at", at >= 0 && at < len(tree.leaves), uc.NewErrOutOfBounds(at, 0, len(tree.leaves)))

	leaf := tree.leaves[at]

	// Make the subtree
	leaf.LinkChildren(values)

	// Update the size of the tree
	tree.size += len(values) - 1

	// Replace the current leaf with the leaf's children
	sub_leaves := GetNodeLeaves(leaf)

	if at == len(tree.leaves)-1 {
		tree.leaves = append(tree.leaves[:at], sub_leaves...)
	} else if at == 0 {
		tree.leaves = append(sub_leaves, tree.leaves[at+1:]...)
	} else {
		tree.leaves = append(tree.leaves[:at], append(sub_leaves, tree.leaves[at+1:]...)...)
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
//   - The function must return a slice of values of type N.
//   - If the function returns an error, the process stops and the error is returned.
//   - The leaves are replaced with the children returned by the function.
func ProcessLeaves[N Noder](tree *Tree[N], f func(node N) ([]N, error)) error {
	if f == nil {
		return nil
	}

	for i, leaf := range tree.leaves {
		children, err := f(leaf)
		if err != nil {
			return err
		}

		if len(children) != 0 {
			conv := make([]Noder, 0, len(children))

			for _, child := range children {
				conv = append(conv, child)
			}

			replaceLeafWithTree(tree, i, conv)
		}
	}

	return nil
}
