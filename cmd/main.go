// This command generates a tree node with the given fields that uses first child/next sibling pointers.
//
// To use it, run the following command:
//
// //go:generate go run github.com/PlayerR9/tree/cmd/tree -name=<type_name> -fields=<field_list> [ -g=<generics>] [ -o=<output_file> ]
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
//	//go:generate tree -name=TreeNode -fields=a/int,b/int,name/string
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
//	//go:generate tree -name=TreeNode -fields=a/MyType[T],b/MyType[C] -g=T/any,C/int
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

	ggen "github.com/PlayerR9/MyGoLib/Generator"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

var (
	// Logger is the logger for the package.
	Logger *log.Logger

	// t is the template for the tree node.
	t *template.Template
)

func init() {
	Logger = ggen.InitLogger("tree")

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

	type_name := uc.AssertDerefNil(TypeNameFlag, "TypeNameFlag")

	err = ggen.IsValidName(type_name, nil, ggen.Exported)
	if err != nil {
		Logger.Fatalf("Invalid type name: %s", err.Error())
	}

	filename, err := ggen.FixOutputLoc(type_name, "_treenode.go")
	if err != nil {
		Logger.Fatalf("Could not fix output location: %s", err.Error())
	}

	err = ggen.Generate(filename, GenData{}, t,
		func(data *GenData) error {
			data.TypeName = type_name

			return nil
		},
		func(data *GenData) error {
			tn_type_sig, err := ggen.MakeTypeSig(type_name, "")
			if err != nil {
				return err
			}

			data.TypeSig = tn_type_sig

			return nil
		},
		func(data *GenData) error {
			tn_iterator_sig, err := ggen.MakeTypeSig(type_name, "Iterator")
			if err != nil {
				return err
			}

			data.IteratorSig = tn_iterator_sig

			return nil
		},
		func(data *GenData) error {
			data.Generics = ggen.GenericsSigFlag.String()

			return nil
		},
		func(data *GenData) error {
			var builder strings.Builder

			builder.WriteString(type_name)
			builder.WriteString("Iterator")

			data.IteratorName = builder.String()

			return nil
		},
		func(data *GenData) error {
			param_list, err := ggen.MakeParameterList()
			if err != nil {
				return err
			}

			data.ParamList = param_list

			return nil
		},
		func(data *GenData) error {
			assignment_map, err := ggen.MakeAssignmentList()
			if err != nil {
				return err
			}

			data.AssignmentMap = assignment_map

			return nil
		},
		func(data *GenData) error {
			data.Fields = ggen.StructFieldsFlag.GetFields()

			return nil
		},
	)
	if err != nil {
		Logger.Fatalf("Could not generate code: %s", err.Error())
	}

	Logger.Printf("Generated %s", filename)
}

