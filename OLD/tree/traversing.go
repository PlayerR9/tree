package tree

import (
	"iter"
	"slices"
)

// Pair is a pair of a node and its info.
type Pair[A, B any] struct {
	// Node is the node of the pair.
	Node A

	// Info is the info of the pair.
	Info B
}

// NewPair creates a new pair of a node and its info.
//
// Parameters:
//   - node: The node of the pair.
//   - info: The info of the pair.
//
// Returns:
//   - Pair[A, B]: The new pair.
func NewPair[A, B any](node A, info B) Pair[A, B] {
	return Pair[A, B]{
		Node: node,
		Info: info,
	}
}

// Traverser is the traverser that holds the traversal logic.
type Traverser[T interface {
	Child() iter.Seq[T]
	BackwardChild() iter.Seq[T]
	Copy() T
	LinkChildren(children []T)
	Noder
}, I any] struct {
	// InitFn is the function that initializes the traversal info.
	//
	// Parameters:
	//   - root: The root node of the tree.
	//
	// Returns:
	//   - I: The initial traversal info.
	InitFn func(root T) I

	// DoFn is the function that performs the traversal logic.
	//
	// Parameters:
	//   - node: The current node of the tree.
	//   - info: The traversal info.
	//
	// Returns:
	//   - []Pair[T, I]: The next traversal info.
	//   - error: The error that might occur during the traversal.
	DoFn func(node T, info I) ([]Pair[T, I], error)
}

// ApplyDFS applies the DFS traversal logic to the tree.
//
// Parameters:
//   - t: The tree to apply the traversal logic to.
//   - trav: The traverser that holds the traversal logic.
//
// Returns:
//   - error: The error that might occur during the traversal.
func ApplyDFS[T interface {
	BackwardChild() iter.Seq[T]
	Child() iter.Seq[T]
	Cleanup() []T
	Copy() T
	LinkChildren(children []T)
	Noder
}, I any](t *Tree[T], trav Traverser[T, I]) (I, error) {
	if t == nil {
		return *new(I), nil
	}

	info := trav.InitFn(t.root)

	stack := []Pair[T, I]{NewPair(t.root, info)}

	for len(stack) > 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		nexts, err := trav.DoFn(top.Node, top.Info)
		if err != nil {
			return info, err
		}

		if len(nexts) > 0 {
			slices.Reverse(nexts)
			stack = append(stack, nexts...)
		}
	}

	return info, nil
}

// ApplyBFS applies the BFS traversal logic to the tree.
//
// Parameters:
//   - t: The tree to apply the traversal logic to.
//   - trav: The traverser that holds the traversal logic.
//
// Returns:
//   - error: The error that might occur during the traversal.
func ApplyBFS[T interface {
	BackwardChild() iter.Seq[T]
	Child() iter.Seq[T]
	Cleanup() []T
	Copy() T
	LinkChildren(children []T)
	Noder
}, I any](t *Tree[T], trav Traverser[T, I]) (I, error) {
	if t == nil {
		return *new(I), nil
	}

	info := trav.InitFn(t.root)

	queue := []Pair[T, I]{NewPair(t.root, info)}

	for len(queue) > 0 {
		top := queue[0]
		queue = queue[1:]

		nexts, err := trav.DoFn(top.Node, top.Info)
		if err != nil {
			return info, err
		}

		if len(nexts) > 0 {
			queue = append(queue, nexts...)
		}
	}

	return info, nil
}
