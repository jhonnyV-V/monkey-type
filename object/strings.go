package object

import (
	"strings"
)

var Split = &Builtin{
	Fn: func(args ...Object) Object {
		if len(args) != 2 {
			return NewError("wrong number of arguments. got=%d, want=2", len(args))
		}
		if args[1].Type() != STRING_OBJ {
			return NewError("argument to `split` must be STRING, got %s", args[1].Type())
		}
		separator := args[1].(*String).Value
		switch arg := args[0].(type) {
		case *String:
			elements := []Object{}
			for _, v := range strings.Split(arg.Value, separator) {
				elements = append(elements, &String{Value: v})
			}
			return &Array{Elements: elements}
		default:
			return NewError("argument to `split` not supported, got %s", args[0].Type())
		}
	},
}

var Replace = &Builtin{
	Fn: func(args ...Object) Object {
		if len(args) != 3 {
			return NewError("wrong number of arguments. got=%d, want=3", len(args))
		}
		if args[1].Type() != STRING_OBJ || args[2].Type() != STRING_OBJ {
			return NewError("argument to `replace` must be STRING, got %s", args[1].Type())
		}
		substring := args[1].(*String).Value
		replacement := args[2].(*String).Value
		switch arg := args[0].(type) {
		case *String:
			return &String{
				Value: strings.ReplaceAll(arg.Value, substring, replacement),
			}
		default:
			return NewError("argument to `replace` not supported, got %s", args[0].Type())
		}
	},
}
