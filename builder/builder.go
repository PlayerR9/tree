package builder

import (
	uc "github.com/PlayerR9/MyGoLib/Units/common"
	lls "github.com/PlayerR9/stack"
	tn "github.com/PlayerR9/tree"
	tr "github.com/PlayerR9/tree/tree"
)

// NextsFunc is a function that returns the next elements.
//
// Parameters:
//   - elem: The element to get the next elements from.
//   - info: The info of the element.
//
// Returns:
//   - []*tn.TreeNode[T]: The next elements.
//   - error: An error if the function fails.
type NextsFunc[T any] func(elem *tn.TreeNode[T], info tr.Infoer) ([]*tn.TreeNode[T], error)

// Builder is a struct that builds a tree.
type Builder[T any] struct {
	// info is the info of the builder.
	info tr.Infoer

	// f is the next function.
	f NextsFunc[T]
}

// SetInfo sets the info of the builder.
//
// Parameters:
//   - info: The info to set.
func (b *Builder[T]) SetInfo(info tr.Infoer) {
	b.info = info
}

// SetNextFunc sets the next function of the builder.
//
// Parameters:
//   - f: The function to set.
func (b *Builder[T]) SetNextFunc(f NextsFunc[T]) {
	b.f = f
}

// MakeTree creates a tree from the given element.
//
// Parameters:
//   - elem: The element to start the tree from.
//   - info: The info of the element.
//   - f: The function that, given an element and info, returns the next elements.
//     (i.e., the children of the element).
//
// Returns:
//   - *Tree: The tree created from the element.
//   - error: An error if the next function fails.
//
// Behaviors:
//   - The 'info' parameter is copied for each node and it specifies the initial info
//     before traversing the tree.
func (b *Builder[T]) Build(elem T) (*tr.Tree[*tn.TreeNode[T]], error) {
	if b.f == nil {
		return nil, nil
	}

	// 1. Handle the root node
	root := tn.NewTreeNode(elem)

	nexts, err := b.f(root, b.info)
	if err != nil {
		return nil, err
	}

	tree := tr.NewTree(root)

	if len(nexts) == 0 {
		return tree, nil
	}

	S := lls.NewLinkedStack[*stack_element[T]]()

	for _, next := range nexts {
		se := new_stack_element[T](tree.Root(), next, b.info)

		S.Push(se)
	}

	for {
		top, ok := S.Pop()
		if !ok {
			break
		}

		data, ok := top.get_data()
		uc.Assert(ok, "Missing data")

		top_inf := top.get_info()

		nexts, err := b.f(data, top_inf)
		if err != nil {
			return nil, err
		}

		ok = top.link_to_prev()
		uc.Assert(ok, "Cannot link to previous node")

		if len(nexts) == 0 {
			continue
		}

		top_elem := top.get_elem()

		for _, next := range nexts {
			se := new_stack_element(top_elem, next, top_inf)

			S.Push(se)
		}
	}

	b.Reset()

	tr.RegenerateLeaves(tree)

	return tree, nil
}

// Reset resets the builder.
func (b *Builder[T]) Reset() {
	b.info = nil
	b.f = nil
}
