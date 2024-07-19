package Tree

/*

// IsSingleton returns true if the tree has only one node.
//
// Returns:
//   - bool: True if the tree has only one node, false otherwise.
func (t *Tree[N]) IsSingleton() bool {
	return t.size == 1
}

// GetChildren returns all the children of the tree in a DFS order.
//
// Returns:
//   - children: A slice of the values of the children of the tree.
//
// Behaviors:
//   - The root is the first element in the slice.
//   - If the tree does not have a root, it returns nil.
func (t *Tree) GetChildren() (children []N) {
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
