package evaluator

import (
	"fmt"
	"mokey-type/object"
)

var builtins = map[string]*object.Builtin{
	"puts": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return NULL
		},
	},
	"typeOf": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			return &object.String{Value: string(args[0].Type())}
		},
	},
	"len":     Len,
	"first":   First,
	"last":    Last,
	"rest":    Rest,
	"push":    Push,
	"pop":     Pop,
	"reverse": Reverse,
	"join":    Join,
	"split":   Split,
	"replace": Replace,
}
