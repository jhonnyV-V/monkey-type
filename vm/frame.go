package vm

import (
	"mokey-type/code"
	"mokey-type/object"
)

type Frame struct {
	fn          *object.CompiledFunction
	ip          int
	basePointer int
}

func NewFrame(fn *object.CompiledFunction, basePointer int) *Frame {
	//TODO: change ip to minus 1 to fix bug with indexing in the vm
	return &Frame{fn: fn, ip: 0, basePointer: basePointer}
}

func (f *Frame) Instructions() code.Instructions {
	return f.fn.Instructions
}
