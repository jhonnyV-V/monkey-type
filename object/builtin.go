package object

import (
	"fmt"
)

var Builtins = []struct {
	Name    string
	Builtin *Builtin
}{
	{
		"puts",
		&Builtin{
			Fn: func(args ...Object) Object {
				for _, arg := range args {
					fmt.Println(arg.Inspect())
				}
				return nil
			},
		},
	},
	{
		"typeOf",
		&Builtin{
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return NewError("wrong number of arguments. got=%d, want=1", len(args))
				}
				return &String{Value: string(args[0].Type())}
			},
		},
	},
	{
		"len",
		Len,
	},
	{
		"first",
		First,
	},
	{
		"last",
		Last,
	},
	{
		"rest",
		Rest,
	},
	{
		"push",
		Push,
	},
	{
		"pop",
		Pop,
	},
	{
		"reverse",
		Reverse,
	},
	{
		"join",
		Join,
	},
	{
		"split",
		Split,
	},
	{
		"replace",
		Replace,
	},
	{
		"toLower",
		ToLower,
	},
	{
		"toUpper",
		ToUpper,
	},
	{
		"trim",
		Trim,
	},
	{
		"trimLeft",
		TrimLeft,
	},
	{
		"trimRight",
		TrimRight,
	},
	{
		"contains",
		Contains,
	},
	{
		"merge",
		Merge,
	},
	{
		"findIndex",
		FindIndex,
	},
}

func NewError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}

func GetBuiltinByName(name string) *Builtin {
	for _, def := range Builtins {
		if def.Name == name {
			return def.Builtin
		}
	}
	return nil
}
