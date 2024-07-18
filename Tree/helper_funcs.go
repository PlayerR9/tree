package Tree

import (
	"errors"

	uc "github.com/PlayerR9/MyGoLib/Units/common"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
	tn "github.com/PlayerR9/treenode"
)

// rec_snake_traversal is an helper function that returns all the paths
// from n to the leaves of the tree rooted at n.
//
// Returns:
//   - [][]T: A slice of slices of elements.
//   - error: An error if the iterator fails.
//
// Behaviors:
//   - The paths are returned in the order of a BFS traversal.
//   - It is a recursive function.
func rec_snake_traversal[T tn.Noder](n T) ([][]T, error) {
	uc.AssertParam("n", n != nil, errors.New("recSnakeTraversal: n is nil"))

	ok := n.IsLeaf()
	if ok {
		return [][]T{
			{n},
		}, nil
	}

	iter := n.Iterator()
	if iter == nil {
		return nil, nil
	}

	var result [][]T

	for {
		value, err := iter.Consume()
		ok := uc.IsDone(err)
		if ok {
			break
		} else if err != nil {
			return nil, err
		}

		subResults, err := rec_snake_traversal(value)
		if err != nil {
			return nil, err
		}

		for _, tmp := range subResults {
			tmp = append([]T{n}, tmp...)
			result = append(result, tmp)
		}
	}

	return result, nil
}

// SnakeTraversal returns all the paths from the root to the leaves of the tree.
//
// Returns:
//   - [][]T: A slice of slices of elements.
//   - error: An error if the iterator fails.
//
// Behaviors:
//   - The paths are returned in the order of a BFS traversal.
func (t *Tree[T]) SnakeTraversal() ([][]T, error) {
	root := t.root
	if root == nil {
		return nil, nil
	}

	sol, err := rec_snake_traversal(root)
	return sol, err
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
func rec_prune_func[T tn.Noder](filter us.PredicateFilter[T], highest T, n T) (T, bool) {
	ok := filter(n)

	if ok {
		// Delete all children
		n.Cleanup()

		ancestors := FindCommonAncestor(highest, n)

		return ancestors, true
	}

	iter := n.Iterator()
	if iter == nil {
		return highest, false
	}

	for {
		value, err := iter.Consume()
		ok := uc.IsDone(err)
		if ok {
			break
		} else if err != nil {
			return highest, false
		}

		high, ok := rec_prune_func(filter, highest, value)
		if !ok {
			continue
		}

		n.DeleteChild(value)

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
func (t *Tree[T]) PruneFunc(filter us.PredicateFilter[T]) bool {
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
