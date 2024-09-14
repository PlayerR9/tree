package internal

import (
	"errors"
	"flag"
	"fmt"

	gcgen "github.com/PlayerR9/go-commons/generator"
)

var (
	// OutputFlag is the flag for the output.
	OutputFlag *gcgen.OutputLocVal

	// StructFieldsFlag is the flag for the struct fields.
	StructFieldsFlag *gcgen.StructFieldsVal

	// GenericsSignFlag is the flag for the generics sign.
	GenericsSignFlag *gcgen.GenericsSignVal

	// TypeNameFlag is the flag for the type name.
	TypeNameFlag *string
)

func init() {
	TypeNameFlag = flag.String("name", "",
		"The name of the struct to generate the tree node for. It must be set."+
			" Must start with an upper case letter and must be a valid Go identifier.",
	)

	OutputFlag = gcgen.NewOutputFlag("<type_name>_treenode.go", false)
	StructFieldsFlag = gcgen.NewStructFieldsFlag("fields", true, -1, "The fields to generate the code for.")
	GenericsSignFlag = gcgen.NewGenericsSignFlag("g", false, -1)
}

func Parse() (string, error) {
	gcgen.ParseFlags()

	if TypeNameFlag == nil {
		return "", errors.New("flag TypeNameFlag must be set")
	}

	type_name := *TypeNameFlag

	err := gcgen.IsValidVariableName(type_name, nil, gcgen.Exported)
	if err != nil {
		return "", fmt.Errorf("invalid type name: %w", err)
	}

	return type_name, nil
}
