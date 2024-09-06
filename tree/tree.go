package tree

import (
	gcslc "github.com/PlayerR9/go-commons/slices"
)

// Tree is a generic data structure that represents a tree.
type Tree[N Noder] struct {
	// root is the root of the tree.
	root N

	// leaves is the leaves of the tree.
	leaves []N

	// size is the number of nodes in the tree.
	size int
}

// Copy is a method that copies the tree.
//
// Returns:
//   - *Tree: A copy of the tree.
func (t *Tree[N]) Copy() *Tree[N] {
	root := t.root

	var tree *Tree[N]

	root_copy := root.Copy().(N)

	tree = &Tree[N]{
		root:   root_copy,
		leaves: GetNodeLeaves(root_copy),
		size:   t.size,
	}

	return tree
}

// String implements the fmt.Stringer interface.
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
func (t *Tree[N]) String() string {
	str, _ := PrintTree(t.root)
	// uc.AssertErr(err, "PrintTree(%s)", t.root.String())

	return str
}

// Cleanup implements the object.Cleaner interface.
func (t *Tree[N]) Cleanup() {
	root := t.root

	to_delete := []Noder{root}

	for len(to_delete) > 0 {
		n := to_delete[0]
		to_delete = to_delete[1:]

		to_add := n.Cleanup()

		if len(to_add) > 0 {
			to_delete = append(to_delete, to_add...)
		}
	}
}

// NewTree creates a new tree from the given root.
//
// Parameters:
//   - root: The root of the tree.
//
// Returns:
//   - *Tree[N]: A pointer to the newly created tree. Never nil.
func NewTree[N Noder](root N) *Tree[N] {
	leaves := GetNodeLeaves(root)
	size := GetNodeSize(root)

	tree := &Tree[N]{
		root:   root,
		leaves: leaves,
		size:   size,
	}

	return tree
}

// Root returns the root of the tree.
//
// Returns:
//   - N: The root of the tree.
func (t *Tree[N]) Root() N {
	return t.root
}

// Leaves returns the leaves of the tree.
//
// Returns:
//   - []N: The leaves of the tree.
func (t *Tree[N]) Leaves() []N {
	return t.leaves
}

// Size returns the number of nodes in the tree.
//
// Returns:
//   - int: The number of nodes in the tree.
func (t *Tree[N]) Size() int {
	return t.size
}

// SetChildren sets the children of the root of the tree.
//
// Parameters:
//   - children: The children to set.
//
// Returns:
//   - error: An error of type *ErrMissingRoot if the tree does not have a root.
func (t *Tree[N]) SetChildren(children []*Tree[N]) error {
	children = gcslc.FilterNilValues(children)
	if len(children) == 0 {
		return nil
	}

	root := t.root

	var leaves []N
	var sub_children []Noder

	t.size = 1

	for _, child := range children {
		leaves = append(leaves, child.leaves...)
		t.size += child.Size()

		sub_children = append(sub_children, child.root)
	}

	root.LinkChildren(sub_children)

	t.leaves = leaves

	return nil
}

// GetDirectChildren returns the direct children of the root of the tree.
//
// Children are never nil.
//
// Returns:
//   - []N: A slice of the direct children of the root. Nil if the tree does not have a root.
func (t *Tree[N]) GetDirectChildren() []N {
	var children []N

	iter := t.root.Iterator()
	// uc.Assert(iter != nil, "Unexpected nil iterator")

	for {
		node, err := iter.Consume()
		if err != nil {
			break
		}

		tmp := node.(N)
		// uc.AssertF(ok, "node should be of type %T, got %T", *new(N), node)

		children = append(children, tmp)
	}

	return children
}
