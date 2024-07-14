package pkg

import (
	"errors"
	"flag"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"unicode"

	uc "github.com/PlayerR9/MyGoLib/Units/common"
	utgo "github.com/PlayerR9/MyGoLib/Utility/Go"
	// Add the following when MyGoLib is updated
	// ggen "github.com/PlayerR9/MyGoLib/go_generator"
)

var (
	// TypeNameFlag is the flag for the type name.
	TypeNameFlag *string

	// Remove this once MyGoLib is updated
	// OutputFileFlag is the flag for the output file.
	OutputFileFlag *string

	// FieldsFlag is the flag for the data file.
	FieldsFlag *FieldsValue

	// GenericsFlag is the flag for the generics.
	GenericsFlag *GenericsValue
)

func init() {
	TypeNameFlag = flag.String("type", "",
		"The type name to generate the tree node for. It must be set."+
			" Must start with an upper case letter and must be a valid Go identifier.",
	)

	// Remove this once MyGoLib is updated
	OutputFileFlag = flag.String("output", "",
		"The output file to write the generated code to. If not set, the default file name is used."+
			" That is \"<type_name>_treenode.go\".",
	)

	// Add this once MyGoLib is updated
	// ggen.SetOutputFlag("<type_name>_treenode.go", false)

	FieldsFlag = NewFieldsValue()

	flag.Var(FieldsFlag, "fields",
		"The fields to generate the code for. It must be set."+
			" The syntax of the field's argument is described in the documentation.",
	)

	GenericsFlag = NewGenericsValue()

	flag.Var(GenericsFlag, "g",
		"The generics to generate the code for. It is optional."+
			" The syntax of the generics's argument is described in the documentation.",
	)
}

func ParseFlags() (string, string, error) {
	// Remove these assertions once MyGoLib is updated
	uc.Assert(TypeNameFlag != nil, "TypeNameFlag must not be nil")
	uc.Assert(FieldsFlag != nil, "FieldsFlag must not be nil")
	uc.Assert(GenericsFlag != nil, "GenericsFlag must not be nil")

	// Add these assertions once MyGoLib is updated
	// type_name := uc.AssertNil(TypeNameFlag, "TypeNameFlag")
	// fields_value := uc.AssertNil(FieldsFlag, "FieldsFlag")
	// generics_value := uc.AssertNil(GenericsFlag, "GenericsFlag")

	flag.Parse()

	// Remove this once MyGoLib is updated
	if *TypeNameFlag == "" {
		return "", "", errors.New("the type name must be set")
	}

	// Remove this once MyGoLib is updated
	type_name := *TypeNameFlag

	var filename string

	if *OutputFileFlag == "" {
		var builder strings.Builder

		str := strings.ToLower(type_name)
		builder.WriteString(str)
		builder.WriteString("_treenode.go")

		filename = builder.String()
	} else {
		filename = *OutputFileFlag
	}

	align_generics(GenericsFlag, FieldsFlag)

	return type_name, filename, nil
}

type FieldsValue struct {
	fields   map[string]string
	generics map[rune]string
}

func (s *FieldsValue) String() string {
	var values []string
	var builder strings.Builder

	for name, value := range s.fields {
		builder.WriteString(value)
		builder.WriteRune(' ')
		builder.WriteString(name)

		str := builder.String()
		values = append(values, str)

		builder.Reset()
	}

	joined_str := strings.Join(values, ", ")
	quoted := strconv.Quote(joined_str)

	return quoted
}

func (s *FieldsValue) Set(value string) error {
	if value == "" {
		return errors.New("the value must be set")
	}

	fields := strings.Split(value, ",")

	parsed := make(map[string]string)

	for i, field := range fields {
		if field == "" {
			continue
		}

		sub_fields := strings.Split(field, "/")

		if len(sub_fields) == 1 {
			reason := errors.New("missing type")
			err := uc.NewErrAt(i+1, "field", reason)
			return err
		} else if len(sub_fields) > 2 {
			reason := errors.New("too many fields")
			err := uc.NewErrAt(i+1, "field", reason)
			return err
		}

		parsed[sub_fields[0]] = sub_fields[1]
	}

	s.fields = parsed

	// Find generics
	generics := make(map[rune]string)

	for _, field_type := range s.fields {
		chars, err := utgo.ParseGenerics(field_type)
		ok := utgo.IsErrNotGeneric(err)
		if ok {
			continue
		} else if err != nil {
			err := fmt.Errorf("syntax error for type %q: %w", field_type, err)
			return err
		}

		for _, char := range chars {
			_, ok := generics[char]
			if !ok {
				generics[char] = ""
			}
		}
	}

	return nil
}

