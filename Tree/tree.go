package Tree

import (
	"slices"
	"strings"
	"unicode/utf8"

	ffs "github.com/PlayerR9/MyGoLib/Formatting/FString"
	"github.com/PlayerR9/MyGoLib/ListLike/Queuer"
	lls "github.com/PlayerR9/MyGoLib/ListLike/Stacker"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
	us "github.com/PlayerR9/MyGoLib/Units/slice"
)

// TreeFormatter is a formatter that formats the tree.
type TreeFormatter struct {
	// spacing is the spacing between nodes.
	spacing string

	// leaf_prefix is the prefix for leaves.
	leaf_prefix string

	// node_prefix is the prefix for nodes.
	node_prefix string
}

// WithSpacing sets the spacing between nodes.
//
// If spacing is an empty string, it is set to three spaces.
//
// Parameters:
//   - spacing: The spacing between nodes.
//
// Returns:
//   - ffs.Option: The option function.
func WithSpacing(spacing string) ffs.Option {
	size := utf8.RuneCountInString(spacing)
	if size <= 1 {
		spacing = "   "
	}

	p1 := strings.Repeat("─", size-1)
	p2 := strings.Repeat(spacing, size)

	return func(s ffs.Settinger) {
		tf, ok := s.(*TreeFormatter)
		if !ok {
			return
		}

		var builder strings.Builder

		builder.WriteRune('|')
		builder.WriteString(p2)

		tf.spacing = builder.String()
		builder.Reset()

		builder.WriteRune('├')
		builder.WriteString(p1)
		builder.WriteRune(' ')

		tf.leaf_prefix = builder.String()
		builder.Reset()

		builder.WriteRune('└')
		builder.WriteString(p1)
		builder.WriteRune(' ')

		tf.node_prefix = builder.String()
	}
}

// Tree is a generic data structure that represents a tree.
type Tree[T any] struct {
	// root is the root of the tree.
	root *TreeNode[T]

	// leaves is the leaves of the tree.
	leaves []*TreeNode[T]

	// size is the number of nodes in the tree.
	size int
}

// FString implements the FString.FStringer interface.
//
// By default, it uses a three-space indentation.
//
// Format:
//
//	root
//	├── node1
//	│   ├── node2
//	│   └── node3
//	└── node4
//	|   └── node5
//	|
//	| // ...
//	// ...
func (t *Tree[T]) FString(trav *ffs.Traversor, opts ...ffs.Option) error {
	if trav == nil {
		return nil
	}

	root := t.root
	if root == nil {
		return nil
	}

	tf := &TreeFormatter{
		spacing:     "|   ",
		leaf_prefix: "├── ",
		node_prefix: "└── ",
	}

	for _, opt := range opts {
		opt(tf)
	}

	// Deal with root.
	str := root.String()

	err := trav.AddLine(str)
	if err != nil {
		return err
	}

	type stack_element struct {
		node  *TreeNode[T]
		depth string
	}

	stack := lls.NewArrayStack[*stack_element]()

	for c := root.FirstChild; c != nil; c = c.NextSibling {
		new_se := &stack_element{
			node:  c,
			depth: tf.spacing,
		}

		stack.Push(new_se)
	}

	// Deal with the rest of the tree.

	for {
		top, ok := stack.Pop()
		if !ok {
			break
		}

		var builder strings.Builder

		builder.WriteString(top.depth)

		ok = top.node.IsLeaf()
		if ok {
			builder.WriteString(tf.leaf_prefix)
		} else {
			builder.WriteString(tf.node_prefix)
		}

		str := top.node.String()

		builder.WriteString(str)

		str = builder.String()

		err := trav.AddLine(str)
		if err != nil {
			return err
		}

		for c := top.node.FirstChild; c != nil; c = c.NextSibling {
			new_se := &stack_element{
				node:  c,
				depth: top.depth + tf.spacing,
			}

			stack.Push(new_se)
		}
	}

	return nil
}

// Cleanup implements the object.Cleaner interface.
func (t *Tree[T]) Cleanup() {
	root := t.root
	if root == nil {
		return
	}

	root.Cleanup()

	t.root = nil
}

// Copy implements the common.Copier interface.
func (t *Tree[T]) Copy() uc.Copier {
	root := t.root
	if root == nil {
		tree := &Tree[T]{
			root:   nil,
			leaves: nil,
			size:   0,
		}

		return tree
	}

	var tree *Tree[T]

	root_copy := root.Copy().(*TreeNode[T])

	tree = &Tree[T]{
		root:   root_copy,
		leaves: root_copy.GetLeaves(),
		size:   t.size,
	}

	return tree
}

