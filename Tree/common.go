package Tree

import (
	"slices"

	tn "github.com/PlayerR9/treenode"
)

// FindCommonAncestor returns the first common ancestor of the two nodes.
//
// Parameters:
//   - n1: The first node.
//   - n2: The second node.
//
// Returns:
//   - tn.Noder: A pointer to the common ancestor. Nil if no such node is found.
func FindCommonAncestor(n1, n2 tn.Noder) tn.Noder {
	if n1 == nil {
		return n2
	} else if n2 == nil {
		return n1
	} else if n1 == n2 {
		return n1
	}

	ancestors1 := n1.GetAncestors()
	ancestors2 := n2.GetAncestors()

	if len(ancestors1) > len(ancestors2) {
		ancestors1, ancestors2 = ancestors2, ancestors1
	}

	for _, node := range ancestors1 {
		ok := slices.Contains(ancestors2, node)
		if ok {
			return node
		}
	}

	return nil
}

// FindBranchingPoint returns the first node in the path from n to the root
// such that has more than one sibling.
//
// Returns:
//   - tn.Noder: The branching point.
//   - tn.Noder: The parent of the branching point.
//   - bool: True if the node has a branching point, false otherwise.
//
// Behaviors:
//   - If there is no branching point, it returns the root of the tree. However,
//     if n is nil, it returns nil, nil, false and if the node has no parent, it
//     returns nil, n, false.
func FindBranchingPoint(n tn.Noder) (tn.Noder, tn.Noder, bool) {
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

/*

// ExtractData returns the values of the nodes in the slice. This only works if the
// nodes are of type *TreeNode[T].
//
// Parameters:
//   - nodes: The nodes to extract the values from.
//
// Returns:
//   - []T: A slice of the values of the nodes.
func ExtractData[T any](nodes []tn.Noder) []T {
	if len(nodes) == 0 {
		return nil
	}

	var data []T

	for _, node := range nodes {
		val, ok := node.(*TreeNode[T])
		if !ok {
			continue
		}

		data = append(data, val.Data)
	}

	return data
}
*/
