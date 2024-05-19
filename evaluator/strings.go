package evaluator

import (
	"mokey-type/object"
	"strings"
)

var Split = &object.Builtin{
	Fn: func(args ...object.Object) object.Object {
		if len(args) != 2 {
			return newError("wrong number of arguments. got=%d, want=2", len(args))
		}
		if args[1].Type() != object.STRING_OBJ {
			return newError("argument to `split` must be STRING, got %s", args[1].Type())
		}
		separator := args[1].(*object.String).Value
		switch arg := args[0].(type) {
		case *object.String:
			elements := []object.Object{}
			for _, v := range strings.Split(arg.Value, separator) {
				elements = append(elements, &object.String{Value: v})
			}
			return &object.Array{Elements: elements}
		default:
			return newError("argument to `split` not supported, got %s", args[0].Type())
		}
	},
}

var Replace = &object.Builtin{
	Fn: func(args ...object.Object) object.Object {
		if len(args) != 3 {
			return newError("wrong number of arguments. got=%d, want=3", len(args))
		}
		if args[1].Type() != object.STRING_OBJ || args[2].Type() != object.STRING_OBJ {
			return newError("argument to `replace` must be STRING, got %s", args[1].Type())
		}
		substring := args[1].(*object.String).Value
		replacement := args[2].(*object.String).Value
		switch arg := args[0].(type) {
		case *object.String:
			return &object.String{
				Value: strings.ReplaceAll(arg.Value, substring, replacement),
			}
		default:
			return newError("argument to `replace` not supported, got %s", args[0].Type())
		}
	},
}
