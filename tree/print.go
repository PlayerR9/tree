package tree

import (
	"iter"
	"strings"
)

// _TreePrinterTrav is the stack element of the tree stringer.
type _TreePrinterTrav[T TreeNoder] struct {
	// seen is the seen map of the tree stringer.
	seen map[T]struct{}

	// builder is the builder of the tree stringer.
	builder *strings.Builder

	// indent is the indentation string.
	indent string

	// is_last is the flag that indicates whether the node is the last node in the level.
	is_last bool

	// same_level is the flag that indicates whether the node is in the same level.
	same_level bool
}

// String implements the fmt.Stringer interface.
func (tse _TreePrinterTrav[T]) String() string {
	str := tse.builder.String()
	return strings.TrimSuffix(str, "\n")
}

// set_is_last is a helper function that sets the is_last flag.
//
// Assumes that the receiver is not nil.
func (tse *_TreePrinterTrav[T]) set_is_last() {
	tse.is_last = true
}

// set_same_level is a helper function that sets the same_level flag.
//
// Assumes that the receiver is not nil.
func (tse *_TreePrinterTrav[T]) set_same_level() {
	tse.same_level = true
}

// print_fn returns the print function of the tree stringer.
//
// Parameters:
//   - root: The root node of the tree.
//
// Returns:
//   - Traverser[T]: The print function of the tree stringer.
func print_fn[T interface {
	Child() iter.Seq[T]
	BackwardChild() iter.Seq[T]
	Copy() T
	LinkChildren(children []T)
	TreeNoder
}]() Traverser[T] {
	init_fn := func(root T) any {
		var builder strings.Builder

		return &_TreePrinterTrav[T]{
			seen:       make(map[T]struct{}),
			builder:    &builder,
			indent:     "",
			is_last:    true,
			same_level: false,
		}
	}

	fn := func(node T, info any) ([]Pair[T], error) {
		inf := info.(*_TreePrinterTrav[T])

		if inf.indent != "" {
			inf.builder.WriteString(inf.indent)

			if !node.IsLeaf() || inf.is_last {
				inf.builder.WriteString("└── ")
			} else {
				inf.builder.WriteString("├── ")
			}
		}

		// Prevent cycles.
		_, ok := inf.seen[node]
		if ok {
			inf.builder.WriteString("... WARNING: Cycle detected!\n")

			return nil, nil
		}

		inf.builder.WriteString(node.String())
		inf.builder.WriteRune('\n')

		inf.seen[node] = struct{}{}

		if node.IsLeaf() {
			return nil, nil
		}

		var indent strings.Builder

		indent.WriteString(inf.indent)

		if inf.same_level && !inf.is_last {
			indent.WriteString("│   ")
		} else {
			indent.WriteString("    ")
		}

		var elems []Pair[T]

		for c := range node.Child() {
			se := &_TreePrinterTrav[T]{
				seen:       inf.seen,
				builder:    inf.builder,
				indent:     indent.String(),
				is_last:    false,
				same_level: false,
			}

			elems = append(elems, NewPair(c, se))
		}

		if len(elems) >= 2 {
			for i := 0; i < len(elems); i++ {
				elems[i].Info.(*_TreePrinterTrav[T]).set_same_level()
			}
		}

		elems[len(elems)-1].Info.(*_TreePrinterTrav[T]).set_is_last()

		return elems, nil
	}

	return Traverser[T]{
		InitFn: init_fn,
		DoFn:   fn,
	}
}
