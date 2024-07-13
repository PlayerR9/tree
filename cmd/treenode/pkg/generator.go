package pkg

import (
	"bytes"
	"errors"
	"fmt"
	"go/build"
	"path/filepath"
	"slices"
	"strings"
	"text/template"
	"unicode"

	uc "github.com/PlayerR9/MyGoLib/Units/common"
)

// generator is a code generator for a tree node.
type generator struct {
	// package_name is the name of the package.
	package_name string

	// type_name is the name of the type.
	type_name string

	// variable_name is the name of the variable.
	variable_name string

	// fields is the fields to generate the code for.
	fields map[string]string
}

func make_output_type(type_name string, suffix string, gv *GenericsValue) string {
	uc.AssertParam("gv", gv != nil, errors.New("generics must be set"))
	uc.AssertParam("type_name", type_name != "", errors.New("type_name must be set"))

	var builder strings.Builder

	builder.WriteString(type_name)

	if suffix != "" {
		builder.WriteString(suffix)
	}

	if len(gv.letters) > 0 {
		str := gv.GetGenericsList()

		builder.WriteString(str)
	}

	output_name := builder.String()

	return output_name
}

func make_parameter_list(fields map[string]string) (string, error) {
	var field_list []string
	var type_list []string

	for k, v := range fields {
		if k == "" {
			err := errors.New("found type name with empty name")
			return "", err
		}

		first_letter := rune(k[0])

		ok := unicode.IsLetter(first_letter)
		if !ok {
			err := fmt.Errorf("type name %q must start with a letter", k)
			return "", err
		}

		ok = unicode.IsUpper(first_letter)
		if !ok {
			continue
		}

		pos, ok := slices.BinarySearch(field_list, k)
		uc.AssertF(!ok, "%q must be unique", k)

		field_list = slices.Insert(field_list, pos, k)
		type_list = slices.Insert(type_list, pos, v)
	}

	var values []string
	var builder strings.Builder

	for i := 0; i < len(field_list); i++ {
		param := strings.ToLower(field_list[i])

		builder.WriteString(param)
		builder.WriteRune(' ')
		builder.WriteString(type_list[i])

		str := builder.String()
		values = append(values, str)

		builder.Reset()
	}

	joined_str := strings.Join(values, ", ")

	return joined_str, nil
}

func make_assignment_list(fields map[string]string) (map[string]string, error) {
	var field_list []string
	var type_list []string

	for k, v := range fields {
		if k == "" {
			err := errors.New("found type name with empty name")
			return nil, err
		}

		first_letter := rune(k[0])

		ok := unicode.IsLetter(first_letter)
		if !ok {
			err := fmt.Errorf("type name %q must start with a letter", k)
			return nil, err
		}

		ok = unicode.IsUpper(first_letter)
		if !ok {
			continue
		}

		pos, ok := slices.BinarySearch(field_list, k)
		uc.AssertF(!ok, "%q must be unique", k)

		field_list = slices.Insert(field_list, pos, k)
		type_list = slices.Insert(type_list, pos, v)
	}

	assignment_map := make(map[string]string)

	for i := 0; i < len(field_list); i++ {
		param := strings.ToLower(field_list[i])

		assignment_map[field_list[i]] = param
	}

	return assignment_map, nil
}

