package pkg

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"unicode"

	uc "github.com/PlayerR9/MyGoLib/Units/common"
	ggen "github.com/PlayerR9/MyGoLib/go_generator"
)

func MakeParameterList(fields map[string]string) (string, error) {
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

func MakeAssignmentList(fields map[string]string) (map[string]string, error) {
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

type GenData struct {
	PackageName   string
	TypeName      string
	TypeSig       string
	Fields        map[string]string
	IteratorName  string
	IteratorSig   string
	ParamList     string
	AssignmentMap map[string]string
	Generics      string
	Noder         string
}

func (g GenData) SetPackageName(pkg_name string) ggen.Generater {
	g.PackageName = pkg_name
	return g
}
