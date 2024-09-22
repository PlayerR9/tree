package tree

import (
	"iter"

	gcers "github.com/PlayerR9/errors"
)

// NextsFunc is a function that returns the next elements.
//
// Parameters:
//   - elem: The element to get the next elements from.
//   - info: The info of the element.
//
// Returns:
//   - []T: The next elements.
//   - error: An error if the function fails.
type NextsFunc[T interface {
	AddChild(child T)
	BackwardChild() iter.Seq[T]
	Child() iter.Seq[T]
	Copy() T
	LinkChildren(children []T)
	TreeNoder
}, I interface {
	Copy() I
}] func(elem T, info I) ([]T, error)

// builder_stack_element is a stack element.
type builder_stack_element[T interface {
	AddChild(child T)
	BackwardChild() iter.Seq[T]
	Child() iter.Seq[T]
	Copy() T
	LinkChildren(children []T)
	TreeNoder
}, I interface {
	Copy() I
}] struct {
	// prev is the previous node.
	prev T

	// elem is the current node.
	elem T

	// info is the info of the current node.
	info I
}

// Builder is a struct that builds a tree.
type Builder[T interface {
	AddChild(child T)
	BackwardChild() iter.Seq[T]
	Child() iter.Seq[T]
	Cleanup() []T
	Copy() T
	LinkChildren(children []T)
	TreeNoder
}, I interface {
	Copy() I
}] struct {
	// info is the info of the builder.
	info I

	// f is the next function.
	f NextsFunc[T, I]
}

// NewBuilder creates a new builder that builds a tree from the given function.
//
// Parameters:
//   - info: The info of the builder.
//   - f: The function that, given an element and info, returns the next elements.
//     (i.e., the children of the element).
//
// Returns:
//   - *Builder: The builder created from the function.
//   - error: An error if the function is nil.
func NewBuilder[T interface {
	AddChild(child T)
	BackwardChild() iter.Seq[T]
	Child() iter.Seq[T]
	Cleanup() []T
	Copy() T
	LinkChildren(children []T)
	TreeNoder
}, I interface {
	Copy() I
}](info I, f NextsFunc[T, I]) (*Builder[T, I], error) {
	if f == nil {
		return nil, gcers.NewErrNilParameter("f")
	}

	return &Builder[T, I]{
		info: info,
		f:    f,
	}, nil
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
func (b *Builder[T, I]) Build(root T) (*Tree[T], error) {
	// 1. Handle the root node
	nexts, err := b.f(root, b.info)
	if err != nil {
		return nil, err
	}

	tree := NewTree(root)

	if len(nexts) == 0 {
		return tree, nil
	}

	stack := make([]builder_stack_element[T, I], 0, len(nexts))

	for _, next := range nexts {
		se := builder_stack_element[T, I]{
			prev: tree.Root(),
			elem: next,
			info: b.info.Copy(),
		}

		stack = append(stack, se)
	}

	for len(stack) > 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		nexts, err := b.f(top.elem, top.info)
		if err != nil {
			return nil, err
		}

		top.prev.AddChild(top.elem)

		if len(nexts) == 0 {
			continue
		}

		for _, next := range nexts {
			se := builder_stack_element[T, I]{
				prev: top.elem,
				elem: next,
				info: top.info.Copy(),
			}

			stack = append(stack, se)
		}
	}

	tree.RegenerateLeaves()

	return tree, nil
}

// Reset resets the builder.
func (b *Builder[T, I]) Reset() {
	b.f = nil
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
func Build[T interface {
	AddChild(child T)
	BackwardChild() iter.Seq[T]
	Child() iter.Seq[T]
	Cleanup() []T
	Copy() T
	LinkChildren(children []T)
	TreeNoder
}](root T, fn func(elem T) ([]T, error)) (*Tree[T], error) {
	if fn == nil {
		return nil, gcers.NewErrInvalidUsage(
			"no next function is set",
			"Please call Builder.SetNextFunc() before building the tree",
		)
	}

	// 1. Handle the root node
	nexts, err := fn(root)
	if err != nil {
		return nil, err
	}

	tree := NewTree(root)

	if len(nexts) == 0 {
		return tree, nil
	}

	// StackElement is a stack element.
	type StackElement struct {
		// prev is the previous node.
		prev T

		// elem is the current node.
		elem T
	}

	stack := make([]StackElement, 0, len(nexts))

	for _, next := range nexts {
		se := StackElement{
			prev: tree.Root(),
			elem: next,
		}

		stack = append(stack, se)
	}

	for len(stack) > 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		nexts, err := fn(top.elem)
		if err != nil {
			return nil, err
		}

		top.prev.AddChild(top.elem)

		if len(nexts) == 0 {
			continue
		}

		for _, next := range nexts {
			se := StackElement{
				prev: top.elem,
				elem: next,
			}

			stack = append(stack, se)
		}
	}

	tree.RegenerateLeaves()

	return tree, nil
}
