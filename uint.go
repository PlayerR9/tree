// Code generated by go generate; EDIT THIS FILE DIRECTLY

package tree

import (
	"slices"

	"github.com/PlayerR9/MyGoLib/Units/common"
	"github.com/PlayerR9/stack"
	"github.com/PlayerR9/tree/tree"
)

// UintNodeIterator is a pull-based iterator that iterates
// over the children of a UintNode.
type UintNodeIterator struct {
	parent, current *UintNode
}

// Consume implements the common.Iterater interface.
//
// The only error type that can be returned by this function is the *common.ErrExhaustedIter type.
//
// Moreover, the return value is always of type *UintNode and never nil; unless the iterator
// has reached the end of the branch.
func (iter *UintNodeIterator) Consume() (tree.Noder, error) {
	if iter.current == nil {
		return nil, common.NewErrExhaustedIter()
	}

	node := iter.current
	iter.current = iter.current.NextSibling

	return node, nil
}

// Restart implements the common.Iterater interface.
func (iter *UintNodeIterator) Restart() {
	iter.current = iter.parent.FirstChild
}

// UintNode is a node in a tree.
type UintNode struct {
	Parent, FirstChild, NextSibling, LastChild, PrevSibling *UintNode
	Data uint
}

// IsLeaf implements the tree.Noder interface.
func (tn *UintNode) IsLeaf() bool {
	return tn.FirstChild == nil
}

// GetParent implements the tree.Noder interface.
func (tn *UintNode) GetParent() tree.Noder {
	return tn.Parent
}

// IsSingleton implements the tree.Noder interface.
func (tn *UintNode) IsSingleton() bool {
	return tn.FirstChild != nil && tn.FirstChild == tn.LastChild
}


// DeleteChild implements the tree.Noder interface.
func (tn *UintNode) DeleteChild(target tree.Noder) []tree.Noder {
	if target == nil {
		return nil
	}

	tmp, ok := target.(*UintNode)
	if !ok {
		return nil
	}

	children := tn.delete_child(tmp)

	if len(children) == 0 {
		return nil
	}

	for _, child := range children {
		child.PrevSibling = nil
		child.NextSibling = nil
		child.Parent = nil
	}

	tn.FirstChild = nil
	tn.LastChild = nil

	conv := make([]tree.Noder, 0, len(children))

	for _, child := range children {
		conv = append(conv, child)
	}

	return conv
}

// GetFirstChild implements the tree.Noder interface.
func (tn *UintNode) GetFirstChild() tree.Noder {
	return tn.FirstChild
}

// AddChild implements the tree.Noder interface.
func (tn *UintNode) AddChild(target tree.Noder) {
	if target == nil {
		return
	}

	tmp, ok := target.(*UintNode)
	if !ok {
		return
	}
	
	tmp.NextSibling = nil
	tmp.PrevSibling = nil

	last_child := tn.LastChild

	if last_child == nil {
		tn.FirstChild = tmp
	} else {
		last_child.NextSibling = tmp
		tmp.PrevSibling = last_child
	}

	tmp.Parent = tn
	tn.LastChild = tmp
}

// LinkChildren implements the tree.Noder interface.
func (tn *UintNode) LinkChildren(children []tree.Noder) {
	var valid_children []*UintNode

	for _, child := range children {
		if child == nil {
			continue
		}

		tmp, ok := child.(*UintNode)
		if !ok {
			continue
		}

		tmp.Parent = tn

		valid_children = append(valid_children, tmp)
	}
	if len(valid_children) == 0 {
		return
	}

	valid_children[0].PrevSibling = nil
	valid_children[len(valid_children)-1].NextSibling = nil

	if len(valid_children) == 1 {
		return
	}

	for i := 0; i < len(valid_children)-1; i++ {
		valid_children[i].NextSibling = valid_children[i+1]
	}

	for i := 1; i < len(valid_children); i++ {
		valid_children[i].PrevSibling = valid_children[i-1]
	}

	tn.FirstChild, tn.LastChild = valid_children[0], valid_children[len(valid_children)-1]
}

