package Tree

import (
	"fmt"
	"strings"

	uc "github.com/PlayerR9/MyGoLib/Units/common"
	tn "github.com/PlayerR9/treenode"
)

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
func PrintTree[T tn.Noder](tree *Tree[T]) ([]string, error) {
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
