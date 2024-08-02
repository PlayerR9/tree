package internal

import (
	"errors"
	"flag"
	"fmt"

	ggen "github.com/PlayerR9/lib_units/generator"
)

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

func Parse() (string, error) {
	err := ggen.ParseFlags()
	if err != nil {
		return "", fmt.Errorf("could not parse flags: %w", err)
	}

	if TypeNameFlag == nil {
		return "", errors.New("flag TypeNameFlag must be set")
	}

	type_name := *TypeNameFlag

	err = ggen.IsValidName(type_name, nil, ggen.Exported)
	if err != nil {
		return "", fmt.Errorf("invalid type name: %w", err)
	}

	return type_name, nil
}
