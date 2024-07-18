package Tree

import (
	"github.com/PlayerR9/MyGoLib/ListLike/Queuer"
	"github.com/PlayerR9/MyGoLib/ListLike/Stacker"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

// ObserverFunc is a function that observes a node.
//
// Parameters:
//   - data: The data of the node.
//   - info: The info of the node.
//
// Returns:
//   - bool: True if the traversal should continue, otherwise false.
//   - error: An error if the observation fails.
type ObserverFunc[T any] func(data *TreeNode[T], info Infoer) (bool, error)

// traversor is a struct that traverses a tree.
type traversor[T any] struct {
	// elem is the current node.
	elem *TreeNode[T]

	// info is the info of the current node.
	info Infoer
}

// new_traversor creates a new traversor for the tree.
//
// Parameters:
//   - tree: The tree to traverse.
//   - init: The initial info.
//
// Returns:
//   - Traversor[T, I]: The traversor.
func new_traversor[T any](node *TreeNode[T], init Infoer) *traversor[T] {
	t := &traversor[T]{
		elem: node,
	}

	if init != nil {
		t.info = init.Copy().(Infoer)
	} else {
		t.info = nil
	}

	return t
}

// get_data returns the data of the traversor.
//
// Returns:
//   - T: The data of the traversor.
//   - bool: True if the data is valid, otherwise false.
func (t *traversor[T]) get_data() (*TreeNode[T], bool) {
	if t.elem == nil {
		return nil, false
	}

	return t.elem, true
}

// get_info returns the info of the traversor.
//
// Returns:
//   - uc.Objecter: The info of the traversor.
func (t *traversor[T]) get_info() Infoer {
	return t.info
}

// DFS traverses the tree in depth-first order.
//
// Parameters:
//   - tree: The tree to traverse.
//   - init: The initial info.
//   - f: The observer function.
//
// Returns:
//   - error: An error if the traversal fails.
func DFS[T any](tree *Tree[T], init Infoer, f ObserverFunc[T]) error {
	if f == nil || tree == nil {
		return nil
	}

	root := tree.Root()
	trav := new_traversor(root, init)

	S := Stacker.NewLinkedStack(trav)

	for {
		top, ok := S.Pop()
		if !ok {
			break
		}

		top_data, ok := top.get_data()
		uc.Assert(ok, "Missing data")

		top_inf := top.get_info()

		ok, err := f(top_data, top_inf)
		if err != nil {
			return err
		} else if !ok {
			continue
		}

		for c := top.elem.FirstChild; c != nil; c = c.NextSibling {
			new_t := new_traversor(c, top_inf)

			S.Push(new_t)
		}
	}

	return nil
}

// BFS traverses the tree in breadth-first order.
//
// Parameters:
//   - tree: The tree to traverse.
//   - init: The initial info.
//   - f: The observer function.
//
// Returns:
//   - error: An error if the traversal fails.
func BFS[T any](tree *Tree[T], init Infoer, f ObserverFunc[T]) error {
	if f == nil || tree == nil {
		return nil
	}

	root := tree.Root()
	trav := new_traversor(root, init)

	Q := Queuer.NewLinkedQueue(trav)

	for {
		first, ok := Q.Dequeue()
		if !ok {
			break
		}

		first_data, ok := first.get_data()
		uc.Assert(ok, "Missing data")

		first_inf := first.get_info()

		ok, err := f(first_data, first_inf)
		if err != nil {
			return err
		} else if !ok {
			continue
		}

		for c := first.elem.FirstChild; c != nil; c = c.NextSibling {
			new_t := new_traversor(c, first_inf)

			Q.Enqueue(new_t)
		}
	}

	return nil
}
