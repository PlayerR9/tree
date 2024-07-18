// This command generates a tree node with the given fields that uses first child/next sibling pointers.
//
// To use it, run the following command:
//
// //go:generate go run github.com/PlayerR9/treenode/cmd/treenode -name=<type_name> -fields=<field_list> [ -g=<generics>] [ -o=<output_file> ]
//
// **Flag: Type Name**
//
// The "name" flag is used to specify the name of the tree node struct. As such, it must be set and,
// not only does it have to be a valid Go identifier, but it also must start with an upper case letter.
//
// **Flag: Fields**
//
// The "fields" flag is used to specify the fields that the tree node contains. Because it doesn't make
// a lot of sense to have a tree node without fields, this flag must be set.
//
// Its argument is specified as a list of key-value pairs where each pair is separated by a comma (",") and
// a slash ("/") is used to separate the key and the value.
//
// The key indicates the name of the field while the value indicates the type of the field.
//
// For instance, running the following command:
//
//	//go:generate treenode -name=TreeNode -fields=a/int,b/int,name/string
//
// will generate a tree node with the following fields:
//
//	type TreeNode struct {
//		// Node pointers.
//
//		a int
//		b int
//		name string
//	}
//
// It is important to note that spaces are not allowed.
//
// Also, it is possible to specify generics by following the value with the generics between square brackets;
// like so: "a/MyType[T,C]"
//
// **Flag: Generics**
//
// This optional flag is used to specify the type(s) of the generics. However, this only applies if at least one
// generic type is specified in the fields flag. If none, then this flag is ignored.
//
// As an edge case, if this flag is not specified but the fields flag contains generics, then
// all generics are set to the default value of "any".
//
// As with the fields flag, its argument is specified as a list of key-value pairs where each pair is separated
// by a comma (",") and a slash ("/") is used to separate the key and the value. The key indicates the name of
// the generic and the value indicates the type of the generic.
//
// For instance, running the following command:
//
//	//go:generate treenode -name=TreeNode -fields=a/MyType[T],b/MyType[C] -g=T/any,C/int
//
// will generate a tree node with the following fields:
//
//	type TreeNode[T any, C int] struct {
//		// Node pointers.
//
//		a T
//		b C
//	}
//
// **Flag: Output File**
//
// This optional flag is used to specify the output file. If not specified, the output will be written to
// standard output, that is, the file "<type_name>_treenode.go" in the root of the current directory.
package main

import (
	"flag"
	"log"
	"strings"
	"text/template"

	uc "github.com/PlayerR9/MyGoLib/Units/common"
	ggen "github.com/PlayerR9/MyGoLib/go_generator"
)

var (
	// Logger is the logger for the package.
	Logger *log.Logger

	// t is the template for the tree node.
	t *template.Template
)

func init() {
	Logger = ggen.InitLogger("treenode")

	t = template.Must(template.New("").Parse(templ))
}

var (
	// TypeNameFlag is the flag for the type name.
	TypeNameFlag *string
)

func init() {
	TypeNameFlag = flag.String("name", "",
		"The name of the struct to generate the tree node for. It must be set."+
			" Must start with an upper case letter and must be a valid Go identifier.",
	)

	ggen.SetOutputFlag("<type_name>_treenode.go", false)
	ggen.SetStructFieldsFlag("fields", true, -1, "The fields to generate the code for.")
	ggen.SetGenericsSignFlag("g", false, -1)
}

// GenData is the data for the generator.
type GenData struct {
	// PackageName is the name of the package.
	PackageName string

	// TypeName is the name of the type.
	TypeName string

	// TypeSig is the signature of the type.
	TypeSig string

	// Fields is the list of fields.
	Fields map[string]string

	// IteratorName is the name of the iterator.
	IteratorName string

	// IteratorSig is the signature of the iterator.
	IteratorSig string

	// ParamList is the list of parameters.
	ParamList string

	// AssignmentMap is the map of assignments.
	AssignmentMap map[string]string

	// Generics is the list of generics.
	Generics string
}

// SetPackageName implements the ggen.Generater interface.
func (g GenData) SetPackageName(pkg_name string) ggen.Generater {
	g.PackageName = pkg_name
	return g
}

