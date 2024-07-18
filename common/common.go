package common

import (
	"strings"
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
func FString[T any](t Treer[T]) string {
	iter := t.DFS()

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
