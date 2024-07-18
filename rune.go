// Code generated by go generate; EDIT THIS FILE DIRECTLY

package treenode

import (
	"slices"

	"github.com/PlayerR9/MyGoLib/ListLike/Stacker"
	"github.com/PlayerR9/MyGoLib/Units/common"
)

// RuneNodeIterator is a pull-based iterator that iterates
// over the children of a RuneNode.
type RuneNodeIterator struct {
	parent, current *RuneNode
}

// Consume implements the common.Iterater interface.
//
// *common.ErrExhaustedIter is the only error returned by this function and the returned
// node is never nil.
func (iter *RuneNodeIterator) Consume() (*RuneNode, error) {
	if iter.current == nil {
		return nil, common.NewErrExhaustedIter()
	}

	node := iter.current
	iter.current = iter.current.NextSibling

	return node, nil
}

// Restart implements the common.Iterater interface.
func (iter *RuneNodeIterator) Restart() {
	iter.current = iter.parent.FirstChild
}

// RuneNode is a node in a tree.
type RuneNode struct {
	Parent, FirstChild, NextSibling, LastChild, PrevSibling *RuneNode
	Data rune
}

// Iterator implements the RuneNode interface.
//
// This function iterates over the children of the node, it is a pull-based iterator,
// and never returns nil.
func (tn *RuneNode) Iterator() common.Iterater[*RuneNode] {
	return &RuneNodeIterator{
		parent: tn,
		current: tn.FirstChild,
	}
}

// String implements the RuneNode interface.
func (tn *RuneNode) String() string {
	// WARNING: Implement this function.
	str := common.StringOf(tn.Data)

	return str
}

// Copy implements the RuneNode interface.
//
// It never returns nil and it does not copy the parent or the sibling pointers.
func (tn *RuneNode) Copy() common.Copier {
	var child_copy []*RuneNode	

	for c := tn.FirstChild; c != nil; c = c.NextSibling {
		child_copy = append(child_copy, c.Copy().(*RuneNode))
	}

	// Copy here the data of the node.

	tn_copy := &RuneNode{
	 	// Add here the copied data of the node.
	}

	tn_copy.LinkChildren(child_copy)

	return tn_copy
}

