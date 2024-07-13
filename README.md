# treenode
A Go package used for generating first child/next sibling tree nodes. It also features some already generated tree nodes.


## Table of Contents

1. [Table of Contents](#table-of-contents)
2. [Tool](#tool)
   - [Installation](#installation)
   - [Usage](#usage)
3. [Documentation](#documentation)
4. [Content](#content)


## Tool

### Installation

To install the tool, run the following command:
```
go get -u github.com/PlayerR9/treenode/cmd/treenode
```


### Usage

Once imported, you can use the tool to generate tree nodes for your own types. Like so:
```go
import _ "github.com/PlayerR9/treenode"

//go:generate go run github.com/PlayerR9/treenode/cmd/treenode -type=Foo -fields=value/int
```

This will generate a tree node with the name "Foo" that contains, among other things, the field "value" of type "int."

The type generated will be in the same package as the tool. Make sure to read the documentation of the tool before using it.


## Documentation

```markdown
This command generates a tree node with the given fields that uses first child/next sibling pointers.

To use it, run the following command:

//go:generate go run github.com/PlayerR9/treenode -type=<type_name> -fields=<field_list> [ -g=<generics>] [ -output=<output_file> ]

**Flag: Type Name**

The "type name" flag is used to specify the name of the tree node struct. As such, it must be set and,
not only does it have to be a valid Go identifier, but it also must start with an upper case letter.

**Flag: Fields**

The "fields" flag is used to specify the fields that the tree node contains. Because it doesn't make
a lot of sense to have a tree node without fields, this flag must be set.

Its argument is specified as a list of key-value pairs where each pair is separated by a comma (",") and
a slash ("/") is used to separate the key and the value.

The key indicates the name of the field while the value indicates the type of the field.

For instance, running the following command:

//go:generate treenode -type="TreeNode" -fields=a/int,b/int,name/string

will generate a tree node with the following fields:

type TreeNode struct {
	// Node pointers.

	a int
	b int
	name string
}

It is important to note that spaces are not allowed.

Also, it is possible to specify generics by following the value with the generics between square brackets;
like so: "a/MyType[T,C]"

**Flag: Generics**

This optional flag is used to specify the type(s) of the generics. However, this only applies if at least one
generic type is specified in the fields flag. If none, then this flag is ignored.

As an edge case, if this flag is not specified but the fields flag contains generics, then
all generics are set to the default value of "any".

As with the fields flag, its argument is specified as a list of key-value pairs where each pair is separated
by a comma (",") and a slash ("/") is used to separate the key and the value. The key indicates the name of
the generic and the value indicates the type of the generic.

For instance, running the following command:

//go:generate treenode -type="TreeNode" -fields=a/MyType[T],b/MyType[C] -g=T/any,C/int

will generate a tree node with the following fields:

type TreeNode[T any, C int] struct {
	// Node pointers.

	a T
	b C
}

**Flag: Output File**

This optional flag is used to specify the output file. If not specified, the output will be written to
standard output, that is, the file "<type_name>_treenode.go" in the root of the current directory.
```


## Content

Here are all the pregenerated files:
- [bool](bool.go)
- [byte](byte.go)
- [int](int.go)
- [int8](int8.go)
- [int16](int16.go)
- [int32](int32.go)
- [int64](int64.go)
- [float32](float32.go)
- [float64](float64.go)
- [rune](rune.go)
- [string](string.go)
- [uint](uint.go)
- [uint8](uint8.go)
- [uint16](uint16.go)
- [uint32](uint32.go)
- [uint64](uint64.go)
- [uintptr](uintptr.go)
- [error](error.go)
- [complex128](complex128.go)
- [complex64](complex64.go)
- [generic](generic.go)
- [status](status.go)