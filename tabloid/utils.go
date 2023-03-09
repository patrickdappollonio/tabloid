package tabloid

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types/ref"
)

type DuplicateColumnTitleError struct {
	Title string
}

func (e *DuplicateColumnTitleError) Error() string {
	return fmt.Sprintf("duplicate column title found: %q -- unable to work with non-unique column titles", e.Title)
}

func fnKey(s string) string {
	s = strings.ToLower(s)

	out := make([]rune, 0, len(s))

	for _, v := range s {
		if unicode.IsLetter(v) || unicode.IsDigit(v) || v == ' ' || v == '-' {
			switch v {
			case ' ', '-':
				out = append(out, '_')
			default:
				out = append(out, v)
			}
		}
	}

	return string(out)
}

func defineUnaryFunction(name string, unaryFn func(ref.Val) ref.Val) cel.EnvOption {
	return cel.Function(name, cel.Overload(name+"_global_unary_string", []*cel.Type{cel.StringType}, cel.BoolType, cel.UnaryBinding(unaryFn)))
}

func defineUnaryMethod(name string, unaryFn func(ref.Val) ref.Val) cel.EnvOption {
	return cel.Function(name, cel.MemberOverload(name+"_method_unary_string", []*cel.Type{cel.StringType}, cel.BoolType, cel.UnaryBinding(unaryFn)))
}

func defineBinaryFunction(name string, binaryFn func(ref.Val, ref.Val) ref.Val) cel.EnvOption {
	return cel.Function(name, cel.Overload(name+"_global_binary_string_string", []*cel.Type{cel.StringType, cel.StringType}, cel.BoolType, cel.BinaryBinding(binaryFn)))
}

func defineBinaryMethod(name string, binaryFn func(ref.Val, ref.Val) ref.Val) cel.EnvOption {
	return cel.Function(name, cel.MemberOverload(name+"_method_binary_string_string", []*cel.Type{cel.StringType, cel.StringType}, cel.BoolType, cel.BinaryBinding(binaryFn)))
}
