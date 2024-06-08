package object

import (
	"strings"
)

var Len = &Builtin{
	Fn: func(args ...Object) Object {
		if len(args) != 1 {
			return NewError("wrong number of arguments. got=%d, want=1", len(args))
		}
		switch arg := args[0].(type) {
		case *String:
			return &Integer{Value: int64(len(arg.Value))}
		case *Array:
			return &Integer{Value: int64(len(arg.Elements))}
		default:
			return NewError("argument to `len` not supported, got %s", args[0].Type())
		}
	},
}

var Reverse = &Builtin{
	Fn: func(args ...Object) Object {
		if len(args) != 1 {
			return NewError("wrong number of arguments. got=%d, want=1", len(args))
		}
		switch arg := args[0].(type) {
		case *String:
			ogRunes := []rune(arg.Value)
			length := len(ogRunes)
			runes := make([]rune, length)
			for i, v := range ogRunes {
				runes[length-i-1] = v
			}
			return &String{Value: string(runes)}
		case *Array:
			length := len(arg.Elements)
			elements := make([]Object, length)
			for i, v := range arg.Elements {
				elements[length-i-1] = v
			}

			return &Array{Elements: elements}
		default:
			return NewError("argument to `reverse` not supported, got %s", args[0].Type())
		}
	},
}

var Join = &Builtin{
	Fn: func(args ...Object) Object {
		switch arg := args[0].(type) {
		case *String:
			if len(args) != 3 {
				return NewError("wrong number of arguments. got=%d, want=3", len(args))
			}
			if args[1].Type() != STRING_OBJ {
				return NewError("argument to `join` must be STRING, got %s", args[1].Type())
			}
			if args[2].Type() != STRING_OBJ {
				return NewError("argument to `join` must be STRING, got %s", args[2].Type())
			}
			separator := args[2].(*String).Value
			arg1 := args[1].(*String).Value
			return &String{Value: arg.Value + separator + arg1}
		case *Array:
			separator := args[1].(*String).Value
			if len(args) != 2 {
				return NewError("wrong number of arguments. got=%d, want=2", len(args))
			}
			if args[1].Type() != STRING_OBJ {
				return NewError("argument to `join` must be STRING, got %s", args[1].Type())
			}
			elements := []string{}
			for _, v := range arg.Elements {
				elements = append(elements, v.Inspect())
			}
			return &String{Value: strings.Join(elements, separator)}
		default:
			return NewError("argument to `join` not supported, got %s", args[0].Type())
		}
	},
}
