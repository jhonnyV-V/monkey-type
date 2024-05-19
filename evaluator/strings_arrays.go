package evaluator

import (
	"mokey-type/object"
	"strings"
)

var Len = &object.Builtin{
	Fn: func(args ...object.Object) object.Object {
		if len(args) != 1 {
			return newError("wrong number of arguments. got=%d, want=1", len(args))
		}
		switch arg := args[0].(type) {
		case *object.String:
			return &object.Integer{Value: int64(len(arg.Value))}
		case *object.Array:
			return &object.Integer{Value: int64(len(arg.Elements))}
		default:
			return newError("argument to `len` not supported, got %s", args[0].Type())
		}
	},
}

var Reverse = &object.Builtin{
	Fn: func(args ...object.Object) object.Object {
		if len(args) != 1 {
			return newError("wrong number of arguments. got=%d, want=1", len(args))
		}
		switch arg := args[0].(type) {
		case *object.String:
			ogRunes := []rune(arg.Value)
			length := len(ogRunes)
			runes := make([]rune, length)
			for i, v := range ogRunes {
				runes[length-i-1] = v
			}
			return &object.String{Value: string(runes)}
		case *object.Array:
			length := len(arg.Elements)
			elements := make([]object.Object, length)
			for i, v := range arg.Elements {
				elements[length-i-1] = v
			}

			return &object.Array{Elements: elements}
		default:
			return newError("argument to `reverse` not supported, got %s", args[0].Type())
		}
	},
}

var Join = &object.Builtin{
	Fn: func(args ...object.Object) object.Object {
		switch arg := args[0].(type) {
		case *object.String:
			if len(args) != 3 {
				return newError("wrong number of arguments. got=%d, want=3", len(args))
			}
			if args[1].Type() != object.STRING_OBJ {
				return newError("argument to `join` must be STRING, got %s", args[1].Type())
			}
			if args[2].Type() != object.STRING_OBJ {
				return newError("argument to `join` must be STRING, got %s", args[2].Type())
			}
			separator := args[2].(*object.String).Value
			arg1 := args[1].(*object.String).Value
			return &object.String{Value: arg.Value + separator + arg1}
		case *object.Array:
			separator := args[1].(*object.String).Value
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}
			if args[1].Type() != object.STRING_OBJ {
				return newError("argument to `join` must be STRING, got %s", args[1].Type())
			}
			elements := []string{}
			for _, v := range arg.Elements {
				elements = append(elements, v.Inspect())
			}
			return &object.String{Value: strings.Join(elements, separator)}
		default:
			return newError("argument to `join` not supported, got %s", args[0].Type())
		}
	},
}
