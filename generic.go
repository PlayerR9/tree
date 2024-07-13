// Code generated by go generate; EDIT THIS FILE DIRECTLY

package treenode

import (
	"slices"
	"fmt"

	"github.com/PlayerR9/MyGoLib/ListLike/Stacker"
	"github.com/PlayerR9/MyGoLib/Units/common"
)

// TreeNodeIterator is a pull-based iterator that iterates
// over the children of a TreeNode.
type TreeNodeIterator[T any] struct {
	parent, current *TreeNode[T]
}

// Consume implements the common.Iterater interface.
//
// *common.ErrExhaustedIter is the only error returned by this function and the returned
// node is never nil.
func (iter *TreeNodeIterator[T]) Consume() (Noder, error) {
	if iter.current == nil {
		return nil, common.NewErrExhaustedIter()
	}

	node := iter.current
	iter.current = iter.current.NextSibling

	return node, nil
}

// Restart implements the common.Iterater interface.
func (iter *TreeNodeIterator[T]) Restart() {
	iter.current = iter.parent.FirstChild
}

// TreeNode is a node in a tree.
type TreeNode[T any] struct {
	Parent, FirstChild, NextSibling, LastChild, PrevSibling *TreeNode[T]
	Data T
}

// Iterator implements the Noder interface.
//
// This function iterates over the children of the node, it is a pull-based iterator,
// and never returns nil.
func (tn *TreeNode[T]) Iterator() common.Iterater[Noder] {
	return &TreeNodeIterator[T]{
		parent: tn,
		current: tn.FirstChild,
	}
}

// String implements the Noder interface.
func (tn *TreeNode[T]) String() string {
	// WARNING: Implement this function.
	str := fmt.Sprintf("%v", tn.Data)

	return str
}

// Copy implements the Noder interface.
//
// It never returns nil and it does not copy the parent or the sibling pointers.
func (tn *TreeNode[T]) Copy() common.Copier {
	var child_copy []Noder	

	for c := tn.FirstChild; c != nil; c = c.NextSibling {
		child_copy = append(child_copy, c.Copy().(Noder))
	}

	// Copy here the data of the node.

	tn_copy := &TreeNode[T]{
	 	// Add here the copied data of the node.
	}

	tn_copy.LinkChildren(child_copy)

	return tn_copy
}

// SetParent implements the Noder interface.
func (tn *TreeNode[T]) SetParent(parent Noder) bool {
	if parent == nil {
		tn.Parent = nil
		return true
	}

	p, ok := parent.(*TreeNode[T])
	if !ok {
		return false
	}

	tn.Parent = p

	return true
}

// GetParent implements the Noder interface.
func (tn *TreeNode[T]) GetParent() Noder {
	return tn.Parent
}