// NewTree creates a new tree with the given value as the root.
//
// Parameters:
//   - data: The value of the root.
//
// Returns:
//   - *Tree: A pointer to the newly created tree.
func NewTree[T any](root *TreeNode[T]) *Tree[T] {
	if root == nil {
		tree := &Tree[T]{
			root:   nil,
			leaves: nil,
			size:   0,
		}

		return tree
	}

	var leaves []*TreeNode[T]
	var size int

	ok := root.IsLeaf()
	if ok {
		leaves = []*TreeNode[T]{root}
		size = 1
	} else {
		leaves = root.GetLeaves()
		size = root.Size()
	}

	tree := &Tree[T]{
		root:   root,
		leaves: leaves,
		size:   size,
	}

	return tree
}

// SetChildren sets the children of the root of the tree.
//
// Parameters:
//   - children: The children to set.
//
// Returns:
//   - error: An error of type *ErrMissingRoot if the tree does not have a root.
func (t *Tree[T]) SetChildren(children []*Tree[T]) error {
	children = us.SliceFilter(children, FilterNonNilTree)
	if len(children) == 0 {
		return nil
	}

	root := t.root
	if root == nil {
		return NewErrMissingRoot()
	}

	var leaves, sub_children []*TreeNode[T]

	t.size = 1

	for _, child := range children {
		leaves = append(leaves, child.leaves...)
		t.size += child.Size()

		croot := child.root
		croot.Parent = root

		sub_children = append(sub_children, croot)
	}

	root.LinkChildren(sub_children)

	t.leaves = leaves

	return nil
}

// IsSingleton returns true if the tree has only one node.
//
// Returns:
//   - bool: True if the tree has only one node, false otherwise.
func (t *Tree[T]) IsSingleton() bool {
	return t.size == 1
}

// Size returns the number of nodes in the tree.
//
// Returns:
//   - int: The number of nodes in the tree.
func (t *Tree[T]) Size() int {
	return t.size
}

// Root returns the root of the tree.
//
// Returns:
//   - T: The root of the tree. Nil if the tree does not have a root.
func (t *Tree[T]) Root() *TreeNode[T] {
	return t.root
}

/*

// GetChildren returns all the children of the tree in a DFS order.
//
// Returns:
//   - children: A slice of the values of the children of the tree.
//
// Behaviors:
//   - The root is the first element in the slice.
//   - If the tree does not have a root, it returns nil.
func (t *Tree) GetChildren() (children []T) {
	root := t.root
	if root == nil {
		return nil
	}

	S := Stacker.NewLinkedStack(root)

	for {
		node, ok := S.Pop()
		if !ok {
			break
		}

		children = append(children, node.Data)

		for i := 0; i < len(node.children); i++ {
			current := node.children[i]

			S.Push(current)
		}
	}

	return children
}
*/

// GetLeaves returns all the leaves of the tree.
//
// Returns:
//   - []T: A slice of the leaves of the tree. Nil if the tree does not have a root.
//
// Behaviors:
//   - It returns the leaves that are stored in the tree. Make sure to call
//     any update function before calling this function if the tree has been modified
//     unexpectedly.
func (t *Tree[T]) GetLeaves() []*TreeNode[T] {
	return t.leaves
}

// RegenerateLeaves regenerates the leaves of the tree and returns them.
//
// Returns:
//   - []T: The newly generated leaves of the tree.
//   - error: An error if the leaves could not be generated or the nodes are not of type T.
//
// Behaviors:
//   - The leaves are updated in a DFS order.
//   - Expensive operation; use it only when necessary (i.e., leaves changed unexpectedly.)
//   - This also updates the size of the tree.
func (t *Tree[T]) RegenerateLeaves() ([]*TreeNode[T], error) {
	root := t.root
	if root == nil {
		t.leaves = nil
		t.size = 0

		return nil, nil
	}

	var leaves []*TreeNode[T]

	S := lls.NewLinkedStack(root)

	t.size = 0

	for {
		top, ok := S.Pop()
		if !ok {
			break
		}
		uc.Assert(top != nil, "top is nil")

		t.size++

		ok = top.IsLeaf()
		if ok {
			leaves = append(leaves, top)
			continue
		}

		for c := top.FirstChild; c != nil; c = c.NextSibling {
			S.Push(c)
		}
	}

	t.leaves = leaves

	return leaves, nil
}

