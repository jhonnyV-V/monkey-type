package object

var First = &Builtin{
	Fn: func(args ...Object) Object {
		if len(args) != 1 {
			return NewError("wrong number of arguments. got=%d, want=1", len(args))
		}
		if args[0].Type() != ARRAY_OBJ {
			return NewError("argument to `first` must be ARRAY, got %s", args[0].Type())
		}
		arr := args[0].(*Array)
		if len(arr.Elements) > 0 {
			return arr.Elements[0]
		}
		return nil
	},
}

var Last = &Builtin{
	Fn: func(args ...Object) Object {
		if len(args) != 1 {
			return NewError("wrong number of arguments. got=%d, want=1", len(args))
		}
		if args[0].Type() != ARRAY_OBJ {
			return NewError("argument to `last` must be ARRAY, got %s", args[0].Type())
		}
		arr := args[0].(*Array)
		length := len(arr.Elements)
		if length > 0 {
			return arr.Elements[length-1]
		}
		return nil
	},
}

var Rest = &Builtin{
	Fn: func(args ...Object) Object {
		if len(args) != 1 {
			return NewError("wrong number of arguments. got=%d, want=1", len(args))
		}
		if args[0].Type() != ARRAY_OBJ {
			return NewError("argument to `rest` must be ARRAY, got %s", args[0].Type())
		}
		arr := args[0].(*Array)
		length := len(arr.Elements)
		if length > 0 {
			newElements := make([]Object, length-1)
			copy(newElements, arr.Elements[1:length])
			return &Array{Elements: newElements}
		}
		return nil
	},
}

var Push = &Builtin{
	Fn: func(args ...Object) Object {
		if len(args) != 2 {
			return NewError("wrong number of arguments. got=%d, want=2", len(args))
		}
		if args[0].Type() != ARRAY_OBJ {
			return NewError("argument to `push` must be ARRAY, got %s", args[0].Type())
		}
		arr := args[0].(*Array)
		length := len(arr.Elements)
		newElements := make([]Object, length+1)
		copy(newElements, arr.Elements)
		newElements[length] = args[1]
		return &Array{Elements: newElements}
	},
}

var Pop = &Builtin{
	Fn: func(args ...Object) Object {
		if len(args) != 1 {
			return NewError("wrong number of arguments. got=%d, want=1", len(args))
		}
		if args[0].Type() != ARRAY_OBJ {
			return NewError("argument to `pop` must be ARRAY, got %s", args[0].Type())
		}
		arr := args[0].(*Array)
		length := len(arr.Elements)
		if length > 0 {
			newElements := make([]Object, length-1)
			copy(newElements, arr.Elements[:length-1])
			return &Array{Elements: newElements}
		}
		return arr
	},
}

var Merge = &Builtin{
	Fn: func(args ...Object) Object {
		if len(args) != 2 {
			return NewError("wrong number of arguments. got=%d, want=2", len(args))
		}
		if args[0].Type() != ARRAY_OBJ || args[1].Type() != ARRAY_OBJ {
			return NewError("argument to `push` must be ARRAY, got %s, %s", args[0].Type(), args[1].Type())
		}

		arrA := args[0].(*Array)
		arrB := args[1].(*Array)
		newElements := append(arrA.Elements, arrB.Elements...)
		return &Array{Elements: newElements}
	},
}

var FindIndex = &Builtin{
	Fn: func(args ...Object) Object {
		if len(args) != 2 {
			return NewError("wrong number of arguments. got=%d, want=2", len(args))
		}
		switch arg := args[0].(type) {
		case *Array:
			var result Integer

			for i, element := range arg.Elements {
				if element.Type() == args[1].Type() {

					switch element := element.(type) {
					case *Boolean:
						value := args[1].(*Boolean).Value
						if element.Value == value {
							result.Value = int64(i)
						}

					case *Integer:
						value := args[1].(*Integer).Value
						if element.Value == value {
							result.Value = int64(i)
						}

					case *String:
						value := args[1].(*String).Value
						if element.Value == value {
							result.Value = int64(i)
						}

					case *NullValue:
						result.Value = int64(i)

					default:
						return NewError("argument to `contains` not supported, got %s", args[1].Type())
					}
				}
			}
			return &result
		default:
			return NewError("argument to `findIndex` not supported, got %s", args[0].Type())
		}
	},
}
