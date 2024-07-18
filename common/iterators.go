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

	for c := top.Node.FirstChild; c != nil; c = c.NextSibling {
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

	for c := first.Node.FirstChild; c != nil; c = c.NextSibling {
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
