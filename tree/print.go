package tree

import (
	"iter"
	"strings"
)

// TreeStringTraverser is the tree stringer.
type TreeStringTraverser[T TreeNoder] struct {
	// lines is the lines of the tree stringer.
	lines []string

	// seen is the seen map of the tree stringer.
	seen map[T]bool
}

// String implements the fmt.Stringer interface.
func (tst TreeStringTraverser[T]) String() string {
	return strings.Join(tst.lines, "\n")
}

// IsSeen is a helper function that checks if the node is seen.
//
// Parameters:
//   - node: The node to check.
//
// Returns:
//   - bool: The result of the check.
func (tst TreeStringTraverser[T]) IsSeen(node T) bool {
	prev, ok := tst.seen[node]
	return ok && prev
}

// AppendLine is a helper function that appends a line to the tree stringer.
//
// Parameters:
//   - line: The line to append.
func (tst *TreeStringTraverser[T]) AppendLine(line string) {
	tst.lines = append(tst.lines, line)
}

// SetSeen is a helper function that sets the seen flag.
//
// Parameters:
//   - node: The node to set.
func (tst *TreeStringTraverser[T]) SetSeen(node T) {
	tst.seen[node] = true
}

// TreeStackElem is the stack element of the tree stringer.
type TreeStackElem[T TreeNoder] struct {
	// global contains the global info of the tree stringer.
	global *TreeStringTraverser[T]

	// indent is the indentation string.
	indent string

	// is_last is the flag that indicates whether the node is the last node in the level.
	is_last bool

	// same_level is the flag that indicates whether the node is in the same level.
	same_level bool
}

// String implements the fmt.Stringer interface.
func (tse TreeStackElem[T]) String() string {
	return tse.global.String()
}

// set_is_last is a helper function that sets the is_last flag.
func (tse *TreeStackElem[T]) set_is_last() {
	tse.is_last = true
}

// set_same_level is a helper function that sets the same_level flag.
func (tse *TreeStackElem[T]) set_same_level() {
	tse.same_level = true
}

// PrintFn returns the print function of the tree stringer.
//
// Parameters:
//   - root: The root node of the tree.
//
// Returns:
//   - Traverser[T, *TreeStackElem[T]]: The print function of the tree stringer.
func PrintFn[T interface {
	Child() iter.Seq[T]
	BackwardChild() iter.Seq[T]
	Copy() T
	LinkChildren(children []T)
	TreeNoder
}]() Traverser[T, *TreeStackElem[T]] {
	init_fn := func(root T) *TreeStackElem[T] {
		return &TreeStackElem[T]{
			global: &TreeStringTraverser[T]{
				lines: make([]string, 0),
				seen:  make(map[T]bool),
			},
			indent:     "",
			is_last:    true,
			same_level: false,
		}
	}

	fn := func(node T, info *TreeStackElem[T]) ([]Pair[T, *TreeStackElem[T]], error) {
		var builder strings.Builder

		if info.indent != "" {
			builder.WriteString(info.indent)

			if !node.IsLeaf() || info.is_last {
				builder.WriteString("└── ")
			} else {
				builder.WriteString("├── ")
			}
		}

		// Prevent cycles.
		ok := info.global.IsSeen(node)
		if ok {
			builder.WriteString("... WARNING: Cycle detected!")

			info.global.AppendLine(builder.String())

			return nil, nil
		}

		builder.WriteString(node.String())
		info.global.AppendLine(builder.String())

		info.global.SetSeen(node)

		if node.IsLeaf() {
			return nil, nil
		}

		var indent strings.Builder

		indent.WriteString(info.indent)

		if info.same_level && !info.is_last {
			indent.WriteString("│   ")
		} else {
			indent.WriteString("    ")
		}

		var elems []Pair[T, *TreeStackElem[T]]

		for c := range node.Child() {
			se := &TreeStackElem[T]{
				global:     info.global,
				indent:     indent.String(),
				is_last:    false,
				same_level: false,
			}

			elems = append(elems, NewPair(c, se))
		}

		if len(elems) >= 2 {
			for i := 0; i < len(elems); i++ {
				elems[i].Info.set_same_level()
			}
		}

		elems[len(elems)-1].Info.set_is_last()

		return elems, nil
	}

	return Traverser[T, *TreeStackElem[T]]{
		InitFn: init_fn,
		DoFn:   fn,
	}
}
