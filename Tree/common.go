package Tree

import (
	"slices"
)

// FindCommonAncestor returns the first common ancestor of the two nodes.
//
// Parameters:
//   - n1: The first node.
//   - n2: The second node.
//
// Returns:
//   - *TreeNode[T]: A pointer to the common ancestor. Nil if no such node is found.
func FindCommonAncestor[T any](n1, n2 *TreeNode[T]) *TreeNode[T] {
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

/*

// ExtractData returns the values of the nodes in the slice. This only works if the
// nodes are of type *TreeNode[T].
//
// Parameters:
//   - nodes: The nodes to extract the values from.
//
// Returns:
//   - []T: A slice of the values of the nodes.
func ExtractData[T any](nodes []*TreeNode[T]) []T {
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