func main() {
	err := ggen.ParseFlags()
	if err != nil {
		Logger.Fatalf("Could not parse flags: %s", err.Error())
	}

	type_name := uc.AssertNil(TypeNameFlag, "TypeNameFlag")

	err = ggen.IsValidName(type_name, nil, ggen.Exported)
	if err != nil {
		Logger.Fatalf("Invalid type name: %s", err.Error())
	}

	filename, err := ggen.FixOutputLoc(type_name, "_treenode.go")
	if err != nil {
		Logger.Fatalf("Could not fix output location: %s", err.Error())
	}

	err = ggen.Generate(filename, GenData{}, t,
		func(data GenData) GenData {
			data.TypeName = type_name

			return data
		},
		func(data GenData) GenData {
			tn_type_sig, err := ggen.MakeTypeSig(type_name, "")
			if err != nil {
				Logger.Fatalf("Could not generate type signature: %s", err.Error())
			}

			data.TypeSig = tn_type_sig

			return data
		},
		func(data GenData) GenData {
			tn_iterator_sig, err := ggen.MakeTypeSig(type_name, "Iterator")
			if err != nil {
				Logger.Fatalf("Could not generate type signature: %s", err.Error())
			}

			data.IteratorSig = tn_iterator_sig

			return data
		},
		func(data GenData) GenData {
			data.Generics = ggen.GenericsSigFlag.String()

			return data
		},
		func(data GenData) GenData {
			var builder strings.Builder

			builder.WriteString(type_name)
			builder.WriteString("Iterator")

			data.IteratorName = builder.String()

			return data
		},
		func(data GenData) GenData {
			param_list, err := ggen.MakeParameterList()
			if err != nil {
				Logger.Fatalf("Could not generate parameter list: %s", err.Error())
			}

			data.ParamList = param_list

			return data
		},
		func(data GenData) GenData {
			assignment_map, err := ggen.MakeAssignmentList()
			if err != nil {
				Logger.Fatalf("Could not generate assignment map: %s", err.Error())
			}

			data.AssignmentMap = assignment_map

			return data
		},
		func(data GenData) GenData {
			data.Fields = ggen.GetFields()

			return data
		},
	)
	if err != nil {
		Logger.Fatalf("Could not generate code: %s", err.Error())
	}
}

