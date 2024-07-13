// This command generates a tree node with the given fields that uses first child/next sibling pointers.
//
// To use it, run the following command:
//
// //go:generate go run github.com/PlayerR9/treenode/cmd/treenode -type=<type_name> -fields=<field_list> [ -g=<generics>] [ -output=<output_file> ]
//
// **Flag: Type Name**
//
// The "type name" flag is used to specify the name of the tree node struct. As such, it must be set and,
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
//	//go:generate treenode -type="TreeNode" -fields=a/int,b/int,name/string
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
//	//go:generate treenode -type="TreeNode" -fields=a/MyType[T],b/MyType[C] -g=T/any,C/int
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
	"log"
	"os"

	pkg "github.com/PlayerR9/treenode/cmd/treenode/pkg"
)

var (
	// Logger is the logger for the package.
	Logger *log.Logger
)

func init() {
	Logger = log.New(os.Stderr, "[treenode]: ", log.LstdFlags)
}

func main() {
	type_name, filename, err := pkg.ParseFlags()
	if err != nil {
		Logger.Fatalf("Could not parse flags: %s", err.Error())
	}

	// Generate the code.
	g, err := pkg.NewGenerator(type_name, filename)
	if err != nil {
		Logger.Fatal(err.Error())
	}

	generated_data, err := g.Generate()
	if err != nil {
		Logger.Fatalf("Could not generate data: %s", err.Error())
	}

	// Write the code to the file.

	err = os.WriteFile(filename, generated_data, 0644)
	if err != nil {
		Logger.Fatalf("Could not write file: %s", err.Error())
	}
}