// templ is the template for the tree node.
const templ = `// Code generated by go generate; EDIT THIS FILE DIRECTLY

package {{ .PackageName }}

import (
	"slices"

	"github.com/PlayerR9/stack"
	"github.com/PlayerR9/MyGoLib/Units/common"
	"github.com/PlayerR9/tree/tree"
)

// {{ .IteratorName }} is a pull-based iterator that iterates
// over the children of a {{ .TypeName }}.
type {{ .IteratorName }}{{ .Generics }} struct {
	parent, current *{{ .TypeSig }}
}

// Consume implements the common.Iterater interface.
//
// The only error type that can be returned by this function is the *common.ErrExhaustedIter type.
//
// Moreover, the return value is always of type *{{ .TypeSig }} and never nil; unless the iterator
// has reached the end of the branch.
func (iter *{{ .IteratorSig }}) Consume() (tree.Noder, error) {
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

// IsLeaf implements the tree.Noder interface.
func (tn *{{ .TypeSig }}) IsLeaf() bool {
	return tn.FirstChild == nil
}

// GetParent implements the tree.Noder interface.
func (tn *{{ .TypeSig }}) GetParent() tree.Noder {
	return tn.Parent
}

// IsSingleton implements the tree.Noder interface.
func (tn *{{ .TypeSig }}) IsSingleton() bool {
	return tn.FirstChild != nil && tn.FirstChild == tn.LastChild
}


// DeleteChild implements the tree.Noder interface.
func (tn *{{ .TypeSig }}) DeleteChild(target tree.Noder) []tree.Noder {
	if target == nil {
		return nil
	}

	tmp, ok := target.(*{{ .TypeSig }})
	if !ok {
		return nil
	}

	children := tn.delete_child(tmp)

	if len(children) == 0 {
		return nil
	}

	for _, child := range children {
		child.PrevSibling = nil
		child.NextSibling = nil
		child.Parent = nil
	}

	tn.FirstChild = nil
	tn.LastChild = nil

	conv := make([]tree.Noder, 0, len(children))

	for _, child := range children {
		conv = append(conv, child)
	}

	return conv
}

// GetFirstChild implements the tree.Noder interface.
func (tn *{{ .TypeSig }}) GetFirstChild() tree.Noder {
	return tn.FirstChild
}

// AddChild implements the tree.Noder interface.
func (tn *{{ .TypeSig }}) AddChild(target tree.Noder) {
	if target == nil {
		return
	}

	tmp, ok := target.(*{{ .TypeSig }})
	if !ok {
		return
	}
	
	tmp.NextSibling = nil
	tmp.PrevSibling = nil

	last_child := tn.LastChild

	if last_child == nil {
		tn.FirstChild = tmp
	} else {
		last_child.NextSibling = tmp
		tmp.PrevSibling = last_child
	}

	tmp.Parent = tn
	tn.LastChild = tmp
}

// LinkChildren implements the tree.Noder interface.
func (tn *{{ .TypeSig }}) LinkChildren(children []tree.Noder) {
	var valid_children []*{{ .TypeSig }}

	for _, child := range children {
		if child == nil {
			continue
		}

		tmp, ok := child.(*{{ .TypeSig }})
		if !ok {
			continue
		}

		tmp.Parent = tn

		valid_children = append(valid_children, tmp)
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

// delete_child is a helper function to delete the child from the children of the node. No nil
// nodes are returned when this function is called. However, if target is nil, then nothing happens.
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

// RemoveNode implements the tree.Noder interface.
func (tn *{{ .TypeSig }}) RemoveNode() []tree.Noder {
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
		return nil
	}

	for _, child := range sub_roots {
		child.PrevSibling = nil
		child.NextSibling = nil
		child.Parent = nil
	}

	tn.FirstChild = nil
	tn.LastChild = nil

	conv := make([]tree.Noder, 0, len(sub_roots))
	for _, child := range sub_roots {
		conv = append(conv, child)
	}

	return conv
}

// Copy implements the tree.Noder interface.
//
// Although this function never returns nil, it does not copy the parent nor the sibling pointers.
func (tn *{{ .TypeSig }}) Copy() common.Copier {
	var child_copy []tree.Noder

	for c := tn.FirstChild; c != nil; c = c.NextSibling {
		child_copy = append(child_copy, c.Copy().(tree.Noder))
	}

	// Copy here the data of the node.

	tn_copy := &{{ .TypeSig }}{
	 	// Add here the copied data of the node.
	}

	tn_copy.LinkChildren(child_copy)

	return tn_copy
}

// Iterator implements the tree.Noder interface.
//
// This function returns an iterator that iterates over the direct children of the node.
// Implemented as a pull-based iterator, this function never returns nil and any of the
// values is guaranteed to be a non-nil node of type {{ .TypeSig }}.
func (tn *{{ .TypeSig }}) Iterator() common.Iterater[tree.Noder] {
	return &{{ .IteratorSig }}{
		parent: tn,
		current: tn.FirstChild,
	}
}

// Cleanup implements the tree.Noder interface.
//
// This function is expensive as it has to traverse the whole tree to clean up the nodes, one
// by one. While this function is useful for freeing up memory, for large enough trees, it is
// recommended to let the garbage collector handle the cleanup.
//
// Despite the above, this function does not use recursion but it is not safe to use in goroutines
// as pointers may be dereferenced while another goroutine is still using them.
//
// Finally, this function also logically removes the node from the siblings and the parent.
func (tn *{{ .TypeSig }}) Cleanup() {
	type Helper struct {
		previous, current *{{ .TypeSig }}
	}

	lls := stack.NewLinkedStack[*Helper]()

	// Free the first node.
	for c := tn.FirstChild; c != nil; c = c.NextSibling {
		h := &Helper{
			previous:	c.PrevSibling,
			current: 	c,
		}

		lls.Push(h)
	}

	tn.FirstChild = nil
	tn.LastChild = nil
	tn.Parent = nil

	// Free the rest of the nodes.
	for {
		h, ok := lls.Pop()
		if !ok {
			break
		}

		for c := h.current.FirstChild; c != nil; c = c.NextSibling {
			h := &Helper{
				previous:	c.PrevSibling,
				current: 	c,
			}

			lls.Push(h)
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

// String implements the tree.Noder interface.
func (tn *{{ .TypeSig }}) String() string {
	// WARNING: Implement this function.
	str := common.StringOf(tn.Data)

	return str
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

	parents := tree.GetNodeAncestors(target)

	for node := tn; node.Parent != nil; node = node.Parent {
		ok := slices.Contains(parents, node.Parent)
		if ok {
			return true
		}
	}

	return false
}
`