// delete_child is a helper function to delete the child from the children of the node. No nil
// nodes are returned when this function is called. However, if target is nil, then nothing happens.
//
// Parameters:
//   - target: The child to remove.
//
// Returns:
//   - []UintNode: A slice of pointers to the children of the node.
func (tn *UintNode) delete_child(target *UintNode) []*UintNode {
	ok := tn.HasChild(target)
	if !ok {
		return nil
	}

	prev := target.PrevSibling
	next := target.NextSibling

	if prev != nil {
		prev.NextSibling = next
	}

	if next != nil {
		next.PrevSibling = prev
	}

	if target == tn.FirstChild {
		tn.FirstChild = next

		if next == nil {
			tn.LastChild = nil
		}
	} else if target == tn.LastChild {
		tn.LastChild = prev
	}

	target.Parent = nil
	target.PrevSibling = nil
	target.NextSibling = nil

	children := target.GetChildren()

	return children
}

// RemoveNode implements the tree.Noder interface.
func (tn *UintNode) RemoveNode() []tree.Noder {
	prev := tn.PrevSibling
	next := tn.NextSibling
	parent := tn.Parent

	var sub_roots []*UintNode

	if parent == nil {
		for c := tn.FirstChild; c != nil; c = c.NextSibling {
			sub_roots = append(sub_roots, c)
		}
	} else {
		children := parent.delete_child(tn)

		for _, child := range children {
			child.Parent = parent
		}
	}

	if prev != nil {
		prev.NextSibling = next
	} else {
		parent.FirstChild = next
	}

	if next != nil {
		next.PrevSibling = prev
	} else {
		parent.Parent.LastChild = prev
	}

	tn.Parent = nil
	tn.PrevSibling = nil
	tn.NextSibling = nil

	if len(sub_roots) == 0 {
		return nil
	}

	for _, child := range sub_roots {
		child.PrevSibling = nil
		child.NextSibling = nil
		child.Parent = nil
	}

	tn.FirstChild = nil
	tn.LastChild = nil

	conv := make([]tree.Noder, 0, len(sub_roots))
	for _, child := range sub_roots {
		conv = append(conv, child)
	}

	return conv
}

// Copy implements the tree.Noder interface.
//
// Although this function never returns nil, it does not copy the parent nor the sibling pointers.
func (tn *UintNode) Copy() common.Copier {
	var child_copy []tree.Noder

	for c := tn.FirstChild; c != nil; c = c.NextSibling {
		child_copy = append(child_copy, c.Copy().(tree.Noder))
	}

	// Copy here the data of the node.

	tn_copy := &UintNode{
	 	// Add here the copied data of the node.
	}

	tn_copy.LinkChildren(child_copy)

	return tn_copy
}

// Iterator implements the tree.Noder interface.
//
// This function returns an iterator that iterates over the direct children of the node.
// Implemented as a pull-based iterator, this function never returns nil and any of the
// values is guaranteed to be a non-nil node of type UintNode.
func (tn *UintNode) Iterator() common.Iterater[tree.Noder] {
	return &UintNodeIterator{
		parent: tn,
		current: tn.FirstChild,
	}
}

// Cleanup implements the tree.Noder interface.
//
// This function is expensive as it has to traverse the whole tree to clean up the nodes, one
// by one. While this function is useful for freeing up memory, for large enough trees, it is
// recommended to let the garbage collector handle the cleanup.
//
// Despite the above, this function does not use recursion but it is not safe to use in goroutines
// as pointers may be dereferenced while another goroutine is still using them.
//
// Finally, this function also logically removes the node from the siblings and the parent.
func (tn *UintNode) Cleanup() {
	type Helper struct {
		previous, current *UintNode
	}

	lls := stack.NewLinkedStack[*Helper]()

	// Free the first node.
	for c := tn.FirstChild; c != nil; c = c.NextSibling {
		h := &Helper{
			previous:	c.PrevSibling,
			current: 	c,
		}

		lls.Push(h)
	}

	tn.FirstChild = nil
	tn.LastChild = nil
	tn.Parent = nil

	// Free the rest of the nodes.
	for {
		h, ok := lls.Pop()
		if !ok {
			break
		}

		for c := h.current.FirstChild; c != nil; c = c.NextSibling {
			h := &Helper{
				previous:	c.PrevSibling,
				current: 	c,
			}

			lls.Push(h)
		}

		h.previous.NextSibling = nil
		h.previous.PrevSibling = nil

		h.current.FirstChild = nil
		h.current.LastChild = nil
		h.current.Parent = nil
	}

	prev := tn.PrevSibling
	next := tn.NextSibling

	if prev != nil {
		prev.NextSibling = next
	}

	if next != nil {
		next.PrevSibling = prev
	}

	tn.PrevSibling = nil
	tn.NextSibling = nil
}

// String implements the tree.Noder interface.
func (tn *UintNode) String() string {
	// WARNING: Implement this function.
	str := common.StringOf(tn.Data)

	return str
}