func NewFieldsValue() *FieldsValue {
	fv := &FieldsValue{
		fields:   make(map[string]string),
		generics: make(map[rune]string),
	}

	return fv
}

type GenericsValue struct {
	letters []rune
	types   []string
}

func (s *GenericsValue) String() string {
	var values []string
	var builder strings.Builder

	for i, letter := range s.letters {
		builder.WriteRune(letter)
		builder.WriteRune(' ')
		builder.WriteString(s.types[i])

		str := builder.String()
		values = append(values, str)

		builder.Reset()
	}

	joined_str := strings.Join(values, ", ")

	builder.WriteRune('[')
	builder.WriteString(joined_str)
	builder.WriteRune(']')

	str := builder.String()
	return str
}

func (s *GenericsValue) Set(value string) error {
	fields := strings.Split(value, ",")

	for i, field := range fields {
		if field == "" {
			continue
		}

		gu, err := parse_generics_value(field)
		if err != nil {
			err := uc.NewErrAt(i+1, "field", err)
			return err
		}

		err = s.add(gu.letter, gu.g_type)
		if err != nil {
			return err
		}
	}

	return nil
}

func NewGenericsValue() *GenericsValue {
	gv := &GenericsValue{
		letters: make([]rune, 0),
		types:   make([]string, 0),
	}

	return gv
}

func (gv *GenericsValue) add(letter rune, g_type string) error {
	uc.AssertParam("letter", unicode.IsLetter(letter) && unicode.IsUpper(letter), errors.New("letter must be an upper case letter"))
	uc.AssertParam("g_type", g_type != "", errors.New("type must be set"))

	pos, ok := slices.BinarySearch(gv.letters, letter)
	if !ok {
		gv.letters = slices.Insert(gv.letters, pos, letter)
		gv.types = slices.Insert(gv.types, pos, g_type)

		return nil
	}

	if gv.types[pos] != g_type {
		err := fmt.Errorf("duplicate definition for generic %q: %s and %s", string(letter), gv.types[pos], g_type)
		return err
	}

	return nil
}

func (gv *GenericsValue) HasGenerics() bool {
	return len(gv.letters) > 0
}

func (gv *GenericsValue) GetGenericsList() string {
	values := make([]string, 0, len(gv.letters))

	for _, letter := range gv.letters {
		str := string(letter)
		values = append(values, str)
	}

	joined_str := strings.Join(values, ", ")

	var builder strings.Builder

	builder.WriteRune('[')
	builder.WriteString(joined_str)
	builder.WriteRune(']')

	str := builder.String()

	return str
}

type GenericsUnit struct {
	letter rune
	g_type string
}

func parse_generics_value(field string) (*GenericsUnit, error) {
	sub_fields := strings.Split(field, "/")

	if len(sub_fields) == 1 {
		return nil, errors.New("missing type of generic")
	} else if len(sub_fields) > 2 {
		return nil, errors.New("too many fields")
	}

	left := sub_fields[0]

	letter, err := utgo.IsGenericsID(left)
	if err != nil {
		return nil, err
	}

	right := sub_fields[1]

	gu := &GenericsUnit{
		letter: letter,
		g_type: right,
	}

	return gu, nil
}

func align_generics(gv *GenericsValue, fv *FieldsValue) {
	uc.AssertParam("gv", gv != nil, errors.New("gv must be set"))
	uc.AssertParam("fv", fv != nil, errors.New("fv must be set"))

	for generic_id := range fv.generics {
		pos, ok := slices.BinarySearch(gv.letters, generic_id)
		if !ok {
			gv.letters = slices.Insert(gv.letters, pos, generic_id)
			gv.types = slices.Insert(gv.types, pos, "any")
		}
	}
}
