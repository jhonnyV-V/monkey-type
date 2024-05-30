package vm

import (
	"fmt"
	"mokey-type/code"
	"mokey-type/compiler"
	"mokey-type/object"
)

const StackSize = 2048

type VM struct {
	constant     []object.Object
	instructions code.Instructions

	stack []object.Object
	sp    int //Always points to the next value. Top of the stack is stack[sp - 1]
}

func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		constant:     bytecode.Constanst,
		instructions: bytecode.Instructions,

		stack: make([]object.Object, StackSize),
		sp:    0,
	}
}

func (vm *VM) StackTop() object.Object {
	if vm.sp == 0 {
		return nil
	}

	return vm.stack[vm.sp-1]
}

func (vm *VM) Run() error {
	for ip := 0; ip < len(vm.instructions); ip++ {
		op := code.Opcode(vm.instructions[ip])

		switch op {
		case code.OpConstant:
			constIndex := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2
			err := vm.push(vm.constant[constIndex])
			if err != nil {
				return err
			}

		case code.OpAdd:
			right := vm.pop()
			left := vm.pop()
			leftValue := left.(*object.Integer).Value
			rightValue := right.(*object.Integer).Value
			vm.push(&object.Integer{Value: rightValue + leftValue})
		}
	}
	return nil
}

func (vm *VM) push(ob object.Object) error {
	if vm.sp >= StackSize {
		return fmt.Errorf("Stack Overflow")
	}

	vm.stack[vm.sp] = ob
	vm.sp++

	return nil
}

func (vm *VM) pop() object.Object {
	ob := vm.stack[vm.sp-1]
	vm.sp--

	return ob
}