// NewUintNode creates a new node with the given data.
//
// Parameters:
//   - Data: The Data of the node.
//
// Returns:
//   - *UintNode: A pointer to the newly created node. It is
//   never nil.
func NewUintNode(data uint) *UintNode {
	return &UintNode{
		Data: data,
	}
}

// GetLastSibling returns the last sibling of the node. If it has a parent,
// it returns the last child of the parent. Otherwise, it returns the last
// sibling of the node.
//
// As an edge case, if the node has no parent and no next sibling, it returns
// the node itself. Thus, this function never returns nil.
//
// Returns:
//   - *UintNode: A pointer to the last sibling.
func (tn *UintNode) GetLastSibling() *UintNode {
	if tn.Parent != nil {
		return tn.Parent.LastChild
	} else if tn.NextSibling == nil {
		return tn
	}

	last_sibling := tn

	for last_sibling.NextSibling != nil {
		last_sibling = last_sibling.NextSibling
	}

	return last_sibling
}

// GetFirstSibling returns the first sibling of the node. If it has a parent,
// it returns the first child of the parent. Otherwise, it returns the first
// sibling of the node.
//
// As an edge case, if the node has no parent and no previous sibling, it returns
// the node itself. Thus, this function never returns nil.
//
// Returns:
//   - *UintNode: A pointer to the first sibling.
func (tn *UintNode) GetFirstSibling() *UintNode {
	if tn.Parent != nil {
		return tn.Parent.FirstChild
	} else if tn.PrevSibling == nil {
		return tn
	}

	first_sibling := tn

	for first_sibling.PrevSibling != nil {
		first_sibling = first_sibling.PrevSibling
	}

	return first_sibling
}

// AddChildren is a convenience function to add multiple children to the node at once.
// It is more efficient than adding them one by one. Therefore, the behaviors are the
// same as the behaviors of the UintNode.AddChild function.
//
// Parameters:
//   - children: The children to add.
func (tn *UintNode) AddChildren(children []*UintNode) {
	if len(children) == 0 {
		return
	}
	
	var top int

	for i := 0; i < len(children); i++ {
		child := children[i]

		if child != nil {
			children[top] = child
			top++
		}
	}

	children = children[:top]
	if len(children) == 0 {
		return
	}

	// Deal with the first child
	first_child := children[0]

	first_child.NextSibling = nil
	first_child.PrevSibling = nil

	last_child := tn.LastChild

	if last_child == nil {
		tn.FirstChild = first_child
	} else {
		last_child.NextSibling = first_child
		first_child.PrevSibling = last_child
	}

	first_child.Parent = tn
	tn.LastChild = first_child

	// Deal with the rest of the children
	for i := 1; i < len(children); i++ {
		child := children[i]

		child.NextSibling = nil
		child.PrevSibling = nil

		last_child := tn.LastChild
		last_child.NextSibling = child
		child.PrevSibling = last_child

		child.Parent = tn
		tn.LastChild = child
	}
}

// GetChildren returns the immediate children of the node.
//
// The returned nodes are never nil and are not copied. Thus, modifying the returned
// nodes will modify the tree.
//
// Returns:
//   - []*UintNode: A slice of pointers to the children of the node.
func (tn *UintNode) GetChildren() []*UintNode {
	var children []*UintNode

	for c := tn.FirstChild; c != nil; c = c.NextSibling {
		children = append(children, c)
	}

	return children
}

// HasChild returns true if the node has the given child.
//
// Because children of a node cannot be nil, a nil target will always return false.
//
// Parameters:
//   - target: The child to check for.
//
// Returns:
//   - bool: True if the node has the child, false otherwise.
func (tn *UintNode) HasChild(target *UintNode) bool {
	if target == nil || tn.FirstChild == nil {
		return false
	}

	for c := tn.FirstChild; c != nil; c = c.NextSibling {
		if c == target {
			return true
		}
	}

	return false
}

// IsChildOf returns true if the node is a child of the parent. If target is nil,
// it returns false.
//
// Parameters:
//   - target: The target parent to check for.
//
// Returns:
//   - bool: True if the node is a child of the parent, false otherwise.
func (tn *UintNode) IsChildOf(target *UintNode) bool {
	if target == nil {
		return false
	}

	parents := tree.GetNodeAncestors(target)

	for node := tn; node.Parent != nil; node = node.Parent {
		ok := slices.Contains(parents, node.Parent)
		if ok {
			return true
		}
	}

	return false
}