// UpdateLeaves updates the leaves of the tree and returns them.
//
// Returns:
//   - []T: The newly generated leaves of the tree.
//   - error: An error if the leaves could not be generated or the nodes are not of type T.
//
// Behaviors:
//   - The leaves are updated in a DFS order.
//   - Less expensive than RegenerateLeaves. However, if nodes has been deleted
//     from the tree, this may give unexpected results.
//   - This also updates the size of the tree.
func (t *Tree[T]) UpdateLeaves() ([]*TreeNode[T], error) {
	if len(t.leaves) == 0 {
		t.size = 0
		return nil, nil
	}

	var leaves []*TreeNode[T]

	S := lls.NewLinkedStack(t.leaves...)

	t.size -= len(t.leaves)

	for {
		top, ok := S.Pop()
		if !ok {
			break
		}

		uc.Assert(top != nil, "top is nil")

		t.size++

		ok = top.IsLeaf()
		if ok {
			leaves = append(leaves, top)
			continue
		}

		for c := top.FirstChild; c != nil; c = c.NextSibling {
			S.Push(c)
		}
	}

	t.leaves = leaves

	return leaves, nil
}

// HasChild returns true if the tree has the given child in any of its nodes
// in a BFS order.
//
// The filter must assume that the node is never nil.
//
// Parameters:
//   - filter: The filter to apply.
//
// Returns:
//   - bool: True if the tree has the child, false otherwise.
//   - error: An error if the child is not of type T.
func (t *Tree[T]) HasChild(filter us.PredicateFilter[*TreeNode[T]]) (bool, error) {
	if filter == nil {
		return false, nil
	}

	root := t.root
	if root == nil {
		return false, nil
	}

	Q := Queuer.NewLinkedQueue(root)

	for {
		node, ok := Q.Dequeue()
		if !ok {
			break
		}

		uc.Assert(node != nil, "node is nil")

		ok = filter(node)
		if ok {
			return true, nil
		}

		for c := node.FirstChild; c != nil; c = c.NextSibling {
			Q.Enqueue(c)
		}
	}

	return false, nil
}

// FilterChildren returns all the children of the tree that satisfy the given filter
// in a BFS order.
//
// The filter must assume that the node is never nil.
//
// Parameters:
//   - filter: The filter to apply.
//
// Returns:
//   - []T: A slice of the children that satisfy the filter.
//   - error: An error if iterating over the children fails.
func (t *Tree[T]) FilterChildren(filter us.PredicateFilter[*TreeNode[T]]) ([]*TreeNode[T], error) {
	if filter == nil {
		return nil, nil
	}

	root := t.root
	if root == nil {
		return nil, nil
	}

	Q := Queuer.NewLinkedQueue(root)

	var children []*TreeNode[T]

	for {
		node, ok := Q.Dequeue()
		if !ok {
			break
		}

		uc.Assert(node != nil, "node is nil")

		ok = filter(node)
		if ok {
			children = append(children, node)
		}

		for c := node.FirstChild; c != nil; c = c.NextSibling {
			Q.Enqueue(c)
		}
	}

	return children, nil
}

// PruneBranches removes all the children of the node that satisfy the given filter.
// The filter is a function that takes the value of a node and returns a boolean.
// If the filter returns true for a child, the child is removed along with its children.
//
// Parameters:
//   - filter: The filter to apply.
//
// Returns:
//   - bool: True if the whole tree can be deleted, false otherwise.
//
// Behaviors:
//   - If the root satisfies the filter, the tree is cleaned up.
//   - It is a recursive function.
func (t *Tree[T]) PruneBranches(filter us.PredicateFilter[*TreeNode[T]]) bool {
	if filter == nil {
		return false
	}

	root := t.root
	if root == nil {
		return true
	}

	highest, ok := rec_prune_func(filter, nil, root)
	if ok {
		return true
	}

	t.leaves = highest.GetLeaves()
	t.size = highest.Size()

	return false
}

