package tree

import (
	"iter"
	"strings"

	"github.com/PlayerR9/tree/tree/internal"
)

// traversor is a tree traversor.
type traversor[N interface {
	BackwardChild() iter.Seq[N]
	TreeNoder
}] struct {
	// node is the current node.
	node N

	// indent is the current indentation.
	indent string

	// has_sibling is true if the current node has a sibling, false otherwise.
	has_sibling bool
}

// new_traversor is a helper function that creates a new traversor.
//
// Parameters:
//   - node: the current node.
//   - indent: the current indentation.
//   - has_sibling: true if the current node has a sibling, false otherwise.
//
// Returns:
//   - traversor: the new traversor.
func new_traversor[N interface {
	BackwardChild() iter.Seq[N]
	TreeNoder
}](node N, indent string, has_sibling bool) traversor[N] {
	return traversor[N]{
		node:        node,
		indent:      indent,
		has_sibling: has_sibling,
	}
}

// toggle_sibling is a helper function that toggles the sibling flag.
//
// Parameters:
//   - has_sibling: true if the current node has a sibling, false otherwise.
//
// Does nothing if the receiver is nil.
func (t *traversor[N]) toggle_sibling(has_sibling bool) {
	if t == nil {
		return
	}

	t.has_sibling = has_sibling
}

// printer is a tree printer.
type printer[N interface {
	BackwardChild() iter.Seq[N]
	TreeNoder
}] struct {
	// builder is the string builder.
	builder strings.Builder
}

// PrintTree is a function that prints the tree.
//
// Parameters:
//   - root: the root node of the tree.
//
// Returns:
//   - string: the string representation of the tree.
func PrintTree[N interface {
	BackwardChild() iter.Seq[N]
	TreeNoder
}](root N) string {
	if root.IsLeaf() {
		return root.String()
	}

	var p printer[N]

	p.builder.WriteString(root.String())

	var elems []traversor[N]

	var stack internal.Stack[traversor[N]]

	for child := range root.BackwardChild() {
		elems = append(elems, new_traversor(child, "", true))
	}

	elems[0].toggle_sibling(false)

	for _, elem := range elems {
		stack.Push(elem)
	}

	for {
		top, ok := stack.Pop()
		if !ok {
			break
		}

		p.builder.WriteRune('\n')

		p.builder.WriteString(top.indent)

		if top.has_sibling {
			p.builder.WriteString(" ├── ")
		} else {
			p.builder.WriteString(" └── ")
		}

		p.builder.WriteString(top.node.String())

		if top.node.IsLeaf() {
			continue
		}

		var new_indent string

		if top.has_sibling {
			new_indent = top.indent + " │  "
		} else {
			new_indent = top.indent + "    "
		}

		var elems []traversor[N]

		for child := range top.node.BackwardChild() {
			elems = append(elems, new_traversor(child, new_indent, true))
		}

		elems[0].toggle_sibling(false)

		for _, elem := range elems {
			stack.Push(elem)
		}
	}

	return p.builder.String()
}
