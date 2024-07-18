package common

import (
	llq "github.com/PlayerR9/MyGoLib/ListLike/Queuer"
	lls "github.com/PlayerR9/MyGoLib/ListLike/Stacker"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

type IteratorNode struct {
	Node  Noder
	Depth int
}

type DFSIterator struct {
	root  Noder
	stack *lls.LinkedStack[*IteratorNode]
}

// Consume implements the common.Iterater interface.
//
// Only returns the common.ErrExhaustedIter.
func (iter *DFSIterator) Consume() (*IteratorNode, error) {
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

		iter.stack.Push(&IteratorNode{
			Node:  c,
			Depth: top.Depth + 1,
		})
	}

	return top, nil
}

// Restart implements the common.Iterater interface.
func (iter *DFSIterator) Restart() {
	iter.stack.Clear()

	if iter.root != nil {
		iter.stack.Push(&IteratorNode{
			Node:  iter.root,
			Depth: 0,
		})
	}
}

func NewDFSIterator(tree Treer) *DFSIterator {
	var root Noder

	if tree != nil {
		root = tree.Root()
	}

	iter := &DFSIterator{
		root:  root,
		stack: lls.NewLinkedStack[*IteratorNode](),
	}

	if root != nil {
		iter.stack.Push(&IteratorNode{
			Node:  root,
			Depth: 0,
		})
	}

	return iter
}

type BFSIterator struct {
	root  Noder
	queue *llq.LinkedQueue[*IteratorNode]
}

// Consume implements the common.Iterater interface.
//
// Only returns the common.ErrExhaustedIter.
func (iter *BFSIterator) Consume() (*IteratorNode, error) {
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

		iter.queue.Enqueue(&IteratorNode{
			Node:  c,
			Depth: first.Depth + 1,
		})
	}

	return first, nil
}

// Restart implements the common.Iterater interface.
func (iter *BFSIterator) Restart() {
	iter.queue.Clear()

	if iter.root != nil {
		iter.queue.Enqueue(&IteratorNode{
			Node:  iter.root,
			Depth: 0,
		})
	}
}

func NewBFSIterator(tree Treer) *BFSIterator {
	var root Noder

	if tree != nil {
		root = tree.Root()
	}

	iter := &BFSIterator{
		root:  root,
		queue: llq.NewLinkedQueue[*IteratorNode](),
	}

	if root != nil {
		iter.queue.Enqueue(&IteratorNode{
			Node:  root,
			Depth: 0,
		})
	}

	return iter
}