// SearchNodes searches for the first node that satisfies the given filter in a BFS order.
//
// Parameters:
//   - f: The filter to apply.
//
// Returns:
//   - T: The node that satisfies the filter.
//   - error: An error if the node is not found or the iteration fails.
//
// Errors:
//   - *common.ErrNotFound: If the node is not found.
//   - error: The error returned by the iteration function.
func (t *Tree[T]) SearchNodes(f us.PredicateFilter[*TreeNode[T]]) (*TreeNode[T], error) {
	root := t.root
	if root == nil {
		return nil, nil
	}

	Q := Queuer.NewLinkedQueue(root)

	for {
		first, ok := Q.Dequeue()
		if !ok {
			break
		}

		ok = f(first)
		if ok {
			return first, nil
		}

		for c := first.FirstChild; c != nil; c = c.NextSibling {
			Q.Enqueue(c)
		}
	}

	return nil, uc.NewErrNotFound()
}

// DeleteBranchContaining deletes the branch containing the given node.
//
// Parameters:
//   - n: The node to delete.
//
// Returns:
//   - error: An error if the node is not a part of the tree.
func (t *Tree[T]) DeleteBranchContaining(n *TreeNode[T]) error {
	if n == nil {
		return nil
	}

	root := t.root
	if root == nil {
		return NewErrNodeNotPartOfTree()
	}

	child, parent, hasBranching := FindBranchingPoint(n)
	if !hasBranching {
		if parent != root {
			return NewErrNodeNotPartOfTree()
		}

		t.Cleanup()
	}

	children := parent.DeleteChild(child)

	for i := 0; i < len(children); i++ {
		current := children[i]

		current.Cleanup()

		children[i] = nil
	}

	leaves, err := t.RegenerateLeaves()
	if err != nil {
		return err
	}

	t.leaves = leaves

	return nil
}

// PruneTree prunes the tree using the given filter.
//
// Parameters:
//   - filter: The filter to use to prune the tree.
//
// Returns:
//   - bool: True if no nodes were pruned, false otherwise.
//   - error: An error if the iteration fails.
func (t *Tree[T]) Prune(filter us.PredicateFilter[*TreeNode[T]]) (bool, error) {
	for t.Size() != 0 {
		target, err := t.SearchNodes(filter)
		if err != nil {
			return false, err
		}

		if target == nil {
			return true, nil
		}

		t.DeleteBranchContaining(target)
	}

	return false, nil
}

// SkipFunc removes all the children of the tree that satisfy the given filter
// without removing any of their children. Useful for removing unwanted nodes from the tree.
//
// Parameters:
//   - filter: The filter to apply.
//
// Returns:
//   - []*Tree: A slice of pointers to the trees obtained after removing the nodes.
//
// Behaviors:
//   - If this function returns only one tree, this is the updated tree. But, if
//     it returns more than one tree, then we have deleted the root of the tree and
//     obtained a forest.
func (t *Tree[T]) SkipFilter(filter us.PredicateFilter[*TreeNode[T]]) (forest []*Tree[T]) {
	frontier := make([]*TreeNode[T], len(t.leaves))
	copy(frontier, t.leaves)

	seen := make(map[*TreeNode[T]]bool)
	var leaves []*TreeNode[T]

	f := func(n *TreeNode[T]) bool {
		return !seen[n]
	}

	for len(frontier) > 0 {
		leaf := frontier[0]
		seen[leaf] = true

		// Remove any node that has been seen from the frontier.
		frontier = us.SliceFilter(frontier, f)

		ok := filter(leaf)

		parent := leaf.GetParent()

		if !ok {
			if parent == nil {
				// We reached the root
				frontier = frontier[1:]
			} else {
				ok := leaf.IsLeaf()
				if ok {
					leaves = append(leaves, leaf)
				}

				if !seen[parent] {
					frontier[0] = parent
				} else {
					frontier = frontier[1:]
				}
			}
		} else {
			children := leaf.RemoveNode()

			if len(children) != 0 {
				// We obtained a forest as we reached the root

				for i := 0; i < len(children); i++ {
					child := children[i]

					tree := NewTree(child)

					forest = append(forest, tree)
				}

				// We reached the root
				frontier = frontier[1:]
			} else {
				if !seen[parent] {
					frontier[0] = parent
				} else {
					frontier = frontier[1:]
				}

				t.size--
			}
		}
	}

	if len(forest) == 0 {
		t.leaves = leaves

		forest = []*Tree[T]{t}
	}

	return
}