// templ is the template for the tree node.
const templ = `// Code generated by go generate; EDIT THIS FILE DIRECTLY

package {{ .PackageName }}

import (
	"slices"

	"github.com/PlayerR9/MyGoLib/ListLike/Stacker"
	"github.com/PlayerR9/MyGoLib/Units/common"
	
	{{- if ne .PackageName "treenode" }}
		"github.com/PlayerR9/treenode"
	{{- end }}
)

// {{ .IteratorName }} is a pull-based iterator that iterates
// over the children of a {{ .TypeName }}.
type {{ .IteratorName }}{{ .Generics }} struct {
	parent, current *{{ .TypeSig }}
}

// Consume implements the common.Iterater interface.
//
// *common.ErrExhaustedIter is the only error returned by this function and the returned
// node is never nil.
func (iter *{{ .IteratorSig }}) Consume() (*{{ .TypeSig }}, error) {
	if iter.current == nil {
		return nil, common.NewErrExhaustedIter()
	}

	node := iter.current
	iter.current = iter.current.NextSibling

	return node, nil
}

// Restart implements the common.Iterater interface.
func (iter *{{ .IteratorSig }}) Restart() {
	iter.current = iter.parent.FirstChild
}

// {{ .TypeName }} is a node in a tree.
type {{ .TypeName }}{{ .Generics }} struct {
	Parent, FirstChild, NextSibling, LastChild, PrevSibling *{{ .TypeSig }}

	{{- range $key, $value := .Fields }}
	{{ $key }} {{ $value }}
	{{- end }}
}

// Iterator implements the {{ .TypeSig }} interface.
//
// This function iterates over the children of the node, it is a pull-based iterator,
// and never returns nil.
func (tn *{{ .TypeSig }}) Iterator() common.Iterater[*{{ .TypeSig }}] {
	return &{{ .IteratorSig }}{
		parent: tn,
		current: tn.FirstChild,
	}
}

// String implements the {{ .TypeSig }} interface.
func (tn *{{ .TypeSig }}) String() string {
	// WARNING: Implement this function.
	str := common.StringOf(tn.Data)

	return str
}

// Copy implements the {{ .TypeSig }} interface.
//
// It never returns nil and it does not copy the parent or the sibling pointers.
func (tn *{{ .TypeSig }}) Copy() common.Copier {
	var child_copy []*{{ .TypeSig }}	

	for c := tn.FirstChild; c != nil; c = c.NextSibling {
		child_copy = append(child_copy, c.Copy().(*{{ .TypeSig }}))
	}

	// Copy here the data of the node.

	tn_copy := &{{ .TypeSig }}{
	 	// Add here the copied data of the node.
	}

	tn_copy.LinkChildren(child_copy)

	return tn_copy
}

// Cleanup implements the {{ .TypeSig }} interface.
//
// This is expensive as it has to traverse the whole tree to clean up the nodes, one
// by one. While this is useful for freeing up memory, for large enough trees, it is
// recommended to let the garbage collector handle the cleanup.
//
// Despite the above, this function does not use recursion and is safe to use (but
// make sure goroutines are not running on the tree while this function is called).
//
// Finally, it also logically removes the node from the siblings and the parent.
func (tn *{{ .TypeSig }}) Cleanup() {
	type Helper struct {
		previous, current *{{ .TypeSig }}
	}

	stack := Stacker.NewLinkedStack[*Helper]()

	// Free the first node.
	for c := tn.FirstChild; c != nil; c = c.NextSibling {
		h := &Helper{
			previous:	c.PrevSibling,
			current: 	c,
		}

		stack.Push(h)
	}

	tn.FirstChild = nil
	tn.LastChild = nil
	tn.Parent = nil

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

	prev := tn.PrevSibling
	next := tn.NextSibling

	if prev != nil {
		prev.NextSibling = next
	}

	if next != nil {
		next.PrevSibling = prev
	}

	tn.PrevSibling = nil
	tn.NextSibling = nil
}

// IsLeaf implements the {{ .TypeSig }} interface.
func (tn *{{ .TypeSig }}) IsLeaf() bool {
	return tn.FirstChild == nil
}

// IsSingleton implements the {{ .TypeSig }} interface.
func (tn *{{ .TypeSig }}) IsSingleton() bool {
	return tn.FirstChild != nil && tn.FirstChild == tn.LastChild
}

// DeleteChild removes the given child from the children of the node.
//
// Parameters:
//   - target: The child to remove.
//
// Returns:
//   - []*{{ .TypeSig }}: A slice of pointers to the children of the node. Nil if the node has no children.
//
// No nil nodes are returned.
func (tn *{{ .TypeSig }}) DeleteChild(target *{{ .TypeSig }}) []*{{ .TypeSig }} {
	if target == nil {
		return nil
	}

	children := tn.delete_child(target)

	if len(children) == 0 {
		return children
	}

	for _, child := range children {
		child.PrevSibling = nil
		child.NextSibling = nil
		child.Parent = nil
	}

	tn.FirstChild = nil
	tn.LastChild = nil

	return children
}

// Size implements the {{ .TypeSig }} interface.
//
// This is expensive as it has to traverse the whole tree to find the size of the tree.
// Thus, it is recommended to call this function once and then store the size somewhere if needed.
//
// Despite the above, this function does not use recursion and is safe to use.
//
// Finally, the traversal is done in a depth-first manner.
func (tn *{{ .TypeSig }}) Size() int {
	// It is safe to change the stack implementation as long as
	// it is not limited in size. If it is, make sure to check the error
	// returned by the Push and Pop methods.
	stack := Stacker.NewLinkedStack(tn)

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

{{- if eq (len .AssignmentMap) 0 }}

// New{{ .TypeName }} creates a new node with the given data.
//
// Returns:
//   - *{{ .TypeSig }}: A pointer to the newly created node. It is
//   never nil.
func New{{ .TypeName }}{{ .Generics }}() *{{ .TypeSig }} {
	return &{{ .TypeSig }}{}
}
	
{{- else }}

// New{{ .TypeName }} creates a new node with the given data.
//
// Parameters: {{- range $key, $value := .AssignmentMap }}
//   - {{ $key }}: The {{ $key }} of the node.
// {{- end }}
// Returns:
//   - *{{ .TypeSig }}: A pointer to the newly created node. It is
//   never nil.
func New{{ .TypeName }}{{ .Generics }}({{ .ParamList }}) *{{ .TypeSig }} {
	return &{{ .TypeSig }}{
		{{- range $key, $value := .AssignmentMap }}
		{{ $key }}: {{ $value }},
		{{- end }}
	}
}

{{- end }}

// LinkChildren links the parent with the children. It also links the children
// with each other. Nil children are ignored.
//
// Parameters:
//   - children: The children nodes.
func (tn *{{ .TypeSig }}) LinkChildren(children []*{{ .TypeSig }}) {
	if len(children) == 0 {
		return
	}

	var valid_children []*{{ .TypeSig }}

	for _, child := range children {
		if child == nil {
			continue
		}

		child.Parent = tn
		valid_children = append(valid_children, child)		
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

	tn.FirstChild, tn.LastChild = valid_children[0], valid_children[len(valid_children)-1]
}

// AddChild adds a new child to the node. If the child is nil it does nothing.
//
// Parameters:
//   - child: The child to add.
//
// This function clears the parent and sibling pointers of the child and so, it
// does not add relatives to the child.
func (tn *{{ .TypeSig }}) AddChild(child *{{ .TypeSig }}) {
	if child == nil {
		return
	}
	
	child.NextSibling = nil
	child.PrevSibling = nil

	last_child := tn.LastChild

	if last_child == nil {
		tn.FirstChild = child
	} else {
		last_child.NextSibling = child
		child.PrevSibling = last_child
	}

	child.Parent = tn
	tn.LastChild = child
}

// RemoveNode removes the node from the tree while shifting the children up one level to
// maintain the tree structure.
//
// Also, the returned children can be used to create a forest of trees if the root node
// is removed.
//
// Returns:
//   - []*{{ .TypeSig }}: A slice of pointers to the children of the node iff the node is the root.
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
func (tn *{{ .TypeSig }}) RemoveNode() []*{{ .TypeSig }} {
	prev := tn.PrevSibling
	next := tn.NextSibling
	parent := tn.Parent

	var sub_roots []*{{ .TypeSig }}

	if parent == nil {
		for c := tn.FirstChild; c != nil; c = c.NextSibling {
			sub_roots = append(sub_roots, c)
		}
	} else {
		children := parent.delete_child(tn)

		for _, child := range children {
			child.Parent = parent
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

	tn.Parent = nil
	tn.PrevSibling = nil
	tn.NextSibling = nil

	if len(sub_roots) == 0 {
		return sub_roots
	}

	for _, child := range sub_roots {
		child.PrevSibling = nil
		child.NextSibling = nil
		child.Parent = nil
	}

	tn.FirstChild = nil
	tn.LastChild = nil

	return sub_roots
}

// GetLeaves returns all the leaves of the tree rooted at the node.
//
// Returns:
//   - []*{{ .TypeSig }}: A slice of pointers to the leaves of the tree.
//
// This is expensive as leaves are not stored and so, every time this function is called,
// it has to do a DFS traversal to find the leaves. Thus, it is recommended to call
// this function once and then store the leaves somewhere if needed.
//
// Despite the above, this function does not use recursion and is safe to use.
//
// Finally, no nil nodes are returned.
func (tn *{{ .TypeSig }}) GetLeaves() []*{{ .TypeSig }} {
	// It is safe to change the stack implementation as long as
	// it is not limited in size. If it is, make sure to check the error
	// returned by the Push and Pop methods.
	stack := Stacker.NewLinkedStack(tn)

	var leaves []*{{ .TypeSig }}

	for {
		top, ok := stack.Pop()
		if !ok {
			break
		}

		if top.FirstChild == nil {
			leaves = append(leaves, top)
		} else {
			for c := top.FirstChild; c != nil; c = c.NextSibling {
				stack.Push(c)
			}
		}
	}

	return leaves
}

// GetAncestors returns all the ancestors of the node. This does not return the node itself.
//
// Returns:
//   - []*{{ .TypeSig }}: A slice of pointers to the ancestors of the node.
//
// The ancestors are returned in the opposite order of a DFS traversal. Therefore, the first element is the parent
// of the node.
//
// This is expensive since ancestors are not stored and so, every time this
// function is called, it has to traverse the tree to find the ancestors. Thus, it is
// recommended to call this function once and then store the ancestors somewhere if needed.
//
// Despite the above, this function does not use recursion and is safe to use.
//
// Finally, no nil nodes are returned.
func (tn *{{ .TypeSig }}) GetAncestors() []*{{ .TypeSig }} {
	var ancestors []*{{ .TypeSig }}

	for node := tn; node.Parent != nil; node = node.Parent {
		ancestors = append(ancestors, node.Parent)
	}

	slices.Reverse(ancestors)

	return ancestors
}

// GetLastSibling returns the last sibling of the node. If it has a parent,
// it returns the last child of the parent. Otherwise, it returns the last
// sibling of the node.
//
// As an edge case, if the node has no parent and no next sibling, it returns
// the node itself. Thus, this function never returns nil.
//
// Returns:
//   - *{{ .TypeSig }}: A pointer to the last sibling.
func (tn *{{ .TypeSig }}) GetLastSibling() *{{ .TypeSig }} {
	if tn.Parent != nil {
		return tn.Parent.LastChild
	} else if tn.NextSibling == nil {
		return tn
	}

	last_sibling := tn

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
//   - *{{ .TypeSig }}: A pointer to the first sibling.
func (tn *{{ .TypeSig }}) GetFirstSibling() *{{ .TypeSig }} {
	if tn.Parent != nil {
		return tn.Parent.FirstChild
	} else if tn.PrevSibling == nil {
		return tn
	}

	first_sibling := tn

	for first_sibling.PrevSibling != nil {
		first_sibling = first_sibling.PrevSibling
	}

	return first_sibling
}

// IsRoot returns true if the node does not have a parent.
//
// Returns:
//   - bool: True if the node is the root, false otherwise.
func (tn *{{ .TypeSig }}) IsRoot() bool {
	return tn.Parent == nil
}

// AddChildren is a convenience function to add multiple children to the node at once.
// It is more efficient than adding them one by one. Therefore, the behaviors are the
// same as the behaviors of the {{ .TypeName }}.AddChild function.
//
// Parameters:
//   - children: The children to add.
func (tn *{{ .TypeSig }}) AddChildren(children []*{{ .TypeSig }}) {
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

	last_child := tn.LastChild

	if last_child == nil {
		tn.FirstChild = first_child
	} else {
		last_child.NextSibling = first_child
		first_child.PrevSibling = last_child
	}

	first_child.Parent = tn
	tn.LastChild = first_child

	// Deal with the rest of the children
	for i := 1; i < len(children); i++ {
		child := children[i]

		child.NextSibling = nil
		child.PrevSibling = nil

		last_child := tn.LastChild
		last_child.NextSibling = child
		child.PrevSibling = last_child

		child.Parent = tn
		tn.LastChild = child
	}
}

// GetChildren returns the immediate children of the node.
//
// The returned nodes are never nil and are not copied. Thus, modifying the returned
// nodes will modify the tree.
//
// Returns:
//   - []*{{ .TypeSig }}: A slice of pointers to the children of the node.
func (tn *{{ .TypeSig }}) GetChildren() []*{{ .TypeSig }} {
	var children []*{{ .TypeSig }}

	for c := tn.FirstChild; c != nil; c = c.NextSibling {
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
func (tn *{{ .TypeSig }}) HasChild(target *{{ .TypeSig }}) bool {
	if target == nil || tn.FirstChild == nil {
		return false
	}

	for c := tn.FirstChild; c != nil; c = c.NextSibling {
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
//   - []{{ .TypeSig }}: A slice of pointers to the children of the node.
func (tn *{{ .TypeSig }}) delete_child(target *{{ .TypeSig }}) []*{{ .TypeSig }} {
	ok := tn.HasChild(target)
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

	if target == tn.FirstChild {
		tn.FirstChild = next

		if next == nil {
			tn.LastChild = nil
		}
	} else if target == tn.LastChild {
		tn.LastChild = prev
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
func (tn *{{ .TypeSig }}) IsChildOf(target *{{ .TypeSig }}) bool {
	if target == nil {
		return false
	}

	parents := target.GetAncestors()

	for node := tn; node.Parent != nil; node = node.Parent {
		ok := slices.Contains(parents, node.Parent)
		if ok {
			return true
		}
	}

	return false
}

/*

// FindCommonAncestor returns the first common ancestor of the two nodes.
//
// Parameters:
//   - n1: The first node.
//   - n2: The second node.
//
// Returns:
//   - *TreeNode[T]: A pointer to the common ancestor. Nil if no such node is found.
func FindCommonAncestor[T any](n1, n2 *TreeNode[T]) *TreeNode[T] {
	if n1 == nil {
		return n2
	} else if n2 == nil {
		return n1
	} else if n1 == n2 {
		return n1
	}

	ancestors1 := n1.GetAncestors()
	ancestors2 := n2.GetAncestors()

	if len(ancestors1) > len(ancestors2) {
		ancestors1, ancestors2 = ancestors2, ancestors1
	}

	for _, node := range ancestors1 {
		ok := slices.Contains(ancestors2, node)
		if ok {
			return node
		}
	}

	return nil
}

// FindBranchingPoint returns the first node in the path from n to the root
// such that has more than one sibling.
//
// Returns:
//   - *TreeNode[T]: The branching point.
//   - *TreeNode[T]: The parent of the branching point.
//   - bool: True if the node has a branching point, false otherwise.
//
// Behaviors:
//   - If there is no branching point, it returns the root of the tree. However,
//     if n is nil, it returns nil, nil, false and if the node has no parent, it
//     returns nil, n, false.
func FindBranchingPoint[T any](n *TreeNode[T]) (*TreeNode[T], *TreeNode[T], bool) {
	if n == nil {
		return nil, nil, false
	}

	parent := n.Parent
	if parent == nil {
		return nil, n, false
	}

	var has_branching_point bool

	for !has_branching_point {
		grand_parent := parent.Parent
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
*/
`