// LinkWithParent implements the Noder interface.
//
// Children that are not of type *TreeNode[T] or nil are ignored.
func (tn *TreeNode[T]) LinkChildren(children []Noder) {
	if len(children) == 0 {
		return
	}

	var valid_children []*TreeNode[T]

	for _, child := range children {
		if child == nil {
			continue
		}

		c, ok := child.(*TreeNode[T])
		if ok {
			c.Parent = tn
			valid_children = append(valid_children, c)
		}		
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

// GetLeaves implements the Noder interface.
//
// This is expensive as leaves are not stored and so, every time this function is called,
// it has to do a DFS traversal to find the leaves. Thus, it is recommended to call
// this function once and then store the leaves somewhere if needed.
//
// Despite the above, this function does not use recursion and is safe to use.
//
// Finally, no nil nodes are returned.
func (tn *TreeNode[T]) GetLeaves() []Noder {
	// It is safe to change the stack implementation as long as
	// it is not limited in size. If it is, make sure to check the error
	// returned by the Push and Pop methods.
	stack := Stacker.NewLinkedStack[Noder](tn)

	var leaves []Noder

	for {
		top, ok := stack.Pop()
		if !ok {
			break
		}

		node := top.(*TreeNode[T])
		if node.FirstChild == nil {
			leaves = append(leaves, top)
		} else {
			for c := node.FirstChild; c != nil; c = c.NextSibling {
				stack.Push(c)
			}
		}
	}

	return leaves
}

// Cleanup implements the Noder interface.
//
// This is expensive as it has to traverse the whole tree to clean up the nodes, one
// by one. While this is useful for freeing up memory, for large enough trees, it is
// recommended to let the garbage collector handle the cleanup.
//
// Despite the above, this function does not use recursion and is safe to use (but
// make sure goroutines are not running on the tree while this function is called).
//
// Finally, it also logically removes the node from the siblings and the parent.
func (tn *TreeNode[T]) Cleanup() {
	type Helper struct {
		previous, current *TreeNode[T]
	}

	stack := Stacker.NewLinkedStack[*Helper]()

	// Free the first node.
	for c := tn.FirstChild; c != nil; c = c.NextSibling {
		h := &Helper{
			previous:	c.PrevSibling,
			current: 	c,
		}

		stack.Push(h)
	}

	tn.FirstChild = nil
	tn.LastChild = nil
	tn.Parent = nil

	// Free the rest of the nodes.
	for {
		h, ok := stack.Pop()
		if !ok {
			break
		}

		for c := h.current.FirstChild; c != nil; c = c.NextSibling {
			h := &Helper{
				previous:	c.PrevSibling,
				current: 	c,
			}

			stack.Push(h)
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

// GetAncestors implements the Noder interface.
//
// This is expensive since ancestors are not stored and so, every time this
// function is called, it has to traverse the tree to find the ancestors. Thus, it is
// recommended to call this function once and then store the ancestors somewhere if needed.
//
// Despite the above, this function does not use recursion and is safe to use.
//
// Finally, no nil nodes are returned.
func (tn *TreeNode[T]) GetAncestors() []Noder {
	var ancestors []Noder

	for node := tn; node.Parent != nil; node = node.Parent {
		ancestors = append(ancestors, node.Parent)
	}

	slices.Reverse(ancestors)

	return ancestors
}

// IsLeaf implements the Noder interface.
func (tn *TreeNode[T]) IsLeaf() bool {
	return tn.FirstChild == nil
}

// IsSingleton implements the Noder interface.
func (tn *TreeNode[T]) IsSingleton() bool {
	return tn.FirstChild != nil && tn.FirstChild == tn.LastChild
}

// GetFirstChild implements the Noder interface.
func (tn *TreeNode[T]) GetFirstChild() Noder {
	return tn.FirstChild
}

// DeleteChild implements the Noder interface.
//
// No nil nodes are returned.
func (tn *TreeNode[T]) DeleteChild(target Noder) []Noder {
	if target == nil {
		return nil
	}

	n, ok := target.(*TreeNode[T])
	if !ok {
		return nil
	}

	children := tn.delete_child(n)

	if len(children) == 0 {
		return children
	}

	for _, child := range children {
		c := child.(*TreeNode[T])

		c.PrevSibling = nil
		c.NextSibling = nil
		c.Parent = nil
	}

	tn.FirstChild = nil
	tn.LastChild = nil

	return children
}

// Size implements the Noder interface.
//
// This is expensive as it has to traverse the whole tree to find the size of the tree.
// Thus, it is recommended to call this function once and then store the size somewhere if needed.
//
// Despite the above, this function does not use recursion and is safe to use.
//
// Finally, the traversal is done in a depth-first manner.
func (tn *TreeNode[T]) Size() int {
	// It is safe to change the stack implementation as long as
	// it is not limited in size. If it is, make sure to check the error
	// returned by the Push and Pop methods.
	stack := Stacker.NewLinkedStack(tn)

	var size int

	for {
		top, ok := stack.Pop()
		if !ok {
			break
		}

		size++

		for c := top.FirstChild; c != nil; c = c.NextSibling {
			stack.Push(c)
		}
	}

	return size
}

// AddChild adds a new child to the node. If the child is nil or it is not of type
// *TreeNode[T], it does nothing.
//
// This function clears the parent and sibling pointers of the child and so, it
// does not add relatives to the child.
//
// Parameters:
//   - child: The child to add.
func (tn *TreeNode[T]) AddChild(child Noder) {
	if child == nil {
		return
	}

	c, ok := child.(*TreeNode[T])
	if !ok {
		return
	}
	
	c.NextSibling = nil
	c.PrevSibling = nil

	last_child := tn.LastChild

	if last_child == nil {
		tn.FirstChild = c
	} else {
		last_child.NextSibling = c
		c.PrevSibling = last_child
	}

	c.Parent = tn
	tn.LastChild = c
}

// RemoveNode removes the node from the tree while shifting the children up one level to
// maintain the tree structure.
//
// Also, the returned children can be used to create a forest of trees if the root node
// is removed.
//
// Returns:
//   - []Noder: A slice of pointers to the children of the node iff the node is the root.
//     Nil otherwise.
//
// Example:
//
//	// Given the tree:
//	1
//	├── 2
//	└── 3
//		├── 4
//		└── 5
//	└── 6
//
//	// The tree after removing node 3:
//
//	1
//	├── 2
//	└── 4
//	└── 5
//	└── 6
func (tn *TreeNode[T]) RemoveNode() []Noder {
	prev := tn.PrevSibling
	next := tn.NextSibling
	parent := tn.Parent

	var sub_roots []Noder

	if parent == nil {
		for c := tn.FirstChild; c != nil; c = c.NextSibling {
			sub_roots = append(sub_roots, c)
		}
	} else {
		children := parent.delete_child(tn)

		for _, child := range children {
			child.SetParent(parent)
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
		return sub_roots
	}

	for _, child := range sub_roots {
		c := child.(*TreeNode[T])

		c.PrevSibling = nil
		c.NextSibling = nil
		c.Parent = nil
	}

	tn.FirstChild = nil
	tn.LastChild = nil

	return sub_roots
}

// NewTreeNode creates a new node with the given data.
//
// Parameters:
//   - Data: The Data of the node.
//
// Returns:
//   - *TreeNode[T]: A pointer to the newly created node. It is
//   never nil.
func NewTreeNode[T any](data T) *TreeNode[T] {
	return &TreeNode[T]{
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
//   - *TreeNode[T]: A pointer to the last sibling.
func (tn *TreeNode[T]) GetLastSibling() *TreeNode[T] {
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
//   - *TreeNode[T]: A pointer to the first sibling.
func (tn *TreeNode[T]) GetFirstSibling() *TreeNode[T] {
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

// IsRoot returns true if the node does not have a parent.
//
// Returns:
//   - bool: True if the node is the root, false otherwise.
func (tn *TreeNode[T]) IsRoot() bool {
	return tn.Parent == nil
}

// AddChildren is a convenience function to add multiple children to the node at once.
// It is more efficient than adding them one by one. Therefore, the behaviors are the
// same as the behaviors of the TreeNode.AddChild function.
//
// Parameters:
//   - children: The children to add.
func (tn *TreeNode[T]) AddChildren(children []*TreeNode[T]) {
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
//   - []Noder: A slice of pointers to the children of the node.
func (tn *TreeNode[T]) GetChildren() []Noder {
	var children []Noder

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
func (tn *TreeNode[T]) HasChild(target *TreeNode[T]) bool {
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

// delete_child is a helper function to delete the child from the children of the node.
//
// No nil nodes are returned.
//
// Parameters:
//   - target: The child to remove.
//
// Returns:
//   - []Noder: A slice of pointers to the children of the node.
func (tn *TreeNode[T]) delete_child(target *TreeNode[T]) []Noder {
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

// IsChildOf returns true if the node is a child of the parent. If target is nil,
// it returns false.
//
// Parameters:
//   - target: The target parent to check for.
//
// Returns:
//   - bool: True if the node is a child of the parent, false otherwise.
func (tn *TreeNode[T]) IsChildOf(target *TreeNode[T]) bool {
	if target == nil {
		return false
	}

	parents := target.GetAncestors()

	for node := tn; node.Parent != nil; node = node.Parent {
		parent := Noder(node.Parent)

		ok := slices.Contains(parents, parent)
		if ok {
			return true
		}
	}

	return false
}
