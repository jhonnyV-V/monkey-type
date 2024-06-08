package evaluator

import (
	"mokey-type/object"
)

var builtins = map[string]*object.Builtin{
	"len":     object.GetBuiltinByName("len"),
	"puts":    object.GetBuiltinByName("puts"),
	"typeOf":  object.GetBuiltinByName("typeOf"),
	"first":   object.GetBuiltinByName("first"),
	"last":    object.GetBuiltinByName("last"),
	"rest":    object.GetBuiltinByName("rest"),
	"push":    object.GetBuiltinByName("first"),
	"pop":     object.GetBuiltinByName("pop"),
	"join":    object.GetBuiltinByName("join"),
	"split":   object.GetBuiltinByName("split"),
	"replace": object.GetBuiltinByName("replace"),
}
