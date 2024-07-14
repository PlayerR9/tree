package pkg

import (
	"flag"

	ggen "github.com/PlayerR9/MyGoLib/go_generator"
)

var (
	// TypeNameFlag is the flag for the type name.
	TypeNameFlag *string
)

func init() {
	TypeNameFlag = flag.String("type", "",
		"The type name to generate the tree node for. It must be set."+
			" Must start with an upper case letter and must be a valid Go identifier.",
	)

	ggen.SetOutputFlag("<type_name>_treenode.go", false)
	ggen.SetStructFieldsFlag("fields", true, -1, "The fields to generate the code for.")
	ggen.SetGenericsSignFlag("g", false, -1)
}
