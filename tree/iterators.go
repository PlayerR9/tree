package tree

import (
	"fmt"

	uc "github.com/PlayerR9/lib_units/common"
	llq "github.com/PlayerR9/queue/queue"
	lls "github.com/PlayerR9/stack"
)

// IteratorNode is a node in the iterator.
type IteratorNode[N Noder] struct {
	// Node is the node in the iterator.
	Node N

	// Depth is the depth of the node in the iterator.
	Depth int
}

// DFSIterator is the depth-first search iterator.
type DFSIterator[N Noder] struct {
	// root is the root of the iterator.
	root N

	// stack is the stack of the iterator.
	stack *lls.LinkedStack[*IteratorNode[N]]
}

// Consume implements the common.Iterater interface.
//
// Only returns the common.ErrExhaustedIter.
func (iter *DFSIterator[N]) Consume() (*IteratorNode[N], error) {
	top, ok := iter.stack.Pop()
	if !ok {
		return nil, uc.NewErrExhaustedIter()
	}

	sub_iter := top.Node.Iterator()

	for {
		c, err := sub_iter.Consume()
		if err != nil {
			break
		}

		tmp, ok := c.(N)
		uc.AssertF(ok, "child should be of type %T, got %T", *new(N), c)

		iter.stack.Push(&IteratorNode[N]{
			Node:  tmp,
			Depth: top.Depth + 1,
		})
	}

	return top, nil
}

// Restart implements the common.Iterater interface.
func (iter *DFSIterator[N]) Restart() {
	iter.stack.Clear()

	iter.stack.Push(&IteratorNode[N]{
		Node:  iter.root,
		Depth: 0,
	})
}

// NewDFSIterator creates a new DFSIterator.
//
// Parameters:
//   - tree: The tree to iterate over.
//
// Returns:
//   - *common.DFSIterator: The created DFSIterator. Nil if tree is nil.
func NewDFSIterator[N Noder](tree *Tree[N]) *DFSIterator[N] {
	if tree == nil {
		return nil
	}

	root := tree.root

	iter := &DFSIterator[N]{
		root:  root,
		stack: lls.NewLinkedStack[*IteratorNode[N]](),
	}

	iter.stack.Push(&IteratorNode[N]{
		Node:  root,
		Depth: 0,
	})

	return iter
}

// BFSIterator is the breadth-first search iterator.
type BFSIterator[N Noder] struct {
	// root is the root of the iterator.
	root N

	// queue is the queue of the iterator.
	queue *llq.LinkedQueue[*IteratorNode[N]]
}

// Consume implements the common.Iterater interface.
//
// Only returns the common.ErrExhaustedIter.
func (iter *BFSIterator[N]) Consume() (*IteratorNode[N], error) {
	first, ok := iter.queue.Dequeue()
	if !ok {
		return nil, uc.NewErrExhaustedIter()
	}

	sub_iter := first.Node.Iterator()

	for {
		c, err := sub_iter.Consume()
		if err != nil {
			break
		}

		tmp, ok := c.(N)
		uc.AssertF(ok, "child should be of type %T, got %T", *new(N), c)

		iter.queue.Enqueue(&IteratorNode[N]{
			Node:  tmp,
			Depth: first.Depth + 1,
		})
	}

	return first, nil
}

// Restart implements the common.Iterater interface.
func (iter *BFSIterator[N]) Restart() {
	iter.queue.Clear()

	iter.queue.Enqueue(&IteratorNode[N]{
		Node:  iter.root,
		Depth: 0,
	})
}

// NewBFSIterator creates a new BFSIterator.
//
// Parameters:
//   - tree: The tree to iterate over.
//
// Returns:
//   - *common.BFSIterator: The created BFSIterator. Nil if tree is nil.
func NewBFSIterator[N Noder](tree *Tree[N]) *BFSIterator[N] {
	if tree == nil {
		return nil
	}

	root := tree.root

	iter := &BFSIterator[N]{
		root:  root,
		queue: llq.NewLinkedQueue[*IteratorNode[N]](),
	}

	iter.queue.Enqueue(&IteratorNode[N]{
		Node:  root,
		Depth: 0,
	})

	return iter
}

// Infoer is an interface that provides the info of the element.
type Infoer interface {
	uc.Copier
}

// ObserverFunc is a function that observes a node.
//
// Parameters:
//   - data: The data of the node.
//   - info: The info of the node.
//
// Returns:
//   - bool: True if the traversal should continue, otherwise false.
//   - error: An error if the observation fails.
type ObserverFunc[T Noder] func(data T, info Infoer) (bool, error)

// traversor is a struct that traverses a tree.
type traversor[T Noder] struct {
	// elem is the current node.
	elem T

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
func new_traversor[T Noder](node T, init Infoer) *traversor[T] {
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

// DFS traverses the tree in depth-first order.
//
// Parameters:
//   - tree: The tree to traverse.
//   - init: The initial info.
//   - f: The observer function.
//
// Returns:
//   - error: An error if the traversal fails.
func DFS[N Noder](tree *Tree[N], init Infoer, f ObserverFunc[N]) error {
	if f == nil || tree == nil {
		return nil
	}

	trav := new_traversor(tree.root, init)

	S := lls.NewLinkedStack[*traversor[N]]()
	S.Push(trav)

	for {
		top, ok := S.Pop()
		if !ok {
			break
		}

		ok, err := f(top.elem, top.info)
		if err != nil {
			return err
		} else if !ok {
			continue
		}

		iter := top.elem.Iterator()
		uc.Assert(iter != nil, "Iterator is nil")

		for {
			c, err := iter.Consume()
			if err != nil {
				break
			}

			tmp, ok := c.(N)
			if !ok {
				return fmt.Errorf("node is not a tree: %T", c)
			}

			new_t := new_traversor(tmp, top.info)

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
func BFS[N Noder](tree *Tree[N], init Infoer, f ObserverFunc[N]) error {
	if f == nil || tree == nil {
		return nil
	}

	trav := new_traversor(tree.root, init)

	Q := llq.NewLinkedQueue[*traversor[N]]()

	Q.Enqueue(trav)

	for {
		first, ok := Q.Dequeue()
		if !ok {
			break
		}

		ok, err := f(first.elem, first.info)
		if err != nil {
			return err
		} else if !ok {
			continue
		}

		iter := first.elem.Iterator()
		uc.Assert(iter != nil, "Iterator is nil")

		for {
			c, err := iter.Consume()
			if err != nil {
				break
			}

			tmp, ok := c.(N)
			if !ok {
				return fmt.Errorf("node is not a tree: %T", c)
			}

			new_t := new_traversor(tmp, first.info)

			Q.Enqueue(new_t)
		}
	}

	return nil
}