// Generate generates the code for the tree node.
//
// Returns:
//   - []byte: The generated code.
//   - error: An error if the code could not be generated.
func (g *generator) Generate() ([]byte, error) {
	uc.Assert(g.package_name != "", "package name must be set")
	uc.Assert(g.type_name != "", "type name must be set")
	uc.Assert(g.variable_name != "", "variable name must be set")

	t := template.Must(
		template.New("").Parse(templ),
	)

	type GenData struct {
		PackageName   string
		TypeName      string
		TypeSig       string
		VariableName  string
		Fields        map[string]string
		IteratorName  string
		IteratorSig   string
		ParamList     string
		AssignmentMap map[string]string
		Generics      string
	}

	tn_type_sig := make_output_type(g.type_name, "", GenericsFlag)
	tn_iterator_sig := make_output_type(g.type_name, "Iterator", GenericsFlag)

	param_list, err := make_parameter_list(g.fields)
	if err != nil {
		err := fmt.Errorf("could not generate parameter list: %w", err)
		return nil, err
	}

	assignment_map, err := make_assignment_list(g.fields)
	if err != nil {
		err := fmt.Errorf("could not generate assignment map: %w", err)
		return nil, err
	}

	var generics string

	ok := GenericsFlag.HasGenerics()
	if ok {
		generics = GenericsFlag.String()
	}

	data := GenData{
		PackageName:   g.package_name,
		TypeName:      g.type_name,
		TypeSig:       tn_type_sig,
		VariableName:  g.variable_name,
		Fields:        g.fields,
		IteratorName:  g.type_name + "Iterator",
		IteratorSig:   tn_iterator_sig,
		ParamList:     param_list,
		AssignmentMap: assignment_map,
		Generics:      generics,
	}

	var buf bytes.Buffer

	err = t.Execute(&buf, data)
	if err != nil {
		return nil, err
	}

	result := buf.Bytes()

	return result, nil
}

func NewGenerator(type_name, var_name string, dest string) (*generator, error) {
	if dest != "" {
		left := filepath.Dir(dest)
		_, right := filepath.Split(left)
		dest = right
	}

	var pkg_name string

	if dest == "" || dest == "." {
		pkg, err := build.Default.ImportDir(".", 0)
		if err != nil {
			err := fmt.Errorf("could not import directory: %w", err)
			return nil, err
		}

		pkg_name = pkg.Name
	} else {
		pkg_name = dest
	}

	g := &generator{
		package_name:  pkg_name,
		type_name:     type_name,
		variable_name: var_name,
		fields:        FieldsFlag.fields,
	}

	return g, nil
}

