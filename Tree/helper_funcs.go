package Tree

import (
	"errors"

	uc "github.com/PlayerR9/MyGoLib/Units/common"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
)

// rec_snake_traversal is an helper function that returns all the paths
// from n to the leaves of the tree rooted at n.
//
// Returns:
//   - [][]T: A slice of slices of elements.
//
// Behaviors:
//   - The paths are returned in the order of a BFS traversal.
//   - It is a recursive function.
func rec_snake_traversal[T any](n *TreeNode[T]) [][]*TreeNode[T] {
	uc.AssertParam("n", n != nil, errors.New("recSnakeTraversal: n is nil"))

	ok := n.IsLeaf()
	if ok {
		return [][]*TreeNode[T]{
			{n},
		}
	}

	var result [][]*TreeNode[T]

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		subResults := rec_snake_traversal(c)

		for _, tmp := range subResults {
			tmp = append([]*TreeNode[T]{n}, tmp...)
			result = append(result, tmp)
		}
	}

	return result
}

// SnakeTraversal returns all the paths from the root to the leaves of the tree.
//
// Returns:
//   - [][]T: A slice of slices of elements.
//
// Behaviors:
//   - The paths are returned in the order of a BFS traversal.
func (t *Tree[T]) SnakeTraversal() [][]*TreeNode[T] {
	root := t.root
	if root == nil {
		return nil
	}

	sol := rec_snake_traversal(root)
	return sol
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
func rec_prune_func[T any](filter us.PredicateFilter[*TreeNode[T]], highest *TreeNode[T], n *TreeNode[T]) (*TreeNode[T], bool) {
	ok := filter(n)

	if ok {
		// Delete all children
		n.Cleanup()

		ancestors := FindCommonAncestor(highest, n)

		return ancestors, true
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		high, ok := rec_prune_func(filter, highest, c)
		if !ok {
			continue
		}

		n.DeleteChild(c)

		highest = FindCommonAncestor(highest, high)
	}

	return highest, false
}

// PruneFunc removes all the children of the node that satisfy the given filter
// including all of their children.
//
// Parameters:
//   - filter: The filter to apply.
//
// Returns:
//   - bool: True if the node satisfies the filter, false otherwise.
//
// Behaviors:
//   - The root node is not pruned.
func (t *Tree[T]) PruneFunc(filter us.PredicateFilter[*TreeNode[T]]) bool {
	if filter == nil {
		return false
	}

	root := t.root
	if root == nil {
		return false
	}

	highest, ok := rec_prune_func(filter, nil, root)
	if ok {
		return true
	}

	t.leaves = highest.GetLeaves()
	t.size = highest.Size()

	return false
}
