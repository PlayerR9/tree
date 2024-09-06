package internal

import (
	"errors"
	"flag"
	"fmt"

	ggen "github.com/PlayerR9/go-generator/generator"
)

var (
	// OutputFlag is the flag for the output.
	OutputFlag *ggen.OutputLocVal

	// StructFieldsFlag is the flag for the struct fields.
	StructFieldsFlag *ggen.StructFieldsVal

	// GenericsSignFlag is the flag for the generics sign.
	GenericsSignFlag *ggen.GenericsSignVal

	// TypeNameFlag is the flag for the type name.
	TypeNameFlag *string
)

func init() {
	TypeNameFlag = flag.String("name", "",
		"The name of the struct to generate the tree node for. It must be set."+
			" Must start with an upper case letter and must be a valid Go identifier.",
	)

	OutputFlag = ggen.NewOutputFlag("<type_name>_treenode.go", false)
	StructFieldsFlag = ggen.NewStructFieldsFlag("fields", true, -1, "The fields to generate the code for.")
	GenericsSignFlag = ggen.NewGenericsSignFlag("g", false, -1)
}

func Parse() (string, error) {
	ggen.ParseFlags()

	if TypeNameFlag == nil {
		return "", errors.New("flag TypeNameFlag must be set")
	}

	type_name := *TypeNameFlag

	err := ggen.IsValidVariableName(type_name, nil, ggen.Exported)
	if err != nil {
		return "", fmt.Errorf("invalid type name: %w", err)
	}

	return type_name, nil
}
