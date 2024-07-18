package common

type Treer interface {
	// Root returns the root of the tree.
	//
	// Returns:
	//   - T: The root of the tree. Nil if the tree does not have a root.
	Root() Noder

	// Size returns the number of nodes in the tree.
	//
	// Returns:
	//   - int: The number of nodes in the tree.
	Size() int

	// GetLeaves returns all the leaves of the tree.
	//
	// Returns:
	//   - []Noder: A slice of the leaves of the tree. Nil if the tree does not have a root.
	//
	// Behaviors:
	//   - It returns the leaves that are stored in the tree. Make sure to call
	//     any update function before calling this function if the tree has been modified
	//     unexpectedly.
	GetLeaves() []Noder

	/*
		// IsSingleton returns true if the tree has only one node.
		//
		// Returns:
		//   - bool: True if the tree has only one node, false otherwise.
		IsSingleton() bool



		// SetChildren sets the children of the root of the tree.
		//
		// Parameters:
		//   - children: The children to set.
		//
		// Returns:
		//   - error: An error of type *ErrMissingRoot if the tree does not have a root.
		SetChildren(children []*Tree[T]) error





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




		// RegenerateLeaves regenerates the leaves of the tree and returns them.
		//
		// Returns:
		//   - []T: The newly generated leaves of the tree.
		//   - error: An error if the leaves could not be generated or the nodes are not of type T.
		//
		// Behaviors:
		//   - The leaves are updated in a DFS order.
		//   - Expensive operation; use it only when necessary (i.e., leaves changed unexpectedly.)
		//   - This also updates the size of the tree.
		RegenerateLeaves() ([]*TreeNode[T], error)

		// UpdateLeaves updates the leaves of the tree and returns them.
		//
		// Returns:
		//   - []T: The newly generated leaves of the tree.
		//   - error: An error if the leaves could not be generated or the nodes are not of type T.
		//
		// Behaviors:
		//   - The leaves are updated in a DFS order.
		//   - Less expensive than RegenerateLeaves. However, if nodes has been deleted
		//     from the tree, this may give unexpected results.
		//   - This also updates the size of the tree.
		UpdateLeaves() ([]*TreeNode[T], error)
		// HasChild returns true if the tree has the given child in any of its nodes
		// in a BFS order.
		//
		// The filter must assume that the node is never nil.
		//
		// Parameters:
		//   - filter: The filter to apply.
		//
		// Returns:
		//   - bool: True if the tree has the child, false otherwise.
		//   - error: An error if the child is not of type T.
		HasChild(filter us.PredicateFilter[*TreeNode[T]]) (bool, error)

		// FilterChildren returns all the children of the tree that satisfy the given filter
		// in a BFS order.
		//
		// The filter must assume that the node is never nil.
		//
		// Parameters:
		//   - filter: The filter to apply.
		//
		// Returns:
		//   - []T: A slice of the children that satisfy the filter.
		//   - error: An error if iterating over the children fails.
		FilterChildren(filter us.PredicateFilter[*TreeNode[T]]) ([]*TreeNode[T], error)

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
		PruneBranches(filter us.PredicateFilter[*TreeNode[T]]) bool

		// SearchNodes searches for the first node that satisfies the given filter in a BFS order.
		//
		// Parameters:
		//   - f: The filter to apply.
		//
		// Returns:
		//   - T: The node that satisfies the filter.
		//   - error: An error if the node is not found or the iteration fails.
		//
		// Errors:
		//   - *common.ErrNotFound: If the node is not found.
		//   - error: The error returned by the iteration function.
		SearchNodes(f us.PredicateFilter[*TreeNode[T]]) (*TreeNode[T], error)

		// DeleteBranchContaining deletes the branch containing the given node.
		//
		// Parameters:
		//   - n: The node to delete.
		//
		// Returns:
		//   - error: An error if the node is not a part of the tree.
		DeleteBranchContaining(n *TreeNode[T]) error

		// PruneTree prunes the tree using the given filter.
		//
		// Parameters:
		//   - filter: The filter to use to prune the tree.
		//
		// Returns:
		//   - bool: True if no nodes were pruned, false otherwise.
		//   - error: An error if the iteration fails.
		Prune(filter us.PredicateFilter[*TreeNode[T]]) (bool, error)

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
		SkipFilter(filter us.PredicateFilter[*TreeNode[T]]) (forest []*Tree[T])

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
		ProcessLeaves(f uc.EvalManyFunc[*TreeNode[T], *TreeNode[T]]) error
		// GetDirectChildren returns the direct children of the root of the tree.
		//
		// Children are never nil.
		//
		// Returns:
		//   - []T: A slice of the direct children of the root. Nil if the tree does not have a root.
		GetDirectChildren() []*TreeNode[T]

		// ExtractBranch extracts the branch of the tree that contains the given leaf.
		//
		// Parameters:
		//   - leaf: The leaf to extract the branch from.
		//   - delete: If true, the branch is deleted from the tree.
		//
		// Returns:
		//   - *Branch[T]: A pointer to the branch extracted. Nil if the leaf is not a part
		//     of the tree.
		ExtractBranch(leaf *TreeNode[T], delete bool) (*Branch[T], error)

		// InsertBranch inserts the given branch into the tree.
		//
		// Parameters:
		//   - branch: The branch to insert.
		//
		// Returns:
		//   - bool: True if the branch was inserted, false otherwise.
		//   - error: An error if the insertion fails.
		InsertBranch(branch *Branch[T]) (bool, error)

		ob.Cleaner
		uc.Copier
	*/

	// SetLeaves is an internal function that sets the leaves of the tree. Nil
	// or invalid leaves are ignored.
	//
	// WARNING: Never call this function unless you know what you are doing.
	//
	// Parameters:
	//   - leaves: The leaves to set.
	//   - size: The size of the tree.
	SetLeaves(leaves []Noder, size int)
}
