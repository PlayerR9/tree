// Code generated by go generate; EDIT THIS FILE DIRECTLY

package treenode

import (
	"slices"
	"fmt"

	lls "github.com/PlayerR9/MyGoLib/ListLike/Stacker"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

// Float64Iterator is a pull-based iterator that iterates
// over the children of a Float64.
type Float64Iterator struct {
	parent, current *Float64
}

// Consume implements the common.Iterater interface.
//
// *common.ErrExhaustedIter is the only error returned by this function and the returned
// node is never nil.
func (iter *Float64Iterator) Consume() (Noder, error) {
	if iter.current == nil {
		return nil, uc.NewErrExhaustedIter()
	}

	node := iter.current
	iter.current = iter.current.NextSibling

	return node, nil
}

// Restart implements the common.Iterater interface.
func (iter *Float64Iterator) Restart() {
	iter.current = iter.parent.FirstChild
}

// Float64 is a node in a tree.
type Float64 struct {
	Parent, FirstChild, NextSibling, LastChild, PrevSibling *Float64
	Data float64
}

// Iterator implements the Tree.Noder interface.
//
// This function iterates over the children of the node, it is a pull-based iterator,
// and never returns nil.
func (f *Float64) Iterator() uc.Iterater[Noder] {
	return &Float64Iterator{
		parent: f,
		current: f.FirstChild,
	}
}

// String implements the Noder interface.
func (f *Float64) String() string {
	// WARNING: Implement this function.
	str := fmt.Sprintf("%v", f.Data)

	return str
}

// Copy implements the Noder interface.
//
// It never returns nil and it does not copy the parent or the sibling pointers.
func (f *Float64) Copy() uc.Copier {
	var child_copy []Noder	

	for c := f.FirstChild; c != nil; c = c.NextSibling {
		child_copy = append(child_copy, c.Copy().(Noder))
	}

	// Copy here the data of the node.

	f_copy := &Float64{
	 	// Add here the copied data of the node.
	}

	f_copy.LinkChildren(child_copy)

	return f_copy
}

// SetParent implements the Noder interface.
func (f *Float64) SetParent(parent Noder) bool {
	if parent == nil {
		f.Parent = nil
		return true
	}

	p, ok := parent.(*Float64)
	if !ok {
		return false
	}

	f.Parent = p

	return true
}

// GetParent implements the Noder interface.
func (f *Float64) GetParent() Noder {
	return f.Parent
}

// LinkWithParent implements the Noder interface.
//
// Children that are not of type *Float64 or nil are ignored.
func (f *Float64) LinkChildren(children []Noder) {
	if len(children) == 0 {
		return
	}

	var valid_children []*Float64

	for _, child := range children {
		if child == nil {
			continue
		}

		c, ok := child.(*Float64)
		if ok {
			c.Parent = f
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

	f.FirstChild, f.LastChild = valid_children[0], valid_children[len(valid_children)-1]
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
func (f *Float64) GetLeaves() []Noder {
	// It is safe to change the stack implementation as long as
	// it is not limited in size. If it is, make sure to check the error
	// returned by the Push and Pop methods.
	stack := lls.NewLinkedStack[Noder](f)

	var leaves []Noder

	for {
		top, ok := stack.Pop()
		if !ok {
			break
		}

		node := top.(*Float64)
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
func (f *Float64) Cleanup() {
	type Helper struct {
		previous, current *Float64
	}

	stack := lls.NewLinkedStack[*Helper]()

	// Free the first node.
	for c := f.FirstChild; c != nil; c = c.NextSibling {
		h := &Helper{
			previous:	c.PrevSibling,
			current: 	c,
		}

		stack.Push(h)
	}

	f.FirstChild = nil
	f.LastChild = nil
	f.Parent = nil

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

	prev := f.PrevSibling
	next := f.NextSibling

	if prev != nil {
		prev.NextSibling = next
	}

	if next != nil {
		next.PrevSibling = prev
	}

	f.PrevSibling = nil
	f.NextSibling = nil
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
func (f *Float64) GetAncestors() []Noder {
	var ancestors []Noder

	for node := f; node.Parent != nil; node = node.Parent {
		ancestors = append(ancestors, node.Parent)
	}

	slices.Reverse(ancestors)

	return ancestors
}

// IsLeaf implements the Noder interface.
func (f *Float64) IsLeaf() bool {
	return f.FirstChild == nil
}

// IsSingleton implements the Noder interface.
func (f *Float64) IsSingleton() bool {
	return f.FirstChild != nil && f.FirstChild == f.LastChild
}

// GetFirstChild implements the Noder interface.
func (f *Float64) GetFirstChild() Noder {
	return f.FirstChild
}

// DeleteChild implements the Noder interface.
//
// No nil nodes are returned.
func (f *Float64) DeleteChild(target Noder) []Noder {
	if target == nil {
		return nil
	}

	n, ok := target.(*Float64)
	if !ok {
		return nil
	}

	children := f.delete_child(n)

	if len(children) == 0 {
		return children
	}

	for _, child := range children {
		c := child.(*Float64)

		c.PrevSibling = nil
		c.NextSibling = nil
		c.Parent = nil
	}

	f.FirstChild = nil
	f.LastChild = nil

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
func (f *Float64) Size() int {
	// It is safe to change the stack implementation as long as
	// it is not limited in size. If it is, make sure to check the error
	// returned by the Push and Pop methods.
	stack := lls.NewLinkedStack(f)

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
// *Float64, it does nothing.
//
// This function clears the parent and sibling pointers of the child and so, it
// does not add relatives to the child.
//
// Parameters:
//   - child: The child to add.
func (f *Float64) AddChild(child Noder) {
	if child == nil {
		return
	}

	c, ok := child.(*Float64)
	if !ok {
		return
	}
	
	c.NextSibling = nil
	c.PrevSibling = nil

	last_child := f.LastChild

	if last_child == nil {
		f.FirstChild = c
	} else {
		last_child.NextSibling = c
		c.PrevSibling = last_child
	}

	c.Parent = f
	f.LastChild = c
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
func (f *Float64) RemoveNode() []Noder {
	prev := f.PrevSibling
	next := f.NextSibling
	parent := f.Parent

	var sub_roots []Noder

	if parent == nil {
		for c := f.FirstChild; c != nil; c = c.NextSibling {
			sub_roots = append(sub_roots, c)
		}
	} else {
		children := parent.delete_child(f)

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

	f.Parent = nil
	f.PrevSibling = nil
	f.NextSibling = nil

	if len(sub_roots) == 0 {
		return sub_roots
	}

	for _, child := range sub_roots {
		c := child.(*Float64)

		c.PrevSibling = nil
		c.NextSibling = nil
		c.Parent = nil
	}

	f.FirstChild = nil
	f.LastChild = nil

	return sub_roots
}

// NewFloat64 creates a new node with the given data.
//
// Parameters:
//   - Data: The Data of the node.
//
// Returns:
//   - *Float64: A pointer to the newly created node. It is
//   never nil.
func NewFloat64(data float64) *Float64 {
	return &Float64{
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
//   - *Float64: A pointer to the last sibling.
func (f *Float64) GetLastSibling() *Float64 {
	if f.Parent != nil {
		return f.Parent.LastChild
	} else if f.NextSibling == nil {
		return f
	}

	last_sibling := f

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
//   - *Float64: A pointer to the first sibling.
func (f *Float64) GetFirstSibling() *Float64 {
	if f.Parent != nil {
		return f.Parent.FirstChild
	} else if f.PrevSibling == nil {
		return f
	}

	first_sibling := f

	for first_sibling.PrevSibling != nil {
		first_sibling = first_sibling.PrevSibling
	}

	return first_sibling
}

// IsRoot returns true if the node does not have a parent.
//
// Returns:
//   - bool: True if the node is the root, false otherwise.
func (f *Float64) IsRoot() bool {
	return f.Parent == nil
}

// AddChildren is a convenience function to add multiple children to the node at once.
// It is more efficient than adding them one by one. Therefore, the behaviors are the
// same as the behaviors of the Float64.AddChild function.
//
// Parameters:
//   - children: The children to add.
func (f *Float64) AddChildren(children []*Float64) {
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

	last_child := f.LastChild

	if last_child == nil {
		f.FirstChild = first_child
	} else {
		last_child.NextSibling = first_child
		first_child.PrevSibling = last_child
	}

	first_child.Parent = f
	f.LastChild = first_child

	// Deal with the rest of the children
	for i := 1; i < len(children); i++ {
		child := children[i]

		child.NextSibling = nil
		child.PrevSibling = nil

		last_child := f.LastChild
		last_child.NextSibling = child
		child.PrevSibling = last_child

		child.Parent = f
		f.LastChild = child
	}
}

// GetChildren returns the immediate children of the node.
//
// The returned nodes are never nil and are not copied. Thus, modifying the returned
// nodes will modify the tree.
//
// Returns:
//   - []Noder: A slice of pointers to the children of the node.
func (f *Float64) GetChildren() []Noder {
	var children []Noder

	for c := f.FirstChild; c != nil; c = c.NextSibling {
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
func (f *Float64) HasChild(target *Float64) bool {
	if target == nil || f.FirstChild == nil {
		return false
	}

	for c := f.FirstChild; c != nil; c = c.NextSibling {
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
func (f *Float64) delete_child(target *Float64) []Noder {
	ok := f.HasChild(target)
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

	if target == f.FirstChild {
		f.FirstChild = next

		if next == nil {
			f.LastChild = nil
		}
	} else if target == f.LastChild {
		f.LastChild = prev
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
func (f *Float64) IsChildOf(target *Float64) bool {
	if target == nil {
		return false
	}

	parents := target.GetAncestors()

	for node := f; node.Parent != nil; node = node.Parent {
		parent := Noder(node.Parent)

		ok := slices.Contains(parents, parent)
		if ok {
			return true
		}
	}

	return false
}
