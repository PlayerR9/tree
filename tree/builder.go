package tree

import (
	"errors"

	gcers "github.com/PlayerR9/go-commons/errors"
	lls "github.com/PlayerR9/listlike/stack"
)

// NextsFunc is a function that returns the next elements.
//
// Parameters:
//   - elem: The element to get the next elements from.
//   - info: The info of the element.
//
// Returns:
//   - []N: The next elements.
//   - error: An error if the function fails.
type NextsFunc[N Noder] func(elem N, info Infoer) ([]N, error)

// builder_stack_element is a stack element.
type builder_stack_element[N Noder] struct {
	// prev is the previous node.
	prev N

	// elem is the current node.
	elem N

	// info is the info of the current node.
	info Infoer
}

// Builder is a struct that builds a tree.
type Builder[N Noder] struct {
	// info is the info of the builder.
	info Infoer

	// f is the next function.
	f NextsFunc[N]
}

// SetInfo sets the info of the builder.
//
// Parameters:
//   - info: The info to set.
func (b *Builder[N]) SetInfo(info Infoer) {
	b.info = info
}

// SetNextFunc sets the next function of the builder.
//
// Parameters:
//   - f: The function to set.
//
// Returns:
//   - error: An error of type *common.ErrInvalidParameter if 'f' is nil.
func (b *Builder[N]) SetNextFunc(f NextsFunc[N]) error {
	if f == nil {
		return gcers.NewErrNilParameter("f")
	}

	b.f = f

	return nil
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
func (b *Builder[N]) Build(root N) (*Tree[N], error) {
	if b.f == nil {
		return nil, gcers.NewErrInvalidUsage(
			errors.New("no next function is set"),
			"Please call Builder.SetNextFunc() before building the tree",
		)
	}

	// 1. Handle the root node
	nexts, err := b.f(root, b.info)
	if err != nil {
		return nil, err
	}

	tree := NewTree(root)

	if len(nexts) == 0 {
		return tree, nil
	}

	S := lls.NewLinkedStack[*builder_stack_element[N]]()

	if b.info == nil {
		for _, next := range nexts {
			se := &builder_stack_element[N]{
				prev: tree.Root(),
				elem: next,
			}

			S.Push(se)
		}

		for {
			top, ok := S.Pop()
			if !ok {
				break
			}

			nexts, err := b.f(top.elem, nil)
			if err != nil {
				return nil, err
			}

			top.prev.AddChild(top.elem)

			if len(nexts) == 0 {
				continue
			}

			for _, next := range nexts {
				se := &builder_stack_element[N]{
					prev: top.elem,
					elem: next,
					info: nil,
				}

				S.Push(se)
			}
		}

	} else {
		for _, next := range nexts {
			se := &builder_stack_element[N]{
				prev: tree.Root(),
				elem: next,
				info: b.info.Copy(),
			}

			S.Push(se)
		}

		for {
			top, ok := S.Pop()
			if !ok {
				break
			}

			nexts, err := b.f(top.elem, top.info)
			if err != nil {
				return nil, err
			}

			top.prev.AddChild(top.elem)

			if len(nexts) == 0 {
				continue
			}

			for _, next := range nexts {
				se := &builder_stack_element[N]{
					prev: top.elem,
					elem: next,
					info: top.info.Copy(),
				}

				S.Push(se)
			}
		}
	}

	b.Reset()

	RegenerateLeaves(tree)

	return tree, nil
}

// Reset resets the builder.
func (b *Builder[N]) Reset() {
	b.info = nil
	b.f = nil
}
