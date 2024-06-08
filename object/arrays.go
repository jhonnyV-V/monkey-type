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
		if length > 0 {
			newElements := make([]Object, length)
			copy(newElements, arr.Elements)
			newElements[length] = args[1]
			return &Array{Elements: newElements}
		}
		return nil
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