// replaceLeafWithTree is a helper function that replaces a leaf with a tree.
//
// Parameters:
//   - at: The index of the leaf to replace.
//   - children: The children of the leaf.
//
// Behaviors:
//   - The leaf is replaced with the children.
//   - The size of the tree is updated.
func (t *Tree[T]) replaceLeafWithTree(at int, values []*TreeNode[T]) {
	leaf := t.leaves[at]

	// Make the subtree
	leaf.LinkChildren(values)

	// Update the size of the tree
	t.size += len(values) - 1

	// Replace the current leaf with the leaf's children
	sub_leaves := leaf.GetLeaves()

	if at == len(t.leaves)-1 {
		t.leaves = append(t.leaves[:at], sub_leaves...)
	} else if at == 0 {
		t.leaves = append(sub_leaves, t.leaves[at+1:]...)
	} else {
		t.leaves = append(t.leaves[:at], append(sub_leaves, t.leaves[at+1:]...)...)
	}
}

// ProcessLeaves applies the given function to the leaves of the tree and replaces
// the leaves with the children returned by the function.
//
// Parameters:
//   - f: The function to apply to the leaves.
//
// Returns:
//   - error: An error returned by the function.
//
// Behaviors:
//   - The function is applied to the leaves in order.
//   - The function must return a slice of values of type T.
//   - If the function returns an error, the process stops and the error is returned.
//   - The leaves are replaced with the children returned by the function.
func (t *Tree[T]) ProcessLeaves(f uc.EvalManyFunc[*TreeNode[T], *TreeNode[T]]) error {
	for i, leaf := range t.leaves {
		children, err := f(leaf)
		if err != nil {
			return err
		}

		if len(children) != 0 {
			t.replaceLeafWithTree(i, children)
		}
	}

	return nil
}

// GetDirectChildren returns the direct children of the root of the tree.
//
// Children are never nil.
//
// Returns:
//   - []T: A slice of the direct children of the root. Nil if the tree does not have a root.
//   - error: An error if the iteration fails.
func (t *Tree[T]) GetDirectChildren() ([]*TreeNode[T], error) {
	root := t.root
	if root == nil {
		return nil, nil
	}

	var children []*TreeNode[T]

	for c := root.FirstChild; c != nil; c = c.NextSibling {
		uc.Assert(c != nil, "Unexpected nil child")

		children = append(children, c)
	}

	return children, nil
}

// ExtractBranch extracts the branch of the tree that contains the given leaf.
//
// Parameters:
//   - leaf: The leaf to extract the branch from.
//   - delete: If true, the branch is deleted from the tree.
//
// Returns:
//   - *Branch[T]: A pointer to the branch extracted. Nil if the leaf is not a part
//     of the tree.
func (t *Tree[T]) ExtractBranch(leaf *TreeNode[T], delete bool) (*Branch[T], error) {
	found := slices.Contains(t.leaves, leaf)
	if !found {
		return nil, nil
	}

	branch := NewBranch(leaf)

	if !delete {
		return branch, nil
	}

	child, parent, ok := FindBranchingPoint(leaf)
	if !ok {
		parent.DeleteChild(child)
	}

	leaves, err := t.RegenerateLeaves()
	if err != nil {
		return nil, err
	}

	t.leaves = leaves

	return branch, nil
}

// InsertBranch inserts the given branch into the tree.
//
// Parameters:
//   - branch: The branch to insert.
//
// Returns:
//   - bool: True if the branch was inserted, false otherwise.
//   - error: An error if the insertion fails.
func (t *Tree[T]) InsertBranch(branch *Branch[T]) (bool, error) {
	if branch == nil {
		return true, nil
	}

	ref := t.root

	if ref == nil {
		otherTree := NewTree(branch.from_node)

		t.root = otherTree.root
		t.leaves = otherTree.leaves
		t.size = otherTree.size

		return true, nil
	}

	from := branch.from_node
	if ref != from {
		return false, nil
	}

	for from != branch.to_node {
		from = from.FirstChild

		var next *TreeNode[T]

		c := ref.FirstChild

		for c != nil && next == nil {
			if c == from {
				next = c
			}

			c = c.FirstChild
		}

		if next == nil {
			break
		}

		// from is a child of the root. Keep going
		ref = next
	}

	// From this point onward, anything from 'from' up to 'to' must be
	// added in the tree as new children.
	ref.AddChild(from)

	prev_size := t.size

	_, err := t.RegenerateLeaves()
	if err != nil {
		return false, err
	}

	ok := t.size != prev_size
	return ok, nil
}
