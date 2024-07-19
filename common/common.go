package common

import (
	"fmt"
	"slices"
	"strings"

	"github.com/PlayerR9/MyGoLib/ListLike/Queuer"
	lls "github.com/PlayerR9/MyGoLib/ListLike/Stacker"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

/*
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

*/

const (
	// DefaultSpacing is the default spacing between nodes.
	DefaultSpacing string = "|   "

	// DefaultLeafPrefix is the default prefix for leaves.
	DefaultLeafPrefix string = "├── "

	// DefaultNodePrefix is the default prefix for nodes.
	DefaultNodePrefix string = "└── "
)

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
//
// Parameters:
//   - t: The tree to format.
//
// Returns:
//   - string: The formatted string.
func FString(t Treer) string {
	iter := NewDFSIterator(t)

	elem, err := iter.Consume()
	if err != nil {
		return ""
	}

	var builder strings.Builder

	// Deal with root.
	str := elem.Node.String()

	builder.WriteString(str)

	// Deal with children.

	for {
		node, err := iter.Consume()
		if err != nil {
			break
		}

		builder.WriteRune('\n')

		builder.WriteString(strings.Repeat(DefaultSpacing, node.Depth))

		ok := node.Node.IsLeaf()
		if ok {
			builder.WriteString(DefaultLeafPrefix)
		} else {
			builder.WriteString(DefaultNodePrefix)
		}

		builder.WriteString(node.Node.String())
	}

	return builder.String()
}

// RegenerateLeaves regenerates the leaves of the tree. No op if the tree is nil.
//
// Parameters:
//   - tree: The tree to regenerate the leaves of.
//
// Behaviors:
//   - The leaves are updated in a DFS order.
//   - Expensive operation; use it only when necessary (i.e., leaves changed unexpectedly.)
//   - This also updates the size of the tree.
func RegenerateLeaves(tree Treer) {
	if tree == nil {
		return
	}

	root := tree.Root()
	if root == nil {
		return
	}

	var leaves []Noder

	iter := NewDFSIterator(tree)

	var size int

	for {
		elem, err := iter.Consume()
		if err != nil {
			break
		}

		size++

		ok := elem.Node.IsLeaf()
		if ok {
			leaves = append(leaves, elem.Node)
		}
	}

	tree.SetLeaves(leaves, size)
}

// UpdateLeaves updates the leaves of the tree. No op if the tree is nil.
//
// Parameters:
//   - tree: The tree to update.
//
// Behaviors:
//   - The leaves are updated in a DFS order.
//   - Less expensive than RegenerateLeaves. However, if nodes has been deleted
//     from the tree, this may give unexpected results.
//   - This also updates the size of the tree.
func UpdateLeaves(tree Treer) {
	if tree == nil {
		return
	}

	leaves := tree.GetLeaves()
	if len(leaves) == 0 {
		tree.SetLeaves(nil, 0)
		return
	}

	var new_leaves []Noder
	size := tree.Size() - len(leaves)

	stack := lls.NewArrayStack(leaves...)

	for {
		top, ok := stack.Pop()
		if !ok {
			break
		}

		size++

		ok = top.IsLeaf()
		if ok {
			new_leaves = append(new_leaves, top)
		}
	}

	tree.SetLeaves(new_leaves, size)
}

// HasChild returns true if the tree has the given child in any of its nodes
// in a BFS order.
//
// Parameters:
//   - tree: The tree to filter.
//   - filter: The filter to apply. Must return true iff the node is the one we are looking for.
//     This function must assume node is never nil.
//
// Returns:
//   - bool: True if the tree has the child, false otherwise.
//
// If either tree or filter is nil, false is returned.
func HasChild(tree Treer, filter func(node Noder) bool) bool {
	if tree == nil || filter == nil {
		return false
	}

	iter := NewBFSIterator(tree)

	for {
		value, err := iter.Consume()
		if err != nil {
			break
		}

		ok := filter(value.Node)
		if ok {
			return true
		}
	}

	return false
}

// FilterChildren returns all the children of the tree that satisfy the given filter
// in a BFS order.
//
// Parameters:
//   - tree: The tree to filter.
//   - filter: The filter to apply. Must return true iff the node is the one we want to keep.
//     This function must assume node is never nil.
//
// Returns:
//   - []T: A slice of the children that satisfy the filter.
//   - bool: True if all the nodes are of type T, false otherwise.
//
// If either tree or filter is nil, an empty slice and false are returned.
func FilterChildren[T Noder](tree Treer, filter func(node T) bool) ([]T, bool) {
	if tree == nil || filter == nil {
		return nil, true
	}

	iter := NewBFSIterator(tree)

	var children []T

	for {
		value, err := iter.Consume()
		if err != nil {
			break
		}

		tmp, ok := value.Node.(T)
		if !ok {
			return nil, false
		}

		ok = filter(tmp)
		if ok {
			children = append(children, tmp)
		}
	}

	return children, true
}

// SearchNodes searches for the first node that satisfies the given filter in a BFS order.
//
// Parameters:
//   - tree: The tree to search.
//   - filter: The filter to apply. Must return true iff the node is the one we are looking for.
//     This function must assume node is never nil.
//
// Returns:
//   - T: The node that satisfies the filter.
//   - bool: True if the node was found, false otherwise.
//
// Nodes that are not of type T will be ignored. If either tree or filter is nil, false is returned.
func SearchNodes[T Noder](tree Treer, filter func(node T) bool) (T, bool) {
	if tree == nil || filter == nil {
		return *new(T), false
	}

	iter := NewBFSIterator(tree)

	for {
		value, err := iter.Consume()
		if err != nil {
			break
		}

		tmp, ok := value.Node.(T)
		if ok {
			ok = filter(tmp)
			if ok {
				return tmp, true
			}
		}
	}

	return *new(T), false
}

// FindBranchingPoint returns the first node in the path from n to the root
// such that has more than one sibling.
//
// Parameters:
//   - n: The node to start the search.
//
// Returns:
//   - Noder: The branching point. Nil if no branching point was found.
//   - Noder: The parent of the branching point. Nil if n is nil.
//   - bool: True if the node has a branching point, false otherwise.
//
// Behaviors:
//   - If there is no branching point, it returns the root of the tree. However,
//     if n is nil, it returns nil, nil, false and if the node has no parent, it
//     returns nil, n, false.
func FindBranchingPoint(n Noder) (Noder, Noder, bool) {
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

// DeleteBranchContaining deletes the branch containing the given node.
//
// Parameters:
//   - n: The node to delete.
//
// Returns:
//   - error: An error if the node is not a part of the tree.
func DeleteBranchContaining[T Noder](tree Treer, n T) error {
	root := tree.Root()
	if root == nil {
		return NewErrNodeNotPartOfTree()
	}

	child, parent, hasBranching := FindBranchingPoint(n)
	if !hasBranching {
		if parent != root {
			return NewErrNodeNotPartOfTree()
		}

		tree.Cleanup()
	}

	children := parent.DeleteChild(child)

	for i := 0; i < len(children); i++ {
		current := children[i]

		current.Cleanup()

		children[i] = nil
	}

	RegenerateLeaves(tree)

	return nil
}

// PruneTree prunes the tree using the given filter.
//
// Parameters:
//   - tree: The tree to prune.
//   - filter: The filter to use to prune the tree. Must return true iff the node
//     should be pruned.
//
// Returns:
//   - bool: False if the whole tree can be deleted, true otherwise.
func Prune(tree Treer, filter func(node Noder) bool) bool {
	if tree == nil {
		return false
	}

	for tree.Size() != 0 {
		target, ok := SearchNodes(tree, filter)
		if !ok {
			return true
		}

		DeleteBranchContaining(tree, target)
	}

	return false
}

// ExtractBranch extracts the branch of the tree that contains the given leaf.
//
// Parameters:
//   - tree: The tree to search.
//   - leaf: The leaf to extract the branch from.
//   - delete: If true, the branch is deleted from the tree.
//
// Returns:
//   - *Branch[T]: A pointer to the branch extracted. Nil if the leaf is not a part
//     of the tree. Nil if the leaf is not a part of the tree and delete is false.
//
// Behaviors:
//   - If delete is true, then the branch is deleted from the tree.
func ExtractBranch[T Noder](tree Treer, leaf T, delete bool) *Branch[T] {
	if tree == nil {
		return nil
	}

	found := slices.Contains(tree.GetLeaves(), Noder(leaf))
	if !found {
		return nil
	}

	branch, err := NewBranch[T](leaf)
	uc.AssertErr(err, "NewBranch[%T](%s)", leaf, leaf.String())

	if !delete {
		return branch
	}

	child, parent, ok := FindBranchingPoint(leaf)
	if !ok {
		parent.DeleteChild(child)
	}

	RegenerateLeaves(tree)

	return branch
}

// InsertBranch inserts the given branch into the tree.
//
// Parameters:
//   - tree: The tree to insert the branch into.
//   - branch: The branch to insert.
//
// Returns:
//   - T: The updated tree.
//   - error: An error if the insertion fails.
func InsertBranch[T Treer, N Noder](tree T, branch *Branch[N]) (T, error) {
	if branch == nil {
		return tree, nil
	}

	ref := tree.Root()

	if ref == nil {
		other_tree := branch.from_node.ToTree()

		tmp, ok := other_tree.(T)
		if !ok {
			return *new(T), fmt.Errorf("other_tree is not a tree: %T", other_tree)
		}

		return tmp, nil
	}

	var from Noder

	from = branch.from_node

	if ref != from {
		return tree, nil
	}

	for from != Noder(branch.to_node) {
		from = from.GetFirstChild()

		var next Noder

		c := ref.GetFirstChild()

		for c != nil && next == nil {
			if c == from {
				next = c
			}

			c = c.GetFirstChild()
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

	RegenerateLeaves(tree)

	return tree, nil
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
func DFS[T Noder](tree Treer, init Infoer, f ObserverFunc[T]) error {
	if f == nil || tree == nil {
		return nil
	}

	root := tree.Root()

	tmp, ok := root.(T)
	if !ok {
		return fmt.Errorf("root is not a tree: %T", root)
	}

	trav := new_traversor(tmp, init)

	S := lls.NewLinkedStack(trav)

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

			tmp, ok := c.(T)
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
func BFS[T Noder](tree Treer, init Infoer, f ObserverFunc[T]) error {
	if f == nil || tree == nil {
		return nil
	}

	root := tree.Root()

	tmp, ok := root.(T)
	if !ok {
		return fmt.Errorf("root is not a tree: %T", root)
	}

	trav := new_traversor(tmp, init)

	Q := Queuer.NewLinkedQueue(trav)

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

			tmp, ok := c.(T)
			if !ok {
				return fmt.Errorf("node is not a tree: %T", c)
			}

			new_t := new_traversor(tmp, first.info)

			Q.Enqueue(new_t)
		}
	}

	return nil
}

// InfPrinter is a struct that prints the tree.
type InfPrinter struct {
	// indent_level is the level of indentation.
	indent_level int
}

// Copy implements the common.Copier interface.
func (ip *InfPrinter) Copy() uc.Copier {
	ip_copy := &InfPrinter{
		indent_level: ip.indent_level,
	}

	return ip_copy
}

// NewInfPrinter creates a new InfPrinter.
//
// Returns:
//   - *InfPrinter: The new InfPrinter.
func NewInfPrinter() *InfPrinter {
	ip := &InfPrinter{
		indent_level: 0,
	}
	return ip
}

// IncIndent increments the indentation level.
func (ip *InfPrinter) IncIndent() {
	ip.indent_level++
}

// PrintTree prints the tree.
//
// Parameters:
//   - tree: The tree to print.
//
// Returns:
//   - []string: The lines of the tree.
//   - error: An error if the tree cannot be printed.
func PrintTree[T Noder](tree Treer) ([]string, error) {
	if tree == nil {
		return nil, nil
	}

	var lines []string
	var builder strings.Builder

	f := func(elem T, obj Infoer) (bool, error) {
		inf, ok := obj.(*InfPrinter)
		if !ok {
			return false, fmt.Errorf("invalid objecter type: %T", obj)
		}

		builder.WriteString(strings.Repeat("| ", inf.indent_level))
		builder.WriteString(uc.StringOf(elem))
		builder.WriteString("\n")

		inf.IncIndent()

		return true, nil
	}

	ip := NewInfPrinter()

	err := DFS(tree, ip, f)
	if err != nil {
		return nil, err
	}

	return lines, nil
}
