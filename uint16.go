// Code generated by go generate; EDIT THIS FILE DIRECTLY

package treenode

import (
	"slices"
	"fmt"

	lls "github.com/PlayerR9/MyGoLib/ListLike/Stacker"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

// Uint16Iterator is a pull-based iterator that iterates
// over the children of a Uint16.
type Uint16Iterator struct {
	parent, current *Uint16
}

// Consume implements the common.Iterater interface.
//
// *common.ErrExhaustedIter is the only error returned by this function and the returned
// node is never nil.
func (iter *Uint16Iterator) Consume() (Noder, error) {
	if iter.current == nil {
		return nil, uc.NewErrExhaustedIter()
	}

	node := iter.current
	iter.current = iter.current.NextSibling

	return node, nil
}

// Restart implements the common.Iterater interface.
func (iter *Uint16Iterator) Restart() {
	iter.current = iter.parent.FirstChild
}

// Uint16 is a node in a tree.
type Uint16 struct {
	Parent, FirstChild, NextSibling, LastChild, PrevSibling *Uint16
	Data uint16
}

// Iterator implements the Tree.Noder interface.
//
// This function iterates over the children of the node, it is a pull-based iterator,
// and never returns nil.
func (u *Uint16) Iterator() uc.Iterater[Noder] {
	return &Uint16Iterator{
		parent: u,
		current: u.FirstChild,
	}
}

// String implements the Noder interface.
func (u *Uint16) String() string {
	// WARNING: Implement this function.
	str := fmt.Sprintf("%v", u.Data)

	return str
}

// Copy implements the Noder interface.
//
// It never returns nil and it does not copy the parent or the sibling pointers.
func (u *Uint16) Copy() uc.Copier {
	var child_copy []Noder	

	for c := u.FirstChild; c != nil; c = c.NextSibling {
		child_copy = append(child_copy, c.Copy().(Noder))
	}

	// Copy here the data of the node.

	u_copy := &Uint16{
	 	// Add here the copied data of the node.
	}

	u_copy.LinkChildren(child_copy)

	return u_copy
}

// SetParent implements the Noder interface.
func (u *Uint16) SetParent(parent Noder) bool {
	if parent == nil {
		u.Parent = nil
		return true
	}

	p, ok := parent.(*Uint16)
	if !ok {
		return false
	}

	u.Parent = p

	return true
}

// GetParent implements the Noder interface.
func (u *Uint16) GetParent() Noder {
	return u.Parent
}

// LinkWithParent implements the Noder interface.
//
// Children that are not of type *Uint16 or nil are ignored.
func (u *Uint16) LinkChildren(children []Noder) {
	if len(children) == 0 {
		return
	}

	var valid_children []*Uint16

	for _, child := range children {
		if child == nil {
			continue
		}

		c, ok := child.(*Uint16)
		if ok {
			c.Parent = u
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

	u.FirstChild, u.LastChild = valid_children[0], valid_children[len(valid_children)-1]
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
func (u *Uint16) GetLeaves() []Noder {
	// It is safe to change the stack implementation as long as
	// it is not limited in size. If it is, make sure to check the error
	// returned by the Push and Pop methods.
	stack := lls.NewLinkedStack[Noder](u)

	var leaves []Noder

	for {
		top, ok := stack.Pop()
		if !ok {
			break
		}

		node := top.(*Uint16)
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
func (u *Uint16) Cleanup() {
	type Helper struct {
		previous, current *Uint16
	}

	stack := lls.NewLinkedStack[*Helper]()

	// Free the first node.
	for c := u.FirstChild; c != nil; c = c.NextSibling {
		h := &Helper{
			previous:	c.PrevSibling,
			current: 	c,
		}

		stack.Push(h)
	}

	u.FirstChild = nil
	u.LastChild = nil
	u.Parent = nil

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

	prev := u.PrevSibling
	next := u.NextSibling

	if prev != nil {
		prev.NextSibling = next
	}

	if next != nil {
		next.PrevSibling = prev
	}

	u.PrevSibling = nil
	u.NextSibling = nil
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
func (u *Uint16) GetAncestors() []Noder {
	var ancestors []Noder

	for node := u; node.Parent != nil; node = node.Parent {
		ancestors = append(ancestors, node.Parent)
	}

	slices.Reverse(ancestors)

	return ancestors
}

// IsLeaf implements the Noder interface.
func (u *Uint16) IsLeaf() bool {
	return u.FirstChild == nil
}

// IsSingleton implements the Noder interface.
func (u *Uint16) IsSingleton() bool {
	return u.FirstChild != nil && u.FirstChild == u.LastChild
}

// GetFirstChild implements the Noder interface.
func (u *Uint16) GetFirstChild() Noder {
	return u.FirstChild
}

// DeleteChild implements the Noder interface.
//
// No nil nodes are returned.
func (u *Uint16) DeleteChild(target Noder) []Noder {
	if target == nil {
		return nil
	}

	n, ok := target.(*Uint16)
	if !ok {
		return nil
	}

	children := u.delete_child(n)

	if len(children) == 0 {
		return children
	}

	for _, child := range children {
		c := child.(*Uint16)

		c.PrevSibling = nil
		c.NextSibling = nil
		c.Parent = nil
	}

	u.FirstChild = nil
	u.LastChild = nil

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
func (u *Uint16) Size() int {
	// It is safe to change the stack implementation as long as
	// it is not limited in size. If it is, make sure to check the error
	// returned by the Push and Pop methods.
	stack := lls.NewLinkedStack(u)

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
// *Uint16, it does nothing.
//
// This function clears the parent and sibling pointers of the child and so, it
// does not add relatives to the child.
//
// Parameters:
//   - child: The child to add.
func (u *Uint16) AddChild(child Noder) {
	if child == nil {
		return
	}

	c, ok := child.(*Uint16)
	if !ok {
		return
	}
	
	c.NextSibling = nil
	c.PrevSibling = nil

	last_child := u.LastChild

	if last_child == nil {
		u.FirstChild = c
	} else {
		last_child.NextSibling = c
		c.PrevSibling = last_child
	}

	c.Parent = u
	u.LastChild = c
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
func (u *Uint16) RemoveNode() []Noder {
	prev := u.PrevSibling
	next := u.NextSibling
	parent := u.Parent

	var sub_roots []Noder

	if parent == nil {
		for c := u.FirstChild; c != nil; c = c.NextSibling {
			sub_roots = append(sub_roots, c)
		}
	} else {
		children := parent.delete_child(u)

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

	u.Parent = nil
	u.PrevSibling = nil
	u.NextSibling = nil

	if len(sub_roots) == 0 {
		return sub_roots
	}

	for _, child := range sub_roots {
		c := child.(*Uint16)

		c.PrevSibling = nil
		c.NextSibling = nil
		c.Parent = nil
	}

	u.FirstChild = nil
	u.LastChild = nil

	return sub_roots
}

// NewUint16 creates a new node with the given data.
//
// Parameters:
//   - Data: The Data of the node.
//
// Returns:
//   - *Uint16: A pointer to the newly created node. It is
//   never nil.
func NewUint16(data uint16) *Uint16 {
	return &Uint16{
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
//   - *Uint16: A pointer to the last sibling.
func (u *Uint16) GetLastSibling() *Uint16 {
	if u.Parent != nil {
		return u.Parent.LastChild
	} else if u.NextSibling == nil {
		return u
	}

	last_sibling := u

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
//   - *Uint16: A pointer to the first sibling.
func (u *Uint16) GetFirstSibling() *Uint16 {
	if u.Parent != nil {
		return u.Parent.FirstChild
	} else if u.PrevSibling == nil {
		return u
	}

	first_sibling := u

	for first_sibling.PrevSibling != nil {
		first_sibling = first_sibling.PrevSibling
	}

	return first_sibling
}

// IsRoot returns true if the node does not have a parent.
//
// Returns:
//   - bool: True if the node is the root, false otherwise.
func (u *Uint16) IsRoot() bool {
	return u.Parent == nil
}

// AddChildren is a convenience function to add multiple children to the node at once.
// It is more efficient than adding them one by one. Therefore, the behaviors are the
// same as the behaviors of the Uint16.AddChild function.
//
// Parameters:
//   - children: The children to add.
func (u *Uint16) AddChildren(children []*Uint16) {
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

	last_child := u.LastChild

	if last_child == nil {
		u.FirstChild = first_child
	} else {
		last_child.NextSibling = first_child
		first_child.PrevSibling = last_child
	}

	first_child.Parent = u
	u.LastChild = first_child

	// Deal with the rest of the children
	for i := 1; i < len(children); i++ {
		child := children[i]

		child.NextSibling = nil
		child.PrevSibling = nil

		last_child := u.LastChild
		last_child.NextSibling = child
		child.PrevSibling = last_child

		child.Parent = u
		u.LastChild = child
	}
}

// GetChildren returns the immediate children of the node.
//
// The returned nodes are never nil and are not copied. Thus, modifying the returned
// nodes will modify the tree.
//
// Returns:
//   - []Noder: A slice of pointers to the children of the node.
func (u *Uint16) GetChildren() []Noder {
	var children []Noder

	for c := u.FirstChild; c != nil; c = c.NextSibling {
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
func (u *Uint16) HasChild(target *Uint16) bool {
	if target == nil || u.FirstChild == nil {
		return false
	}

	for c := u.FirstChild; c != nil; c = c.NextSibling {
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
func (u *Uint16) delete_child(target *Uint16) []Noder {
	ok := u.HasChild(target)
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

	if target == u.FirstChild {
		u.FirstChild = next

		if next == nil {
			u.LastChild = nil
		}
	} else if target == u.LastChild {
		u.LastChild = prev
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
func (u *Uint16) IsChildOf(target *Uint16) bool {
	if target == nil {
		return false
	}

	parents := target.GetAncestors()

	for node := u; node.Parent != nil; node = node.Parent {
		parent := Noder(node.Parent)

		ok := slices.Contains(parents, parent)
		if ok {
			return true
		}
	}

	return false
}