// templ is the template for the tree node.
const templ = `// Code generated by go generate; EDIT THIS FILE DIRECTLY

package {{ .PackageName }}

import (
	"slices"
	"fmt"

	lls "github.com/PlayerR9/MyGoLib/ListLike/Stacker"
	uc "github.com/PlayerR9/MyGoLib/Units/common"
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
func (iter *{{ .IteratorSig }}) Consume() (Noder, error) {
	if iter.current == nil {
		return nil, uc.NewErrExhaustedIter()
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

// Iterator implements the Tree.Noder interface.
//
// This function iterates over the children of the node, it is a pull-based iterator,
// and never returns nil.
func ({{ .VariableName }} *{{ .TypeSig }}) Iterator() uc.Iterater[Noder] {
	return &{{ .IteratorSig }}{
		parent: {{ .VariableName }},
		current: {{ .VariableName }}.FirstChild,
	}
}

// String implements the Noder interface.
func ({{ .VariableName }} *{{ .TypeSig }}) String() string {
	// WARNING: Implement this function.
	str := fmt.Sprintf("%v", {{ .VariableName }}.Data)

	return str
}

// Copy implements the Noder interface.
//
// It never returns nil and it does not copy the parent or the sibling pointers.
func ({{ .VariableName }} *{{ .TypeSig }}) Copy() uc.Copier {
	var child_copy []Noder	

	for c := {{ .VariableName }}.FirstChild; c != nil; c = c.NextSibling {
		child_copy = append(child_copy, c.Copy().(Noder))
	}

	// Copy here the data of the node.

	{{ .VariableName}}_copy := &{{ .TypeSig }}{
	 	// Add here the copied data of the node.
	}

	{{ .VariableName }}_copy.LinkChildren(child_copy)

	return {{ .VariableName }}_copy
}

// SetParent implements the Noder interface.
func ({{ .VariableName }} *{{ .TypeSig }}) SetParent(parent Noder) bool {
	if parent == nil {
		{{ .VariableName }}.Parent = nil
		return true
	}

	p, ok := parent.(*{{ .TypeSig }})
	if !ok {
		return false
	}

	{{ .VariableName }}.Parent = p

	return true
}

// GetParent implements the Noder interface.
func ({{ .VariableName }} *{{ .TypeSig }}) GetParent() Noder {
	return {{ .VariableName }}.Parent
}

// LinkWithParent implements the Noder interface.
//
// Children that are not of type *{{ .TypeSig }} or nil are ignored.
func ({{ .VariableName }} *{{ .TypeSig }}) LinkChildren(children []Noder) {
	if len(children) == 0 {
		return
	}

	var valid_children []*{{ .TypeSig }}

	for _, child := range children {
		if child == nil {
			continue
		}

		c, ok := child.(*{{ .TypeSig }})
		if ok {
			c.Parent = {{ .VariableName }}
			valid_children = append(valid_children, c)
		}		
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

	{{ .VariableName }}.FirstChild, {{ .VariableName }}.LastChild = valid_children[0], valid_children[len(valid_children)-1]
}

// GetLeaves implements the Noder interface.
//
// This is expensive as leaves are not stored and so, every time this function is called,
// it has to do a DFS traversal to find the leaves. Thus, it is recommended to call
// this function once and then store the leaves somewhere if needed.
//
// Despite the above, this function does not use recursion and is safe to use.
//
// Finally, no nil nodes are returned.
func ({{ .VariableName }} *{{ .TypeSig }}) GetLeaves() []Noder {
	// It is safe to change the stack implementation as long as
	// it is not limited in size. If it is, make sure to check the error
	// returned by the Push and Pop methods.
	stack := lls.NewLinkedStack[Noder]({{ .VariableName }})

	var leaves []Noder

	for {
		top, ok := stack.Pop()
		if !ok {
			break
		}

		node := top.(*{{ .TypeSig }})
		if node.FirstChild == nil {
			leaves = append(leaves, top)
		} else {
			for c := node.FirstChild; c != nil; c = c.NextSibling {
				stack.Push(c)
			}
		}
	}

	return leaves
}

// Cleanup implements the Noder interface.
//
// This is expensive as it has to traverse the whole tree to clean up the nodes, one
// by one. While this is useful for freeing up memory, for large enough trees, it is
// recommended to let the garbage collector handle the cleanup.
//
// Despite the above, this function does not use recursion and is safe to use (but
// make sure goroutines are not running on the tree while this function is called).
//
// Finally, it also logically removes the node from the siblings and the parent.
func ({{ .VariableName }} *{{ .TypeSig }}) Cleanup() {
	type Helper struct {
		previous, current *{{ .TypeSig }}
	}

	stack := lls.NewLinkedStack[*Helper]()

	// Free the first node.
	for c := {{ .VariableName }}.FirstChild; c != nil; c = c.NextSibling {
		h := &Helper{
			previous:	c.PrevSibling,
			current: 	c,
		}

		stack.Push(h)
	}

	{{ .VariableName }}.FirstChild = nil
	{{ .VariableName }}.LastChild = nil
	{{ .VariableName }}.Parent = nil

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

	prev := {{ .VariableName }}.PrevSibling
	next := {{ .VariableName }}.NextSibling

	if prev != nil {
		prev.NextSibling = next
	}

	if next != nil {
		next.PrevSibling = prev
	}

	{{ .VariableName }}.PrevSibling = nil
	{{ .VariableName }}.NextSibling = nil
}

// GetAncestors implements the Noder interface.
//
// This is expensive since ancestors are not stored and so, every time this
// function is called, it has to traverse the tree to find the ancestors. Thus, it is
// recommended to call this function once and then store the ancestors somewhere if needed.
//
// Despite the above, this function does not use recursion and is safe to use.
//
// Finally, no nil nodes are returned.
func ({{ .VariableName }} *{{ .TypeSig }}) GetAncestors() []Noder {
	var ancestors []Noder

	for node := {{ .VariableName }}; node.Parent != nil; node = node.Parent {
		ancestors = append(ancestors, node.Parent)
	}

	slices.Reverse(ancestors)

	return ancestors
}

// IsLeaf implements the Noder interface.
func ({{ .VariableName }} *{{ .TypeSig }}) IsLeaf() bool {
	return {{ .VariableName }}.FirstChild == nil
}

// IsSingleton implements the Noder interface.
func ({{ .VariableName }} *{{ .TypeSig }}) IsSingleton() bool {
	return {{ .VariableName }}.FirstChild != nil && {{ .VariableName }}.FirstChild == {{ .VariableName }}.LastChild
}

// GetFirstChild implements the Noder interface.
func ({{ .VariableName }} *{{ .TypeSig }}) GetFirstChild() Noder {
	return {{ .VariableName }}.FirstChild
}

// DeleteChild implements the Noder interface.
//
// No nil nodes are returned.
func ({{ .VariableName }} *{{ .TypeSig }}) DeleteChild(target Noder) []Noder {
	if target == nil {
		return nil
	}

	n, ok := target.(*{{ .TypeSig }})
	if !ok {
		return nil
	}

	children := {{ .VariableName }}.delete_child(n)

	if len(children) == 0 {
		return children
	}

	for _, child := range children {
		c := child.(*{{ .TypeSig }})

		c.PrevSibling = nil
		c.NextSibling = nil
		c.Parent = nil
	}

	{{ .VariableName }}.FirstChild = nil
	{{ .VariableName }}.LastChild = nil

	return children
}

// Size implements the Noder interface.
//
// This is expensive as it has to traverse the whole tree to find the size of the tree.
// Thus, it is recommended to call this function once and then store the size somewhere if needed.
//
// Despite the above, this function does not use recursion and is safe to use.
//
// Finally, the traversal is done in a depth-first manner.
func ({{ .VariableName }} *{{ .TypeSig }}) Size() int {
	// It is safe to change the stack implementation as long as
	// it is not limited in size. If it is, make sure to check the error
	// returned by the Push and Pop methods.
	stack := lls.NewLinkedStack({{ .VariableName }})

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

// AddChild adds a new child to the node. If the child is nil or it is not of type
// *{{ .TypeSig }}, it does nothing.
//
// This function clears the parent and sibling pointers of the child and so, it
// does not add relatives to the child.
//
// Parameters:
//   - child: The child to add.
func ({{ .VariableName }} *{{ .TypeSig }}) AddChild(child Noder) {
	if child == nil {
		return
	}

	c, ok := child.(*{{ .TypeSig }})
	if !ok {
		return
	}
	
	c.NextSibling = nil
	c.PrevSibling = nil

	last_child := {{ .VariableName }}.LastChild

	if last_child == nil {
		{{ .VariableName }}.FirstChild = c
	} else {
		last_child.NextSibling = c
		c.PrevSibling = last_child
	}

	c.Parent = {{ .VariableName }}
	{{ .VariableName }}.LastChild = c
}

// RemoveNode removes the node from the tree while shifting the children up one level to
// maintain the tree structure.
//
// Also, the returned children can be used to create a forest of trees if the root node
// is removed.
//
// Returns:
//   - []Noder: A slice of pointers to the children of the node iff the node is the root.
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
func ({{ .VariableName }} *{{ .TypeSig }}) RemoveNode() []Noder {
	prev := {{ .VariableName }}.PrevSibling
	next := {{ .VariableName }}.NextSibling
	parent := {{ .VariableName }}.Parent

	var sub_roots []Noder

	if parent == nil {
		for c := {{ .VariableName }}.FirstChild; c != nil; c = c.NextSibling {
			sub_roots = append(sub_roots, c)
		}
	} else {
		children := parent.delete_child({{ .VariableName }})

		for _, child := range children {
			child.SetParent(parent)
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

	{{ .VariableName }}.Parent = nil
	{{ .VariableName }}.PrevSibling = nil
	{{ .VariableName }}.NextSibling = nil

	if len(sub_roots) == 0 {
		return sub_roots
	}

	for _, child := range sub_roots {
		c := child.(*{{ .TypeSig }})

		c.PrevSibling = nil
		c.NextSibling = nil
		c.Parent = nil
	}

	{{ .VariableName }}.FirstChild = nil
	{{ .VariableName }}.LastChild = nil

	return sub_roots
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
func ({{ .VariableName }} *{{ .TypeSig }}) GetLastSibling() *{{ .TypeSig }} {
	if {{ .VariableName }}.Parent != nil {
		return {{ .VariableName }}.Parent.LastChild
	} else if {{ .VariableName }}.NextSibling == nil {
		return {{ .VariableName }}
	}

	last_sibling := {{ .VariableName }}

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
func ({{ .VariableName }} *{{ .TypeSig }}) GetFirstSibling() *{{ .TypeSig }} {
	if {{ .VariableName }}.Parent != nil {
		return {{ .VariableName }}.Parent.FirstChild
	} else if {{ .VariableName }}.PrevSibling == nil {
		return {{ .VariableName }}
	}

	first_sibling := {{ .VariableName }}

	for first_sibling.PrevSibling != nil {
		first_sibling = first_sibling.PrevSibling
	}

	return first_sibling
}

// IsRoot returns true if the node does not have a parent.
//
// Returns:
//   - bool: True if the node is the root, false otherwise.
func ({{ .VariableName }} *{{ .TypeSig }}) IsRoot() bool {
	return {{ .VariableName }}.Parent == nil
}

// AddChildren is a convenience function to add multiple children to the node at once.
// It is more efficient than adding them one by one. Therefore, the behaviors are the
// same as the behaviors of the {{ .TypeName }}.AddChild function.
//
// Parameters:
//   - children: The children to add.
func ({{ .VariableName }} *{{ .TypeSig }}) AddChildren(children []*{{ .TypeSig }}) {
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

	last_child := {{ .VariableName }}.LastChild

	if last_child == nil {
		{{ .VariableName }}.FirstChild = first_child
	} else {
		last_child.NextSibling = first_child
		first_child.PrevSibling = last_child
	}

	first_child.Parent = {{ .VariableName }}
	{{ .VariableName }}.LastChild = first_child

	// Deal with the rest of the children
	for i := 1; i < len(children); i++ {
		child := children[i]

		child.NextSibling = nil
		child.PrevSibling = nil

		last_child := {{ .VariableName }}.LastChild
		last_child.NextSibling = child
		child.PrevSibling = last_child

		child.Parent = {{ .VariableName }}
		{{ .VariableName }}.LastChild = child
	}
}

// GetChildren returns the immediate children of the node.
//
// The returned nodes are never nil and are not copied. Thus, modifying the returned
// nodes will modify the tree.
//
// Returns:
//   - []Noder: A slice of pointers to the children of the node.
func ({{ .VariableName }} *{{ .TypeSig }}) GetChildren() []Noder {
	var children []Noder

	for c := {{ .VariableName }}.FirstChild; c != nil; c = c.NextSibling {
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
func ({{ .VariableName }} *{{ .TypeSig }}) HasChild(target *{{ .TypeSig }}) bool {
	if target == nil || {{ .VariableName }}.FirstChild == nil {
		return false
	}

	for c := {{ .VariableName }}.FirstChild; c != nil; c = c.NextSibling {
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
//   - []Noder: A slice of pointers to the children of the node.
func ({{ .VariableName }} *{{ .TypeSig }}) delete_child(target *{{ .TypeSig }}) []Noder {
	ok := {{ .VariableName }}.HasChild(target)
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

	if target == {{ .VariableName }}.FirstChild {
		{{ .VariableName }}.FirstChild = next

		if next == nil {
			{{ .VariableName }}.LastChild = nil
		}
	} else if target == {{ .VariableName }}.LastChild {
		{{ .VariableName }}.LastChild = prev
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
func ({{ .VariableName }} *{{ .TypeSig }}) IsChildOf(target *{{ .TypeSig }}) bool {
	if target == nil {
		return false
	}

	parents := target.GetAncestors()

	for node := {{ .VariableName }}; node.Parent != nil; node = node.Parent {
		parent := Noder(node.Parent)

		ok := slices.Contains(parents, parent)
		if ok {
			return true
		}
	}

	return false
}
`
