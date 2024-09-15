package tree

import "errors"

var (
	// NodeNotPartOfTree is an error that is returned when a node is not part of a tree.
	NodeNotPartOfTree error
)

func init() {
	NodeNotPartOfTree = errors.New("node is not part of the tree")
}
