package vm

import (
	"mokey-type/code"
	"mokey-type/object"
)

type Frame struct {
	cl          *object.Closure
	ip          int
	basePointer int
}

func NewFrame(fn *object.Closure, basePointer int) *Frame {
	//TODO: change ip to minus 1 to fix bug with indexing in the vm
	return &Frame{cl: fn, ip: 0, basePointer: basePointer}
}

func (f *Frame) Instructions() code.Instructions {
	return f.cl.Fn.Instructions
}
