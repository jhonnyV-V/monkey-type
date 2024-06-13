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

var ToLower = &Builtin{
	Fn: func(args ...Object) Object {
		if len(args) != 1 {
			return NewError("wrong number of arguments. got=%d, want=1", len(args))
		}
		if args[0].Type() != STRING_OBJ {
			return NewError("argument to `toLower` must be STRING, got %s", args[0].Type())
		}

		return &String{Value: strings.ToLower(args[0].(*String).Value)}
	},
}

var ToUpper = &Builtin{
	Fn: func(args ...Object) Object {
		if len(args) != 1 {
			return NewError("wrong number of arguments. got=%d, want=1", len(args))
		}
		if args[0].Type() != STRING_OBJ {
			return NewError("argument to `toUpper` must be STRING, got %s", args[0].Type())
		}

		return &String{Value: strings.ToUpper(args[0].(*String).Value)}
	},
}

var Trim = &Builtin{
	Fn: func(args ...Object) Object {
		if len(args) != 2 {
			return NewError("wrong number of arguments. got=%d, want=2", len(args))
		}
		if args[0].Type() != STRING_OBJ && args[1].Type() != STRING_OBJ {
			return NewError("argument to `trim` must be STRING, got %s %s", args[0].Type(), args[1].Type())
		}

		value := args[0].(*String).Value
		cutset := args[1].(*String).Value

		return &String{Value: strings.Trim(value, cutset)}
	},
}

var TrimLeft = &Builtin{
	Fn: func(args ...Object) Object {
		if len(args) != 2 {
			return NewError("wrong number of arguments. got=%d, want=2", len(args))
		}
		if args[0].Type() != STRING_OBJ && args[1].Type() != STRING_OBJ {
			return NewError("argument to `trimLeft` must be STRING, got %s %s", args[0].Type(), args[1].Type())
		}

		value := args[0].(*String).Value
		cutset := args[1].(*String).Value

		return &String{Value: strings.TrimLeft(value, cutset)}
	},
}

var TrimRight = &Builtin{
	Fn: func(args ...Object) Object {
		if len(args) != 2 {
			return NewError("wrong number of arguments. got=%d, want=2", len(args))
		}
		if args[0].Type() != STRING_OBJ && args[1].Type() != STRING_OBJ {
			return NewError("argument to `trimRight` must be STRING, got %s %s", args[0].Type(), args[1].Type())
		}

		value := args[0].(*String).Value
		cutset := args[1].(*String).Value

		return &String{Value: strings.TrimRight(value, cutset)}
	},
}