// Cleanup implements the RuneNode interface.
//
// This is expensive as it has to traverse the whole tree to clean up the nodes, one
// by one. While this is useful for freeing up memory, for large enough trees, it is
// recommended to let the garbage collector handle the cleanup.
//
// Despite the above, this function does not use recursion and is safe to use (but
// make sure goroutines are not running on the tree while this function is called).
//
// Finally, it also logically removes the node from the siblings and the parent.
func (tn *RuneNode) Cleanup() {
	type Helper struct {
		previous, current *RuneNode
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

// IsLeaf implements the RuneNode interface.
func (tn *RuneNode) IsLeaf() bool {
	return tn.FirstChild == nil
}

// IsSingleton implements the RuneNode interface.
func (tn *RuneNode) IsSingleton() bool {
	return tn.FirstChild != nil && tn.FirstChild == tn.LastChild
}

// DeleteChild removes the given child from the children of the node.
//
// Parameters:
//   - target: The child to remove.
//
// Returns:
//   - []*RuneNode: A slice of pointers to the children of the node. Nil if the node has no children.
//
// No nil nodes are returned.
func (tn *RuneNode) DeleteChild(target *RuneNode) []*RuneNode {
	if target == nil {
		return nil
	}

	children := tn.delete_child(target)

	if len(children) == 0 {
		return children
	}

	for _, child := range children {
		child.PrevSibling = nil
		child.NextSibling = nil
		child.Parent = nil
	}

	tn.FirstChild = nil
	tn.LastChild = nil

	return children
}

// Size implements the RuneNode interface.
//
// This is expensive as it has to traverse the whole tree to find the size of the tree.
// Thus, it is recommended to call this function once and then store the size somewhere if needed.
//
// Despite the above, this function does not use recursion and is safe to use.
//
// Finally, the traversal is done in a depth-first manner.
func (tn *RuneNode) Size() int {
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

// NewRuneNode creates a new node with the given data.
//
// Parameters:
//   - Data: The Data of the node.
//
// Returns:
//   - *RuneNode: A pointer to the newly created node. It is
//   never nil.
func NewRuneNode(data rune) *RuneNode {
	return &RuneNode{
		Data: data,
	}
}

// LinkChildren links the parent with the children. It also links the children
// with each other. Nil children are ignored.
//
// Parameters:
//   - children: The children nodes.
func (tn *RuneNode) LinkChildren(children []*RuneNode) {
	if len(children) == 0 {
		return
	}

	var valid_children []*RuneNode

	for _, child := range children {
		if child == nil {
			continue
		}

		child.Parent = tn
		valid_children = append(valid_children, child)		
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

// AddChild adds a new child to the node. If the child is nil it does nothing.
//
// Parameters:
//   - child: The child to add.
//
// This function clears the parent and sibling pointers of the child and so, it
// does not add relatives to the child.
func (tn *RuneNode) AddChild(child *RuneNode) {
	if child == nil {
		return
	}
	
	child.NextSibling = nil
	child.PrevSibling = nil

	last_child := tn.LastChild

	if last_child == nil {
		tn.FirstChild = child
	} else {
		last_child.NextSibling = child
		child.PrevSibling = last_child
	}

	child.Parent = tn
	tn.LastChild = child
}

// RemoveNode removes the node from the tree while shifting the children up one level to
// maintain the tree structure.
//
// Also, the returned children can be used to create a forest of trees if the root node
// is removed.
//
// Returns:
//   - []*RuneNode: A slice of pointers to the children of the node iff the node is the root.
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
func (tn *RuneNode) RemoveNode() []*RuneNode {
	prev := tn.PrevSibling
	next := tn.NextSibling
	parent := tn.Parent

	var sub_roots []*RuneNode

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
		return sub_roots
	}

	for _, child := range sub_roots {
		child.PrevSibling = nil
		child.NextSibling = nil
		child.Parent = nil
	}

	tn.FirstChild = nil
	tn.LastChild = nil

	return sub_roots
}

// GetLeaves returns all the leaves of the tree rooted at the node.
//
// Returns:
//   - []*RuneNode: A slice of pointers to the leaves of the tree.
//
// This is expensive as leaves are not stored and so, every time this function is called,
// it has to do a DFS traversal to find the leaves. Thus, it is recommended to call
// this function once and then store the leaves somewhere if needed.
//
// Despite the above, this function does not use recursion and is safe to use.
//
// Finally, no nil nodes are returned.
func (tn *RuneNode) GetLeaves() []*RuneNode {
	// It is safe to change the stack implementation as long as
	// it is not limited in size. If it is, make sure to check the error
	// returned by the Push and Pop methods.
	stack := Stacker.NewLinkedStack(tn)

	var leaves []*RuneNode

	for {
		top, ok := stack.Pop()
		if !ok {
			break
		}

		if top.FirstChild == nil {
			leaves = append(leaves, top)
		} else {
			for c := top.FirstChild; c != nil; c = c.NextSibling {
				stack.Push(c)
			}
		}
	}

	return leaves
}

// GetAncestors returns all the ancestors of the node. This does not return the node itself.
//
// Returns:
//   - []*RuneNode: A slice of pointers to the ancestors of the node.
//
// The ancestors are returned in the opposite order of a DFS traversal. Therefore, the first element is the parent
// of the node.
//
// This is expensive since ancestors are not stored and so, every time this
// function is called, it has to traverse the tree to find the ancestors. Thus, it is
// recommended to call this function once and then store the ancestors somewhere if needed.
//
// Despite the above, this function does not use recursion and is safe to use.
//
// Finally, no nil nodes are returned.
func (tn *RuneNode) GetAncestors() []*RuneNode {
	var ancestors []*RuneNode

	for node := tn; node.Parent != nil; node = node.Parent {
		ancestors = append(ancestors, node.Parent)
	}

	slices.Reverse(ancestors)

	return ancestors
}

// GetLastSibling returns the last sibling of the node. If it has a parent,
// it returns the last child of the parent. Otherwise, it returns the last
// sibling of the node.
//
// As an edge case, if the node has no parent and no next sibling, it returns
// the node itself. Thus, this function never returns nil.
//
// Returns:
//   - *RuneNode: A pointer to the last sibling.
func (tn *RuneNode) GetLastSibling() *RuneNode {
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
//   - *RuneNode: A pointer to the first sibling.
func (tn *RuneNode) GetFirstSibling() *RuneNode {
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
func (tn *RuneNode) IsRoot() bool {
	return tn.Parent == nil
}

// AddChildren is a convenience function to add multiple children to the node at once.
// It is more efficient than adding them one by one. Therefore, the behaviors are the
// same as the behaviors of the RuneNode.AddChild function.
//
// Parameters:
//   - children: The children to add.
func (tn *RuneNode) AddChildren(children []*RuneNode) {
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
//   - []*RuneNode: A slice of pointers to the children of the node.
func (tn *RuneNode) GetChildren() []*RuneNode {
	var children []*RuneNode

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
func (tn *RuneNode) HasChild(target *RuneNode) bool {
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
//   - []RuneNode: A slice of pointers to the children of the node.
func (tn *RuneNode) delete_child(target *RuneNode) []*RuneNode {
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
func (tn *RuneNode) IsChildOf(target *RuneNode) bool {
	if target == nil {
		return false
	}

	parents := target.GetAncestors()

	for node := tn; node.Parent != nil; node = node.Parent {
		ok := slices.Contains(parents, node.Parent)
		if ok {
			return true
		}
	}

	return false
}

/*

// FindCommonAncestor returns the first common ancestor of the two nodes.
//
// Parameters:
//   - n1: The first node.
//   - n2: The second node.
//
// Returns:
//   - *TreeNode[T]: A pointer to the common ancestor. Nil if no such node is found.
func FindCommonAncestor[T any](n1, n2 *TreeNode[T]) *TreeNode[T] {
	if n1 == nil {
		return n2
	} else if n2 == nil {
		return n1
	} else if n1 == n2 {
		return n1
	}

	ancestors1 := n1.GetAncestors()
	ancestors2 := n2.GetAncestors()

	if len(ancestors1) > len(ancestors2) {
		ancestors1, ancestors2 = ancestors2, ancestors1
	}

	for _, node := range ancestors1 {
		ok := slices.Contains(ancestors2, node)
		if ok {
			return node
		}
	}

	return nil
}

// FindBranchingPoint returns the first node in the path from n to the root
// such that has more than one sibling.
//
// Returns:
//   - *TreeNode[T]: The branching point.
//   - *TreeNode[T]: The parent of the branching point.
//   - bool: True if the node has a branching point, false otherwise.
//
// Behaviors:
//   - If there is no branching point, it returns the root of the tree. However,
//     if n is nil, it returns nil, nil, false and if the node has no parent, it
//     returns nil, n, false.
func FindBranchingPoint[T any](n *TreeNode[T]) (*TreeNode[T], *TreeNode[T], bool) {
	if n == nil {
		return nil, nil, false
	}

	parent := n.GetParent()
	if parent == nil {
		return nil, n, false
	}

	var has_branching_point bool

	for !has_branching_point {
		grand_parent := parent.GetParent()
		if grand_parent == nil {
			break
		}

		ok := parent.IsSingleton()
		if !ok {
			has_branching_point = true
		} else {
			n = parent
			parent = grand_parent
		}
	}

	return n, parent, has_branching_point
}
*/
